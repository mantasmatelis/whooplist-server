package whooplist

import (
	"code.google.com/p/go.crypto/scrypt"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"github.com/imdario/mergo"
	"io"
	"log"
	"strings"
	"time"
)

type User struct {
	Id       *int64
	Email    *string
	Name     *string
	Fname    *string    `json:",omitempty"`
	Lname    *string    `json:",omitempty"`
	Birthday *time.Time `json:",omitempty"`
	School   *string    `json:",omitempty"`
	Picture  *string    `json:",omitempty"`
	Gender   *int64     `json:",omitempty"` /* Follows ISO/IEC 5218 */

	/* Private (i.e. not passed to client ever) */
	PasswordHash string `json:"-"`
	Role         string `json:"-"`

	/* For the client to pass to user only */
	Password    *string `json:",omitempty"`
	OldPassword *string `json:",omitempty"`

	/* For friend suggestion */
	Score *int `json:",omitempty"`
}

/* Sessions offer long-term authenticated communication between
   server and client. Identity is proved with a long key
   (that must be kept absolutely secret at all times).
   Sessions are signed out after 28 days of inactivity
   (not making requests under the key) for security. */
type Session struct {
	Id       int64
	UserId   int64
	Key      string
	LastAuth time.Time
	LastUse  time.Time
}

var getUserDataStmt, createUserStmt, updateUserStmt, deleteUserStmt,
	authUserStmt, loginUserVerifyStmt, loginUserStmt,
	deleteSessionStmt, existsUserStmt *sql.Stmt

func prepareUser() {
	stmt(&getUserDataStmt,
		"SELECT id, email, name, fname, lname, birthday, "+
			"school, picture, gender, password_hash, role "+
			"FROM wl.user WHERE id = $1 OR email = $2;")

	stmt(&createUserStmt,
		"INSERT INTO wl.user (email, name, fname, lname, birthday, "+
			"school, picture, gender, password_hash, role) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) "+
			"RETURNING id;")

	stmt(&updateUserStmt,
		"UPDATE wl.user SET name = $2, fname = $3, lname = $4, "+
			"birthday = $5, school = $6, "+
			"picture = $7, gender = $8,  password_hash = $9, "+
			"role = $10 WHERE email = $1;")

	stmt(&deleteUserStmt,
		"DELETE FROM wl.user WHERE id = $1;")

	stmt(&authUserStmt,
		"SELECT * FROM wl.session WHERE key = $1;")

	stmt(&loginUserVerifyStmt,
		"SELECT id, email, name, fname, lname, birthday, school, "+
			"picture, gender, role "+
			"FROM wl.user "+
			"WHERE lower(email) = lower($1) AND password_hash = $2;")

	stmt(&loginUserStmt,
		"INSERT INTO wl.session (user_id, key, last_auth, last_use) "+
			"VALUES ($1, $2, NOW(), NOW());")

	stmt(&deleteSessionStmt,
		"DELETE FROM wl.session WHERE key = $1;")

	stmt(&existsUserStmt,
		"SELECT id FROM wl.user WHERE lower(email)= $1;")
}

func Hash(email, password string) (hash string, err error) {
	secret := "aYdZYlE9ybGXn5CldvQ3f/shKxNshtAOvqDlaw/wbUBHwc5r9zBal9hf9CDkGxSgddAMtNm+uz1G"
	secretData, err := base64.StdEncoding.DecodeString(secret)

	if err != nil {
		return
	}

	salt := make([]byte, len(secretData)+len(email))

	copy(salt, secretData)
	copy(salt[len(secretData):], []byte(strings.ToLower(email)))

	hash_data, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	hash = base64.StdEncoding.EncodeToString(hash_data)

	return
}

func CheckPassword(password string) bool {
	return len(password) > 6
}

func GetUserData(id int64, email string) (user *User, err error) {
	res := getUserDataStmt.QueryRow(id, email)
	user = new(User)
	err = res.Scan(&user.Id, &user.Email, &user.Name, &user.Fname,
		&user.Lname, &user.Birthday, &user.School, &user.Picture, &user.Gender,
		&user.PasswordHash, &user.Role)

	if err == sql.ErrNoRows {
		user = nil
		err = nil
	} else if err != nil {
		user = nil
	}

	return

}

func CreateUser(user *User) (err error) {
	if !CheckPassword(*user.Password) {
		return WeakPassword
	}
	if user.Picture != nil {
		str, err := WriteFileBase64("profile.jpg", user.Picture, false)
		if err != nil {
			return err
		}
		user.Picture = &str
	}

	hash, err := Hash(*user.Email, *user.Password)

	if err != nil {
		return
	}

	res := createUserStmt.QueryRow(user.Email, user.Name, user.Fname,
		user.Lname, user.Birthday, user.School, user.Picture,
		user.Gender, hash, user.Role)

	err = res.Scan(&user.Id)

	if err != nil {
		return
	}

	//TODO: Get USER Id on query
	AddNewsfeedItem(
		&FeedItem{Type: NfWelcome, UserId: *user.Id, Picture: *user.Picture})

	return
}

func UpdateUser(oldUser, user *User) (err error) {
	user.Email = oldUser.Email
	user.Id = oldUser.Id
	user.PasswordHash = oldUser.PasswordHash
	user.Role = oldUser.Role

	if user.Password != nil {
		if !CheckPassword(*user.Password) {
			return WeakPassword
		}

		hash, err := Hash(*user.Email, *user.OldPassword)
		if err != nil {
			return err
		}

		var blankUser User
		res := loginUserVerifyStmt.QueryRow(user.Email, hash)
		err = res.Scan(&blankUser.Id, &blankUser.Email, &blankUser.Name,
			&blankUser.Fname, &blankUser.Lname, &blankUser.Birthday,
			&blankUser.School, &blankUser.Picture, &blankUser.Gender,
			&blankUser.Role)

		if err != nil {
			if err == sql.ErrNoRows {
				return BadPassword
			}
			return err
		}

		user.PasswordHash, err = Hash(*user.Email, *user.Password)
		if err != nil {
			return err
		}
	}

	if user.Picture != nil && !strings.HasPrefix(*user.Picture, baseUrl) {
		log.Print("updating picture")
		str, err := WriteFileBase64("profile.jpg", user.Picture, false)
		if err != nil {
			return err
		}
		user.Picture = &str
		log.Print("newuser.picture = ", user.Picture)
	}

	err = mergo.Merge(user, *oldUser)

	log.Print("after merge, newuser.picture = ", user.Picture)

	if err != nil {
		return err
	}

	_, err = updateUserStmt.Exec(user.Email, user.Name, user.Fname,
		user.Lname, user.Birthday, user.School, user.Picture,
		user.Gender, user.PasswordHash, user.Role)

	if err != nil {
		return err
	}

	if user.Picture != oldUser.Picture {
		AddNewsfeedItem(&FeedItem{Type: NfProfilePictureUpdated,
			UserId: *user.Id, Picture: *user.Picture})
	}

	if user.School != oldUser.School {
		AddNewsfeedItem(&FeedItem{Type: NfSchoolUpdated,
			UserId: *user.Id, Picture: *user.Picture, AuxString: *user.School})
	}
	return
}

func DeleteUser(userId int64) (err error) {
	//TODO: Delete a user's lists as well!
	_, err = deleteUserStmt.Exec(userId)
	return
}

func AuthUser(key string) (user *User, session *Session, err error) {
	res := authUserStmt.QueryRow(key)
	session = new(Session)
	err = res.Scan(&session.Id, &session.UserId, &session.Key,
		&session.LastAuth, &session.LastUse)

	if err == sql.ErrNoRows {
		session = nil
		err = nil
		return
	}
	if err != nil {
		return
	}

	user, err = GetUserData(session.UserId, "")

	return
}

func LoginUser(username, password string) (user *User,
	session *Session, err error) {

	hash, err := Hash(username, password)

	if err != nil {
		return
	}

	res := loginUserVerifyStmt.QueryRow(username, hash)
	user = new(User)
	err = res.Scan(&user.Id, &user.Email, &user.Name, &user.Fname,
		&user.Lname, &user.Birthday, &user.School, &user.Picture,
		&user.Gender, &user.Role)

	if err == sql.ErrNoRows {
		err = nil
		user = nil
		return
	}
	if err != nil {
		user = nil
		return
	}

	/* Generate a session key. Currently 192 bits of entropy. */
	data := make([]byte, 24)
	n, err := io.ReadFull(rand.Reader, data)

	if n != len(data) || err != nil {
		user = nil
		return
	}

	key := base64.StdEncoding.EncodeToString(data)

	_, err = loginUserStmt.Exec(user.Id, key)

	if err != nil {
		return
	}

	/* TODO: We don't fill the whole session,
	we should identify if this is a problem. */

	session = new(Session)
	session.UserId = *user.Id
	session.Key = key

	return
}

func DeleteSession(key string) (exist bool, err error) {
	res, err := deleteSessionStmt.Exec(key)
	if err == nil {
		rows, _ := res.RowsAffected()
		if rows == 1 {
			exist = true
		}
	}
	return
}

func UserExists(email string) (exist bool, err error) {
	res, err := existsUserStmt.Exec(strings.ToLower(email))
	if err == nil {
		rows, _ := res.RowsAffected()
		if rows == 1 {
			exist = true
		}
	}
	return
}

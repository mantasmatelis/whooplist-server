package whooplist

import (
	"code.google.com/p/go.crypto/scrypt"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/imdario/mergo"
	"github.com/kisielk/sqlstruct"
	_ "github.com/lib/pq"
	"io"
	"log"
	"strconv"
	"strings"
)

var UserError = errors.New("whooplist: user error")
var WeakPassword = errors.New("whooplist: weak password error")
var BadPassword = errors.New("whooplist: bad password error")

/* Application global database connection pool */
var db *sql.DB

var getUserDataStmt, createUserStmt, updateUserStmt, deleteUserStmt,
	authUserStmt, loginUserVerifyStmt, loginUserStmt,
	deleteSessionStmt, existsUserStmt, getUserListsStmt, getUserListStmt,
	getUserListPlaceStmt, putUserListStmt, deleteUserListStmt,
	networkUserFriendsStmt, suggestUserFriendsStmt, contactsUserFriendsStmt,
	getUserFriendsStmt, addUserFriendStmt,
	deleteUserFriendStmt, getListTypesStmt, getWhooplistCoordinateStmt,
	getWhooplistLocationStmt, addNewsfeedItemStmt, getNewsfeedStmt,
	getNewsfeedEarlierStmt, getLocationStmt, getLocationCoordinateStmt,
	getPlaceStmt, getPlaceByFactualStmt, addPlaceStmt,
	updatePlaceStmt *sql.Stmt

func prepare() (err error) {
	getUserDataStmt, err = db.Prepare(
		"SELECT id, email, name, fname, lname, birthday, " +
			"school, picture, gender, password_hash, role " +
			"FROM wl.user WHERE id = $1 OR email = $2;")
	if err != nil {
		return
	}

	createUserStmt, err = db.Prepare(
		"INSERT INTO wl.user (email, name, fname, lname, birthday, " +
			"school, picture, gender, password_hash, role) " +
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);")
	if err != nil {
		return
	}

	updateUserStmt, err = db.Prepare(
		"UPDATE wl.user SET name = $2, fname = $3, lname = $4,  birthday = $5, school = $6, " +
			"picture = $7, gender = $8,  password_hash = $9, " +
			"role = $10 WHERE email = $1;")
	if err != nil {
		return
	}

	deleteUserStmt, err = db.Prepare(
		"DELETE FROM wl.user WHERE id = $1;")
	if err != nil {
		return
	}

	authUserStmt, err = db.Prepare(
		"SELECT * FROM wl.session WHERE key = $1;")
	if err != nil {
		return
	}

	loginUserVerifyStmt, err = db.Prepare(
		"SELECT id, email, name, fname, lname, birthday, school, picture, gender, role " +
			"FROM wl.user WHERE lower(email) = lower($1) AND password_hash = $2;")
	if err != nil {
		return
	}

	loginUserStmt, err = db.Prepare(
		"INSERT INTO wl.session (user_id, key, last_auth, last_use) " +
			"VALUES ($1, $2, NOW(), NOW());")
	if err != nil {
		return
	}

	deleteSessionStmt, err = db.Prepare(
		"DELETE FROM wl.session WHERE key = $1;")
	if err != nil {
		return
	}

	existsUserStmt, err = db.Prepare(
		"SELECT id FROM wl.user WHERE lower(email)= $1;")
	if err != nil {
		return
	}

	getUserListsStmt, err = db.Prepare(
		"SELECT list_id FROM wl.list_item WHERE user_id = $1 " +
			"GROUP BY list_id;")
	if err != nil {
		return
	}

	getUserListStmt, err = db.Prepare(
		"SELECT place.id, place.latitude, place.longitude, " +
			"place.factual_id, place.name, place.address, " +
			"place.locality, place.region, place.postcode, place.country, " +
			"place.telephone, place.website, place.email " +
			"FROM wl.list_item JOIN wl.place " +
			"ON list_item.place_id = place.id " +
			"WHERE list_item.user_id = $1 AND list_item.list_id = $2 " +
			"ORDER BY rank")
	if err != nil {
		return
	}

	putUserListStmt, err = db.Prepare(
		"INSERT INTO wl.list_item (place_id, list_id, user_id, rank) " +
			"VALUES ($1, $2, $3, $4)")
	if err != nil {
		return
	}

	deleteUserListStmt, err = db.Prepare(
		"DELETE FROM wl.list_item WHERE user_id = $1 AND list_id = $2")
	if err != nil {
		return
	}

	getUserFriendsStmt, err = db.Prepare(
		"SELECT SUM(direction), id, email, name, fname, lname, birthday, " +
			"school, picture, gender FROM " +
			"(SELECT 1 AS direction, wl.user.id AS id, email, name, fname, " +
			"lname, birthday, school, picture, gender, role FROM wl.user " +
			"JOIN wl.friend ON wl.user.id = friend.from_id " +
			"WHERE friend.to_id = $1 UNION " +
			"SELECT 2 AS direction, wl.user.id, email, name, fname, lname, " +
			"birthday, school, picture, gender, role FROM wl.user " +
			"JOIN wl.friend ON wl.user.id = friend.to_id " +
			"WHERE friend.from_id = $1) AS bothDirections GROUP BY id, " +
			"email, name, fname, lname, birthday, school, picture, gender, role;")
	if err != nil {
		return
	}

	addUserFriendStmt, err = db.Prepare(
		"INSERT INTO wl.friend (from_id, to_id) VALUES ($1, $2);")
	if err != nil {
		return
	}

	deleteUserFriendStmt, err = db.Prepare(
		"DELETE FROM wl.friend WHERE from_id = $1 AND to_id = $2;")
	if err != nil {
		return
	}

	/*	suggestUserFriendsStmt, err = db.Prepare(
		"SELECT COUNT(*) user.id, email, name, fname, lname, birthday, " +
		"school, picture, gender FROM wl.friend " +
		"JOIN user ON user.id= to_id " +
		"WHERE from_id IN " +
		"(SELECT DISTINCT to_id FROM friend WHERE from_id = $1) AND "
		"user.id NOT IN " +
		"(SELECT DISTINCT to_id FROM friend WHERE from_id = $1) " +
		"GROUP BY user.id")*/

	networkUserFriendsStmt, err = db.Prepare(
		"SELECT COUNT(*), wl.user.id, email, name, fname, lname, birthday, " +
			"school, picture, gender FROM wl.friend " +
			"JOIN wl.user ON wl.user.id = to_id " +
			"WHERE from_id IN " +
			"(SELECT DISTINCT to_id FROM friend WHERE from_id = $1) " +
			"AND wl.user.id NOT IN " +
			"(SELECT DISTINCT to_id FROM friend WHERE from_id = $1) " +
			"AND wl.user.id <> $1 " +
			"GROUP BY wl.user.id " +
			"ORDER BY count DESC " +
			"LIMIT 10;")
	if err != nil {
		return
	}

	suggestUserFriendsStmt, err = db.Prepare(
		"SELECT COUNT(*), wl.user.id, email, name, fname, lname, birthday, " +
			"school, picture, gender FROM wl.friend " +
			"JOIN wl.user ON wl.user.id = to_id " +
			"WHERE from_id IN " +
			"(SELECT DISTINCT to_id FROM friend WHERE from_id = $1) " +
			"AND wl.user.id NOT IN " +
			"(SELECT DISTINCT to_id FROM friend WHERE from_id = $1) " +
			"AND wl.user.id <> $1 " +
			"GROUP BY wl.user.id " +
			"ORDER BY count DESC " +
			"LIMIT 10;")
	if err != nil {
		return
	}

	contactsUserFriendsStmt, err = db.Prepare(
		"SELECT id, email, name, fname, lname, " +
			"birthday, school, picture, gender FROM wl.user " +
			"WHERE email IN (SELECT unnest (string_to_array($1, '&'))) " +
			"OR phone IN (SELECT unnest (string_to_array($2, '&'))) " +
			"GROUP BY wl.user.id;")
	if err != nil {
		return
	}

	getListTypesStmt, err = db.Prepare(
		"SELECT * FROM wl.list")
	if err != nil {
		return
	}

	/* Note, the constant below is 360 / (earth's radius in meters) */
	getWhooplistCoordinateStmt, err = db.Prepare(
		"SELECT SUM(10 - rank) AS score, place.id, place.latitude, place.longitude, " +
			"place.factual_id, place.name, place.address, " +
			"place.locality, place.region, place.postcode, place.country, " +
			"place.telephone, place.website, place.email " +
			"FROM list_item JOIN place " +
			"ON list_item.place_id = place.id " +
			"WHERE list_item.list_id = $1 AND " +
			"latitude + (($4 * 0.00000898314) / cos(latitude)) > $2 AND " +
			"latitude - (($4 * 0.00000898314) / cos(latitude)) < $2 AND " +
			"longitude + ($4 * 0.00000898314) > $3 AND " +
			"longitude - ($4 * 0.00000898314) < $3 " +
			"GROUP BY place.id " +
			"LIMIT 10 OFFSET $5")
	if err != nil {
		return
	}

	addNewsfeedItemStmt, err = db.Prepare(
		"INSERT INTO wl.feed_item (user_id, latitude, longitude, " +
			"place_id, list_id, picture, type, aux_string, aux_int) " +
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)")
	if err != nil {
		return
	}

	getNewsfeedStmt, err = db.Prepare(
		"SELECT id, timestamp, user_id, latitude, longitude, place_id, " +
			"list_id, picture, type, aux_string, aux_int " +
			"FROM wl.feed_item " +
			"WHERE (" +
			"(user_id = $1) " +
			"OR " +
			"(latitude + (($5 * 0.00000898314) / cos(latitude)) > $3 AND " +
			"latitude - (($5 * 0.00000898314) / cos(latitude)) < $3 AND " +
			"longitude + ($5 * 0.00000898314) > $4 AND " +
			"longitude - ($5 * 0.00000898314) < $4)) " +
			"AND (id > $2 OR $2 = -1) " +
			"LIMIT 30")
	if err != nil {
		return
	}

	getNewsfeedEarlierStmt, err = db.Prepare(
		"SELECT id, timestamp, user_id, latitude, longitude, place_id, " +
			"list_id, picture, type, aux_string, aux_int " +
			"FROM wl.feed_item " +
			"WHERE (" +
			"(user_id = $1) " +
			"OR " +
			"(latitude + (($5 * 0.00000898314) / cos(latitude)) > $3 AND " +
			"latitude - (($5 * 0.00000898314) / cos(latitude)) < $3 AND " +
			"longitude + ($5 * 0.00000898314) > $4 AND " +
			"longitude - ($5 * 0.00000898314) < $4)) " +
			"AND id < $2 LIMIT 30")
	if err != nil {
		return
	}

	getPlaceStmt, err = db.Prepare(
		"SELECT id, latitude, longitude, factual_id, name, address, locality, " +
			"region, postcode, country, telephone, website, email " +
			"FROM wl.place WHERE id = $1;")
	if err != nil {
		return
	}

	getPlaceByFactualStmt, err = db.Prepare(
		"SELECT latitude, longitude, factual_id, name, address, locality, " +
			"region, postcode, country, telephone, website, email " +
			"FROM wl.place WHERE factual_id = $1;")
	if err != nil {
		return
	}

	updatePlaceStmt, err = db.Prepare(
		"UPDATE wl.place SET latitude=$1, longitude=$2, factual_id=$3, " +
			"name=$4, address=$5, locality=$6, region=$7, postcode=$8, " +
			"country=$9, telephone=$10, website=$11, email=$12 " +
			"WHERE factual_id=$3 RETURNING id;")
	if err != nil {
		return
	}

	addPlaceStmt, err = db.Prepare(
		"INSERT INTO wl.place (latitude, longitude, factual_id, name, " +
			"address, locality, region, postcode, country, telephone, " +
			"website, email) " +
			"SELECT $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12 " +
			"WHERE NOT EXISTS (SELECT 1 FROM place WHERE factual_id=$3) " +
			"RETURNING id;")
	if err != nil {
		return
	}

	return
}

func Hash(username, password string) (hash string, err error) {
	secret := "aYdZYlE9ybGXn5CldvQ3f/shKxNshtAOvqDlaw/wbUBHwc5r9zBal9hf9CDkGxSgddAMtNm+uz1G"
	secretData, err := base64.StdEncoding.DecodeString(secret)

	if err != nil {
		return
	}

	salt := make([]byte, len(secretData)+len(username))

	copy(salt, secretData)
	copy(salt[len(secretData):], []byte(strings.ToLower(username)))

	hash_data, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32)
	hash = base64.StdEncoding.EncodeToString(hash_data)

	return
}

func CheckPassword(password string) bool {
	return len(password) > 6
}

func Initialize() (err error) {
	db, err = sql.Open("postgres",
		"user=whooplist dbname=whooplist password=moteifae0ohcaiCo "+
			"sslmode=disable search_path=wl host=localhost")

	if err == nil {
		err = prepare()
	}

	initializeOauth()

	return
}

func Disconnect() (err error) {
	err = db.Close()
	return
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

	_, err = createUserStmt.Exec(user.Email, user.Name, user.Fname,
		user.Lname, user.Birthday, user.School, user.Picture,
		user.Gender, hash, user.Role)
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

		hash, err := Hash(*user.Email, *user.Password)
		if err != nil {
			return err
		}

		var blankUsers []User
		res, err := loginUserVerifyStmt.Query(user.Email, hash)

		if err != nil {
			return err
		}

		err = sqlstruct.Scan(&blankUsers, res)

		if len(blankUsers) != 1 {
			return BadPassword
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

func GetUserLists(userId int64) (lists []int, err error) {
	rows, err := getUserListsStmt.Query(userId)

	if err == sql.ErrNoRows {
		err = nil
		return
	} else if err != nil {
		return
	}

	lists = make([]int, 0, 10)

	for rows.Next() {
		var list int
		err = rows.Scan(&list)
		lists = append(lists, list)
		if err != nil {
			lists = nil
			return
		}
	}
	return
}

func GetUserList(userId, listId int64) (list []Place, err error) {
	list = make([]Place, 0, 5)

	rows, err := getUserListStmt.Query(userId, listId)

	if err == sql.ErrNoRows {
		list = nil
		err = nil
		return
	} else if err != nil {
		list = nil
		return
	}

	for rows.Next() {
		var curr Place
		err = rows.Scan(&curr.Id, &curr.Latitude, &curr.Longitude,
			&curr.FactualId, &curr.Name, &curr.Address, &curr.Locality,
			&curr.Region, &curr.Postcode, &curr.Country, &curr.Tel,
			&curr.Website, &curr.Email)

		if err != nil {
			list = nil
			return
		}

		list = append(list, curr)
	}
	return
}

func PutUserList(userId, listId int64, places []int64) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}

	_, err = tx.Stmt(deleteUserListStmt).Exec(userId, listId)

	if err != nil {
		tx.Rollback()
		return
	}

	for i, item := range places {
		_, err = tx.Stmt(putUserListStmt).Exec(item, listId, userId, i+1)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	return

}

func DeleteUserList(userId, listId int64) (err error) {
	_, err = deleteUserListStmt.Exec(userId, listId)
	return
}

/* id, email, name, fname, lname, birthday, " +
"school, picture, gender, password_hash, role */

func GetUserFriends(userId int64) (followers,
	following, both []User, err error) {

	rows, err := getUserFriendsStmt.Query(userId)

	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return
	}

	followers = make([]User, 0, 10)
	following = make([]User, 0, 10)
	both = make([]User, 0, 10)

	for rows.Next() {
		var curr User
		var direction int
		err = rows.Scan(&direction, &curr.Id, &curr.Email, &curr.Name,
			&curr.Fname, &curr.Lname, &curr.Birthday, &curr.School,
			&curr.Picture, &curr.Gender)

		log.Print(direction)
		log.Print(curr)

		if direction == 1 {
			followers = append(followers, curr)
		} else if direction == 2 {
			following = append(following, curr)
		} else if direction == 3 {
			both = append(both, curr)
		}

		if err != nil {
			followers = nil
			following = nil
			both = nil
			return
		}
		//users = append(users, curr)
	}
	return
}

func AddUserFriend(fromId, toId int64) (err error) {
	_, err = addUserFriendStmt.Exec(fromId, toId)

	err = AddNewsfeedItem(&FeedItem{UserId: fromId, Type: NfFriendAdded, Picture: "<Picture>", AuxInt: toId, AuxString: "<Name>"})
	return
}

func DeleteUserFriend(fromId, toId int64) (err error) {
	_, err = deleteUserFriendStmt.Exec(fromId, toId)
	return
}

func NetworkUserFriends(userId int64) (users []User, err error) {

	rows, err := networkUserFriendsStmt.Query(userId)

	if err != nil {
		return
	}

	users = make([]User, 0, 20)

	for rows.Next() {
		var curr User

		err = rows.Scan(&curr.Id, &curr.Email, &curr.Name,
			&curr.Fname, &curr.Lname, &curr.Birthday, &curr.School,
			&curr.Picture, &curr.Gender)

		if err != nil {
			rows = nil
			return
		}

		users = append(users, curr)
	}

	return
}

func ContactsUserFriends(userId int64, contacts []string) (users []User, err error) {

	contactsStr := strings.Join(contacts, "&")
	_, err = contactsUserFriendsStmt.Query(userId, contactsStr)

	if err != nil {
		return
	}

	users = make([]User, 0, 20)
	return
}

func SuggestUserFriends(userId int64, contacts []string) (users []User, err error) {

	contactsStr := strings.Join(contacts, "&")
	_, err = suggestUserFriendsStmt.Query(userId, contactsStr)

	if err != nil {
		return
	}

	users = make([]User, 0, 20)
	return
}

func GetListTypes() (lists []List, err error) {
	lists = make([]List, 0, 24)
	rows, err := getListTypesStmt.Query()

	if err != nil {
		lists = nil
		return
	}

	for rows.Next() {
		var curr List
		var children string
		err = rows.Scan(&curr.Id, &curr.Name, &curr.Icon, &children)

		if children != "" {
			childrenSlice := strings.Split(children, ",")
			curr.Children = make([]int64, len(childrenSlice))
			for key, child := range childrenSlice {
				curr.Children[key], err = strconv.ParseInt(
					strings.TrimSpace(child), 10, 64)
				if err != nil {
					return nil, err
				}
			}
		}

		lists = append(lists, curr)
		if err != nil {
			lists = nil
			return
		}
	}
	return
}

func GetWhooplistCoordinate(userId, listId int64, page int32, lat, long,
	radius float64) (places []Place, err error) {

	places = make([]Place, 0, 20)

	rows, err := getWhooplistCoordinateStmt.Query(
		listId, lat, long, radius, (10 * (page - 1)))

	if err != nil {
		return
	}

	for rows.Next() {
		var place Place

		err = rows.Scan(&place.Score, &place.Id, &place.Latitude,
			&place.Longitude, &place.FactualId, &place.Name, &place.Address,
			&place.Locality, &place.Region, &place.Postcode, &place.Country,
			&place.Tel, &place.Website, &place.Email)

		places = append(places, place)

		if err != nil {
			return
		}
	}
	return
}

func AddNewsfeedItem(item *FeedItem) (err error) {
	_, err = addNewsfeedItemStmt.Query(&item.UserId, &item.Latitude,
		&item.Longitude, &item.PlaceId, &item.ListId, &item.Picture,
		&item.Type, &item.AuxString, &item.AuxInt)
	return
}

func GetNewsfeed(user_id, latest_id int64,
	lat, long, radius float64) (items []FeedItem, err error) {

	return getNewsfeed(getNewsfeedStmt.Query(
		user_id, latest_id, lat, long, radius))
}

func GetNewsfeedEarlier(user_id, earliest_id int64,
	lat, long, radius float64) (items []FeedItem, err error) {

	return getNewsfeed(getNewsfeedEarlierStmt.Query(
		user_id, earliest_id, lat, long, radius))
}

func getNewsfeed(rows *sql.Rows, inErr error) (items []FeedItem, err error) {
	if inErr != nil {
		return nil, inErr
	}

	items = make([]FeedItem, 0, 30)

	for rows.Next() {
		var item FeedItem
		err = rows.Scan(&item.Timestamp, &item.UserId, &item.Latitude,
			&item.Longitude, &item.PlaceId, &item.ListId, &item.Picture,
			&item.Type, &item.AuxString, &item.AuxInt)
		items = append(items, item)
		if err != nil {
			items = nil
			return
		}
	}
	return
}

func GetPlace(placeId int64) (place *Place, err error) {
	place = new(Place)

	res := getPlaceStmt.QueryRow(placeId)
	err = res.Scan(&place.Id, &place.Latitude, &place.Longitude,
		&place.FactualId, &place.Name, &place.Address, &place.Locality,
		&place.Region, &place.Postcode, &place.Country,
		&place.Tel, &place.Website, &place.Email)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		place = nil
	}

	return
}

func SearchPlace(str string, listId int64, page int32,
	lat, long, radius float64) (places []Place, err error) {

	places, err = factualSearchPlace(str, lat, long, radius, page)

	if err != nil {
		places = nil
		return
	}

	err = addPlaces(places)

	return
}

func addPlaces(places []Place) (err error) {
	for k, place := range places {
		res := addPlaceStmt.QueryRow(place.Latitude, place.Longitude,
			place.FactualId, place.Name, place.Address, place.Locality,
			place.Region, place.Postcode, place.Country, place.Tel,
			place.Website, place.Email)

		err = res.Scan(&(places[k].Id))

		if err == sql.ErrNoRows {
			res := updatePlaceStmt.QueryRow(place.Latitude, place.Longitude,
				place.FactualId, place.Name, place.Address, place.Locality,
				place.Region, place.Postcode, place.Country, place.Tel,
				place.Website, place.Email)
			err = res.Scan(&(places[k].Id))
			if err != nil {
				return
			}
		}
	}
	return
}

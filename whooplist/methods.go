package whooplist

import (
	"strings"
	"io"
	"log"
	"crypto/rand"
	_ "github.com/lib/pq"
	"database/sql"
	//"time"
	"encoding/base64"
	"code.google.com/p/go.crypto/scrypt"
)

/* Application global database connection pool */
var db *sql.DB

var getUserDataStmt, createUserStmt, updateUserStmt, deleteUserStmt,
	authUserStmt, loginUserVerifyStmt, loginUserStmt,
	deleteSessionStmt, getUserListsStmt, getUserListStmt, getUserListPlaceStmt,
	putUserListStmt, deleteUserListStmt, getListTypesStmt,
	getWhooplistCoordinateStmt, getWhooplistLocationStmt, getLocationStmt,
	getLocationCoordinateStmt, getPlaceStmt *sql.Stmt

func prepare() (err error) {
	getUserDataStmt, err = db.Prepare(
		"SELECT id, email, name, birthday, school, picture, gender, role " +
		"FROM public.user WHERE id = $1;")
	if err != nil { return; }
	
	createUserStmt, err = db.Prepare(
		"INSERT INTO public.user (email, name, birthday, school, picture, " + 
		"gender, password_hash, role) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8);")
	if err != nil { return; }
	
	updateUserStmt, err = db.Prepare(
		"UPDATE public.user SET email = $1, name = $2, birthday = $3, school = $4, " +
		"picture = $5, gender = $6,  password_hash = $7, " +
		"role = $8 WHERE id = $9;")
	if err != nil { return; }

	deleteUserStmt, err = db.Prepare(
		"DELETE FROM public.user WHERE id = $1;")
	if err != nil { return; }	

	authUserStmt, err = db.Prepare(
		"SELECT * FROM public.session WHERE key = $1;")
	if err != nil { return; }
	
	loginUserVerifyStmt, err = db.Prepare(
		"SELECT id, email, name, birthday, school, picture, gender, role " +
		"FROM public.user WHERE email = $1 AND password_hash = $2;")
	if err != nil { return; }

	loginUserStmt, err = db.Prepare(
		"INSERT INTO public.session (user_id, key) " +
		"VALUES ($1, $2);")
	if err != nil { return; }
	
	deleteSessionStmt, err = db.Prepare(
		"DELETE FROM public.session WHERE key = $1;");
	if err != nil { return; }
/*	
	getUserListsStmt, err = db.Prepare(
		"SELECT list_id FROM list_item WHERE user_id = ? " +
		"GROUP BY list_id")
	if err != nil { return; }

	getUserListStmt, err = db.Prepare(
		"SELECT place_id FROM list_item " +
		"WHERE user_id = ? AND list_id = ? " +
		"ORDER BY rank")
	if err != nil { return; }	

	getUserListPlaceStmt, err = db.Prepare(
		"SELECT place.* FROM list_item JOIN place " +
		"ON list_item.place_id = place.id " +
		"WHERE list_item.user_id = ? AND list_item.list_id = ? " +
		"ORDER BY rank")
	if err != nil { return; }	

	putUserListStmt, err = db.Prepare(
		"INSERT INTO list_item (place_id, list_id, user_id, ranking) " +
		"VALUES (?, ?, ?, ?)")
	if err != nil { return; }
	
	deleteUserListStmt, err = db.Prepare(
		"DELETE FROM list_item WHERE user_id = ? AND list_id = ?")
	if err != nil { return; }
	
	getListTypesStmt, err = db.Prepare(
		"SELECT * FROM list")
	if err != nil { return; }
	
	//TODO: fix to coordinate
	getWhooplistCoordinateStmt, err = db.Prepare(
		"SELECT place.id AS place_id, SUM(10 - rank) AS score " +
                "FROM list_item " +
                "JOIN place ON list_item.place_id = place.id " +
                "JOIN place_location ON list_item.place_id = place_location.place_id " +
                "WHERE place_location.location_id = ? " +
                "WHERE list_item.list_id = ? " +
                "ORDER BY score " +
                "GROUP BY list_item.place_id " +
                "LIMIT 10 " +
                "OFFSET ? ")
	if err != nil { return; }
	
	//TODO: check correctness. probably incorrect
	getWhooplistLocationStmt, err = db.Prepare(
		"SELECT place.id AS place_id, SUM(10 - rank) AS score " +
		"FROM list_item " +
		"JOIN place ON list_item.place_id = place.id " +
		"JOIN place_location ON list_item.place_id = place_location.place_id " +
		"WHERE place_location.location_id = ? " +
		"WHERE list_item.list_id = ? " +
		"ORDER BY score " +
		"GROUP BY list_item.place_id " +
		"LIMIT 10 " +
		"OFFSET ? ")	
	if err != nil { return; }
	
	getLocationStmt, err = db.Prepare(
		"SELECT * FROM location WHERE id = ?")
	if err != nil { return; }

	//TODO: Figure out the math on this one properly	
	getLocationCoordinateStmt, err = db.Prepare(
		"SELECT * FROM location WHERE radius + 10 >  " + 
		"|/ ((centre_lat - ?) ^ 2 + (centre_long - ?))" + 
		"")
	if err != nil { return; }
	
	getPlaceStmt, err = db.Prepare(
		"SELECT * FROM place WHERE id = ?")
	if err != nil { return; }
*/
	return
}


func hash(username, password string) (hash string, err error) {
	secret := "aYdZYlE9ybGXn5CldvQ3f/shKxNshtAOvqDlaw/wbUBHwc5r9zBal9hf9CDkGxSgddAMtNm+uz1G"
	secretData, err := base64.StdEncoding.DecodeString(secret)

	if err != nil {
		return
	}	

	salt := make([]byte, len(secretData) + len(username))		

	copy(salt, secretData)
	copy(salt[len(secretData):], []byte(strings.ToLower(username)))	

	hash_data, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, 32) 	

	log.Print(username)
	log.Print(password)

	hash = base64.StdEncoding.EncodeToString(hash_data)	

	return
}

func Connect() (err error) {
	db, err = sql.Open("postgres",
		"user=whooplist dbname=whooplist password=moteifae0ohcaiCo sslmode=disable search_path=public,pg_catalog")

	if err == nil {
		err = prepare()
	}
	return
}

func Disconnect() (err error) {
	err = db.Close()
	return
}

func GetUserData(id int64) (user *User, err error) {
	res := getUserDataStmt.QueryRow(id)
	user = new(User)
	err = res.Scan(&user.Id, &user.Email, &user.Name, &user.Birthday, &user.School, &user.Picture, &user.Gender, &user.Role)

	if err == sql.ErrNoRows {
		user = nil
		err = nil
	} else if err != nil {
		user = nil
	}

	return
	
}

func CreateUser(user User) (err error) {	
	_, err = createUserStmt.Exec(user.Email, user.Name, user.Birthday, user.School, user.Picture, user.Gender, user.PasswordHash, user.Role)
	return
}

func UpdateUser(user User) (err error) {
	_, err = updateUserStmt.Exec(user.Email, user.Name, user.Birthday, user.School, user.Picture, user.Gender, user.PasswordHash, user.Role)
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
	err = res.Scan(&session.Id, &session.UserId, &session.Key, &session.LastAuth, &session.LastUse)

	if err == sql.ErrNoRows {
		session = nil
		err = nil
		return
	}	
	if err != nil {
		return
	}

	user, err = GetUserData(session.UserId)
	
	return
}

func createSession(userId int64) (session *Session, err error) {
	/* Generate a session key. Currently 192 bits of entropy. */
	data := make([]byte, 24)
	
	n, err := io.ReadFull(rand.Reader, data)


	if n != len(data) || err != nil {
		return
	}

	key := base64.StdEncoding.EncodeToString(data)

	_, err = loginUserStmt.Exec(userId, key, nil)

	if err != nil {
		return
	}

	/* TODO: We don't fill the whole session, we should identify if this is a problem. */

	session = new(Session)

	session.UserId = userId
	session.Key = key
	
	return
}

func LoginUser(username, password string) (user *User, session *Session, err error) {
	hash, err := hash(username, password)

	log.Print(hash)	

	if err != nil {
		return
	}

	res := loginUserVerifyStmt.QueryRow(username, hash)	
	user = new(User)
	err = res.Scan(&user.Id, &user.Email, &user.Name, &user.Birthday, &user.School, &user.Picture, &user.Gender, &user.Role)

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
	
        /* TODO: We don't fill the whole session, we should identify if this is a problem. */

	session = new(Session)
        session.UserId = user.Id
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

func GetUserLists(userId int) (lists []int, err error) {
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

func GetUserList(userId, listId int) (list *UserList, err error) {
	list = new(UserList)
	list.UserId = userId
	list.ListId = listId
	list.Items = make([]ListItem, 0, 5)

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
		var curr ListItem
		err = rows.Scan(&curr.Id, &curr.PlaceId, &curr.ListId, &curr.UserId, &curr.Ranking)
		list.Items = append(list.Items, curr)
		if err != nil {
			list = nil
			return
		}
	}
	return
}

func PutUserList(list UserList) (err error) {
	tx, err := db.Begin()

	if err != nil {
		return
	}

	_, err = tx.Stmt(deleteUserListStmt).Exec(list.UserId, list.ListId)

	if err != nil {
		tx.Rollback()
		return
	}

	for _, item := range list.Items {
		_, err = tx.Stmt(putUserListStmt).Exec(item.PlaceId, item.ListId, item.UserId, item.Ranking)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	
	tx.Commit()
	return

}

func DeleteUserList(userId, listId int) (err error) {
	_, err = deleteUserListStmt.Exec(userId, listId)
	return
}

func GetListTypes() (lists []List, err error) {
	lists = make([]List, 0, 20) //TODO: set capacity proper
	rows, err := getListTypesStmt.Query()

	if err != nil {
		lists = nil
		return
	}	

	for rows.Next() {
		var curr List
		err = rows.Scan(&curr.Id, &curr.Name, &curr.Icon)
		lists = append(lists, curr) 	
		if err != nil {
			lists = nil
			return
		}
	}
	return
}

func GetWhooplistCoordinate(userId int, listId int, page int, lat float64, long float64, radius float64) (list *WhooplistCoordinate, err error) {
	//TODO: Implement
	return	
}

func GetWhooplistLocation(userId int, listId int, page int, locationId int) (list *WhooplistLocation, err error) {
	//TODO: Implement
	return
}

func GetLocation(locationId int) (location *Location, err error) {
	location = new(Location)
	
	res := getLocationStmt.QueryRow(locationId)
	err = res.Scan(&location.Id, &location.Name)

	if err != nil {
		location = nil
	}

	return
}

func GetLocationCoordinate(lat, long float64) (locations []Location, err error) {
	//TODO: Implement
	return
}

func GetPlace(placeId int) (place *Place, err error) {
	place = new(Place)

	res := getPlaceStmt.QueryRow(placeId)
	err = res.Scan(&place.Id, &place.Latitude, &place.Longitude,
		&place.TomTomId, &place.Name,
		&place.Address, &place.Phone)

	if err != nil {
		place = nil
	}

	return
}

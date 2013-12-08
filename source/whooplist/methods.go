package whooplist

import (
	"code.google.com/p/go.crypto/scrypt"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	_ "github.com/lib/pq"
	"io"
	"strconv"
	"strings"
)

/* Application global database connection pool */
var db *sql.DB

var getUserDataStmt, createUserStmt, updateUserStmt, deleteUserStmt,
	authUserStmt, loginUserVerifyStmt, loginUserStmt,
	deleteSessionStmt, existsUserStmt, getUserListsStmt, getUserListStmt,
	getUserListPlaceStmt, putUserListStmt, deleteUserListStmt,
	getListTypesStmt, getWhooplistCoordinateStmt, getWhooplistLocationStmt,
	addNewsfeedItemStmt, getNewsfeedStmt, getNewsfeedEarlierStmt,
	getLocationStmt, getLocationCoordinateStmt, getPlaceStmt,
	getPlaceByFactualStmt, addPlaceStmt, updatePlaceStmt *sql.Stmt

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
			"FROM wl.user WHERE email = $1 AND password_hash = $2;")
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
		"SELECT place_id FROM wl.list_item " +
			"WHERE user_id = $1 AND list_id = $2 " +
			"ORDER BY rank")
	if err != nil {
		return
	}

	getUserListPlaceStmt, err = db.Prepare(
		"SELECT place.id, place.latitude, place.longitude, place.factual_id, place.name, " +
			"place.address, place.locality, place.region, place.postcode, place.country, " +
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

	getListTypesStmt, err = db.Prepare(
		"SELECT * FROM wl.list")
	if err != nil {
		return
	}
	/*
	   	getWhooplistStmt, err = db.Prepare(
	   		"SELECT score, place.id, place.latitude, place.longitude, place.factual_id, place.name, " +
	                           "place.address, place.locality, place.region, place.postcode, place.country, " +
	                           "place.telephone, place.website, place.email " +
	   			"FROM whooplist_item JOIN place " +
	   			"ON whooplist_item.place_id = place.id " +
	   			"WHERE whooplist_item.list_id = $1 AND " +
	   			"((place.lat - $2)^2 + (place.long - $3)^2) < ($4 * $4)
	   	);
	*/
	/*	   	//TODO: fix to coordinate
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
	*/

	/*addNewsfeedItemStmt, err = db.Prepare(
			"INSERT INTO newsfeed_item (user_id, location_id, place_id, list_id " +
			"timestamp, picture, is_new_in_list, position, is_trending, is_visiting, " +
			"profile_picture_updated, school_updated, school) VALUES ($1, $2, $3, $4 " +
			"$5, $6, $7, $8, $9, $10, $11, $12, $13)")
		if err != nil {
			return
		}

		getNewsfeedStmt, err = db.Prepare(
			"SELECT user_id, location_id, place_id, list_id, timestamp, picture, is_new_in_list " +
			"position, is_trending, is_visiting, profile_picture_updated FROM newsfeed_item " +
			"WHERE (location_id = $1 OR user_id = $3) AND latest_id > $2 LIMIT 30")
		if err != nil {
			return
		}

		getNewsfeedEarlierStmt, err = db.Prepare(
			"SELECT user_id, location_id, place_id, list_id, timestamp, picture, is_new_in_list " +
	                "position, is_trending, is_visiting, profile_picture_updated FROM newsfeed_item " +
		        "WHERE (location_id = $1 OR user_id = $3) AND latest_id < $2 LIMIT 30")
		if err != nil {
			return
		}*/

	/*
		getLocationStmt, err = db.Prepare(
			"SELECT * FROM location WHERE id = ?")
		if err != nil { return; }

		//TODO: Figure out the math on this one properly
		getLocationCoordinateStmt, err = db.Prepare(
			"SELECT * FROM location WHERE radius + 10 >  " +
			"|/ ((centre_lat - ?) ^ 2 + (centre_long - ?))" +
			"")
		if err != nil { return; }
	*/
	getPlaceStmt, err = db.Prepare(
		"SELECT latitude, longitude, factual_id, name, address, locality, " +
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
	if user.Picture != nil {
		str, err := WriteFileBase64("profile.jpg", user.Picture, false)
		if err != nil {
			return err
		}
		user.Picture = &str
	}
	_, err = createUserStmt.Exec(user.Email, user.Name, user.Fname,
		user.Lname, user.Birthday, user.School, user.Picture,
		user.Gender, user.PasswordHash, user.Role)
	return
}

func CheckUpdateUser(email, password string) (user *User, err error) {
	hash, err := Hash(email, password)

	if err != nil {
		return
	}

	user = new(User)
	res := loginUserVerifyStmt.QueryRow(email, hash)
	err = res.Scan(&user.Id, &user.Email, &user.Name, &user.Fname,
		&user.Lname, &user.Birthday, &user.School, &user.Picture,
		&user.Gender, &user.Role)
	if err == sql.ErrNoRows {
		user = nil
		err = nil
		return
	}
	if err != nil {
		user = nil
		return
	}
	return
}

func UpdateUser(user User) (err error) {
	var userHash string
	if user.Password != "" {
		userHash, _ = Hash(user.Email, user.Password)
	} else {
		userHash = user.PasswordHash
	}

	if user.Picture != nil && !strings.HasPrefix(*user.Picture, "static.whooplist.com") {
		str, err := WriteFileBase64("profile.jpg", user.Picture, false)
		if err != nil {
			return err
		}
		user.Picture = &str
	}
	_, err = updateUserStmt.Exec(user.Email, user.Name, user.Fname,
		user.Lname, user.Birthday, user.School, user.Picture,
		user.Gender, userHash, user.Role)
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
	err = res.Scan(&session.Id, &session.UserId, &session.Key)

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

func GetUserList(userId, listId int64) (list *UserList, err error) {
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

	rank := 1
	for rows.Next() {
		var curr ListItem
		err = rows.Scan(&curr.PlaceId)
		curr.Rank = rank
		rank += 1

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
		_, err = tx.Stmt(putUserListStmt).Exec(item.PlaceId, item.ListId,
			item.UserId, item.Rank)
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

func GetWhooplistCoordinate(userId, listId int64, page int, lat, long,
	radius float64) (list *WhooplistCoordinate, err error) {
	//TODO: Implement
	return
}

func GetWhooplistLocation(userId, listId int64, page int,
	locationId int) (list *WhooplistLocation, err error) {
	//TODO: Implement
	return
}

func AddNewsfeedItem(item *FeedItem) (err error) {
	_, err = addNewsfeedItemStmt.Query(&item.UserId, &item.LocationId,
		&item.PlaceId, &item.ListId, &item.Picture,
		&item.Type, &item.AuxString, &item.AuxInt)
	return
}

func GetNewsfeed(location, latest_id, user_id int64) (items []FeedItem, err error) {
	return getNewsfeed(getNewsfeedStmt.Query(location,
		latest_id, user_id))
}

func GetNewsfeedEarlier(location, earliest_id, user_id int64) (items []FeedItem, err error) {
	return getNewsfeed(getNewsfeedEarlierStmt.Query(location,
		earliest_id, user_id))
}

func getNewsfeed(rows *sql.Rows, inErr error) (items []FeedItem, err error) {
	if inErr != nil {
		return
	}

	items = make([]FeedItem, 0, 30)

	for rows.Next() {
		var item FeedItem
		err = rows.Scan(&item.Timestamp, &item.UserId, &item.LocationId,
			&item.PlaceId, &item.ListId, &item.Picture, &item.Type,
			&item.AuxString, &item.AuxInt)
		items = append(items, item)
		if err != nil {
			items = nil
			return
		}
	}
	return
}

/*func GetLocation(locationId int) (location *Location, err error) {
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
}*/

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
		return
	}

	/*for _, place := range places {
		err = storePlace(&place)
		if err != nil {
			places = nil
			return
		}
	}*/
	return
}

func storePlace(place *Place) (err error) {
	res := addPlaceStmt.QueryRow(place.Latitude, place.Longitude,
		place.FactualId, place.Name, place.Address, place.Locality,
		place.Region, place.Postcode, place.Country, place.Tel,
		place.Website, place.Email)

	err = res.Scan(&place.Id)
	return
}

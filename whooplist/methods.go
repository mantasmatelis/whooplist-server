package whooplist

import (
	_ "github.com/lib/pq"
	"database/sql"
	"time"
)

/* Application global database connection pool */
var db *sql.DB

var getUserDataStmt, createUserStmt, updateUserStmt, authUserStmt
	loginUserInternalStmt, loginUserFacebookStmt, loginUserGoogleStmt,
	deleteSessionStmt, getUserListsStmt, getUserListStmt, getUserListPlace,
	putUserListStmt, deleteUserListStmt, getListTypesStmt,
	getWhooplistCoordinateStmt, getWhooplistLocationStmt, getLocationStmt,
	getLocationCoordinateStmt, getPlaceStmt *sql.Stmt

func prepare() (err error) {
	getUserDataStmt, err = db.Prepare(
		"SELECT * FROM user WHERE id=?")
	if err != nil { return; }
	
	createUserStmt, err = db.Prepare(
		"INSERT INTO user (id, email, name, password_hash, " + 
		"facebook_id, gmail_id, role) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil { return; }
	
	updateUserStmt, err = db.Prepare(
		"UPDATE user SET email = ?, name = ?, password_hash = ?, " +
		"facebook_id = ?, gmail_id = ?, role = ? WHERE id = ?")
	if err != nil { return; }
	
	authUserStmt, err = db.Prepare(
		"SELECT * FROM session WHERE key = ?")
	if err != nil { return; }
	
	loginUserInternalCheckStmt, err = db.Prepare(
		"SELECT * FROM user WHERE username = ? AND password = ?")
	if err != nil { return; }
	
	loginUserInternalCreateStmt, err = db.Prepare(
		"INSERT INTO session (user_id, key) VALUES (?, ?)")
	if err != nil { return; }	

	loginUserFacebookStmt, err = db.Prepare(
		"INSERT INTO session (user_id, key, facebook_token) VALUES (?, ?, ?)");
	if err != nil { return; }
	
	loginUserGoogleStmt, err = db.Prepare(
		"INSERT INTO session (user_id, key, google_token) VALUES (?, ?, ?)");
	if err != nil { return; }
	
	deleteSessionStmt, err = db.Prepare(
		"DELETE FROM session WHERE key = ? LIMIT 1");
	if err != nil { return; }
	
	getUserListsStmt, err = db.Prepare(
		"SELECT list_id FROM list_item WHERE user_id = ?" +
		"GROUP BY list_id")
	if err != nil { return; }

	getUserListStmt, err = db.Prepare(
		"SELECT place_id FROM list_item " +
		"WHERE user_id = ? AND list_id = ? " +
		"ORDER BY rank")
	)
	if err != nil { return; }	

	getUserListPlaceStmt, err = db.Prepare(
		"SELECT place.* FROM list_item JOIN place " +
		"ON list_item.place_id = place.id " +
		"WHERE list_item.user_id = ? AND list_item.list_id = ? " +
		"ORDER BY rank")
	if err != nil { return; }	

	putUserListStmt, err = db.Prepare(
		"INSERT INTO list_item (place_id, list_id, user_id, ranking) VALUES (?, ?, ?, ?)")
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
}

func Connect() (err error) {
	db, err := sql.Open("postgres",
		"user=whooplist dbname=whooplist password=whooplist sslmode=disable")

	if err == nil {
		err = prepare()
	}
}

func Disconnect() (err error) {
	err := db.Close()
}

func GetUserData(id int) (user User, err error) {
}

func CreateUser(user User) (err error) {	
}

func UpdateUser(user User) (err error) {
}

func AuthUser(key string) (user User, session Session, err error) {
}

func LoginUserInternal(username, password string) (user User, session Session, err error) {
}

func LoginUserFacebook(authToken string) (user User, session Session, err error) {
}

func LoginUserGoogle(authToken string) (user User, session Session, err error) {
}

func DeleteSession(key string) (err error) {
}

func GetUserLists(userId int) (lists []List, err error) {
}

func GetUserList(userId, listId int) (list List, err error) {
}

func PutUserList(userId int, list List) (err error) {
}

func DeleteUserList(userId, listId int) (err error) {
}

func GetListTypes() (err error) {
}

func GetWhooplistCoordinate(userId int, listId int, page int, lat float, long float, radius float) (list List, err error) {
}

func GetWhooplistLocation(userId int, listId int, page int, locationId int) (list List, err error) {
}

func GetLocation(locationId int) (location Location, err error) {
}

func GetLocationCoordinate(lat, long float) (locations []Location, err error) {
}

func GetPlace(placeId int) (place Place, err error) {
}

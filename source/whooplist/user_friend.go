package whooplist

import (
	"database/sql"
	"log"
	"strings"
)

var networkUserFriendsStmt, suggestUserFriendsStmt,
	contactsUserFriendsStmt, getUserFriendsStmt,
	addUserFriendStmt, deleteUserFriendStmt *sql.Stmt

func prepareUserFriend() {
	stmt(&getUserFriendsStmt,
		"SELECT SUM(direction), id, email, name, fname, lname, birthday, "+
			"school, picture, gender FROM "+
			"(SELECT 1 AS direction, wl.user.id AS id, email, name, fname, "+
			"lname, birthday, school, picture, gender, role FROM wl.user "+
			"JOIN wl.friend ON wl.user.id = friend.from_id "+
			"WHERE friend.to_id = $1 UNION "+
			"SELECT 2 AS direction, wl.user.id, email, name, fname, lname, "+
			"birthday, school, picture, gender, role FROM wl.user "+
			"JOIN wl.friend ON wl.user.id = friend.to_id "+
			"WHERE friend.from_id = $1) AS bothDirections GROUP BY id, "+
			"email, name, fname, lname, birthday, school, picture, gender, role;")

	stmt(&addUserFriendStmt,
		"INSERT INTO wl.friend (from_id, to_id) VALUES ($1, $2);")

	stmt(&deleteUserFriendStmt,
		"DELETE FROM wl.friend WHERE from_id = $1 AND to_id = $2;")

	stmt(&networkUserFriendsStmt,
		"SELECT COUNT(*), wl.user.id, email, name, fname, lname, birthday, "+
			"school, picture, gender FROM wl.friend "+
			"JOIN wl.user ON wl.user.id = to_id "+
			"WHERE from_id IN "+
			"(SELECT DISTINCT to_id FROM friend WHERE from_id = $1) "+
			"AND wl.user.id NOT IN "+
			"(SELECT DISTINCT to_id FROM friend WHERE from_id = $1) "+
			"AND wl.user.id <> $1 "+
			"GROUP BY wl.user.id "+
			"ORDER BY count DESC "+
			"LIMIT 10;")

	stmt(&suggestUserFriendsStmt,
		"SELECT COUNT(*), wl.user.id, email, name, fname, lname, birthday, "+
			"school, picture, gender FROM wl.friend "+
			"JOIN wl.user ON wl.user.id = to_id "+
			"WHERE from_id IN "+
			"(SELECT DISTINCT to_id FROM friend WHERE from_id = $1) "+
			"AND wl.user.id NOT IN "+
			"(SELECT DISTINCT to_id FROM friend WHERE from_id = $1) "+
			"AND wl.user.id <> $1 "+
			"GROUP BY wl.user.id "+
			"ORDER BY count DESC "+
			"LIMIT 10;")

	stmt(&contactsUserFriendsStmt,
		"SELECT id, email, name, fname, lname, "+
			"birthday, school, picture, gender FROM wl.user "+
			"WHERE email IN (SELECT unnest (string_to_array($1, '&'))) "+
			"OR phone IN (SELECT unnest (string_to_array($2, '&'))) "+
			"GROUP BY wl.user.id;")

}

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

	//TODO: Figure out how to add Picture and AuxStiring as name
	err = AddNewsfeedItem(
		&FeedItem{UserId: fromId, Type: NfFriendAdded, AuxInt: toId})
	if err != nil {
		return
	}
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

		err = rows.Scan(&curr.Score, &curr.Id, &curr.Email, &curr.Name,
			&curr.Fname, &curr.Lname, &curr.Birthday, &curr.School,
			&curr.Picture, &curr.Gender)

		if err != nil {
			users = nil
			return
		}

		users = append(users, curr)
	}

	return
}

func ContactsUserFriends(userId int64, contacts []string) (users []User, err error) {

	contactsStr := strings.Join(contacts, "&")
	rows, err := contactsUserFriendsStmt.Query(userId, contactsStr, contactsStr)

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
			users = nil
			return
		}

		users = append(users, curr)
	}

	return
}

func SuggestUserFriends(userId int64, contacts []string) (users []User, err error) {

	contactsStr := strings.Join(contacts, "&")
	rows, err := suggestUserFriendsStmt.Query(userId, contactsStr, contactsStr)

	if err != nil {
		return
	}

	users = make([]User, 0, 20)

	for rows.Next() {
		var curr User

		err = rows.Scan(&curr.Score, &curr.Id, &curr.Email, &curr.Name,
			&curr.Fname, &curr.Lname, &curr.Birthday, &curr.School,
			&curr.Picture, &curr.Gender)

		if err != nil {
			users = nil
			return
		}

		users = append(users, curr)
	}
	return
}

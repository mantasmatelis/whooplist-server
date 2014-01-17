package whooplist

import (
	"database/sql"
	"strconv"
	"strings"
)

type List struct {
	Id       int64
	Name     string
	Icon     string
	Children []int64
}

var getUserListsStmt, getUserListStmt, getUserListPlaceStmt,
	putUserListStmt, deleteUserListStmt, getListTypesStmt,
	getWhooplistCoordinateStmt *sql.Stmt

func prepareList() {
	stmt(&getUserListsStmt,
		"SELECT list_id FROM wl.list_item WHERE user_id = $1 "+
			"GROUP BY list_id;")

	stmt(&getUserListStmt,
		"SELECT place.id, place.latitude, place.longitude, "+
			"place.factual_id, place.name, place.address, "+
			"place.locality, place.region, place.postcode, place.country, "+
			"place.telephone, place.website, place.email "+
			"FROM wl.list_item JOIN wl.place "+
			"ON list_item.place_id = place.id "+
			"WHERE list_item.user_id = $1 AND list_item.list_id = $2 "+
			"ORDER BY rank")

	stmt(&putUserListStmt,
		"INSERT INTO wl.list_item (place_id, list_id, user_id, rank) "+
			"VALUES ($1, $2, $3, $4)")

	stmt(&deleteUserListStmt,
		"DELETE FROM wl.list_item WHERE user_id = $1 AND list_id = $2")

	stmt(&getListTypesStmt,
		"SELECT * FROM wl.list")

	/* Note, the constant below is 360 / (earth's radius in meters) */
	stmt(&getWhooplistCoordinateStmt,
		"SELECT SUM(10 - rank) AS score, place.id, place.latitude, place.longitude, "+
			"place.factual_id, place.name, place.address, "+
			"place.locality, place.region, place.postcode, place.country, "+
			"place.telephone, place.website, place.email "+
			"FROM list_item JOIN place "+
			"ON list_item.place_id = place.id "+
			"WHERE list_item.list_id = $1 AND "+
			"latitude + (($4 * 0.00000898314) / cos(latitude)) > $2 AND "+
			"latitude - (($4 * 0.00000898314) / cos(latitude)) < $2 AND "+
			"longitude + ($4 * 0.00000898314) > $3 AND "+
			"longitude - ($4 * 0.00000898314) < $3 "+
			"GROUP BY place.id "+
			"LIMIT 10 OFFSET $5")
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

	//TODO: Figure out how to add picture
	AddNewsfeedItem(
		&FeedItem{Type: NfNewInUserList, UserId: userId, ListId: listId})

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

package whooplist

import (
	"database/sql"
	"time"
)

type FeedItem struct {
	/* General Items */
	Id        int64
	Timestamp time.Time

	/* The following group of items do not all
	   have to be present and may be 0 */
	UserId    int64
	Latitude  float64
	Longitude float64
	PlaceId   int64
	ListId    int64
	Picture   string

	/* Type Describes what the Newsfeed event is
	* (No UserId) 1: Trending, 2: NewInWhooplist (AuxInt: Position)
	*    (UserId) 3: Welcome, 4: WhyNotAdd,
	          ... 5: UpdatedUserlist,
	*         ... 6: Visiting 7: ProfilePictureUpdated,
	*         ... 8: SchoolUpdated (AuxString: Name),
	*         ... 9: FriendAdded (AuxInt: Id, AuxString: Name)
	*/
	Type      int64
	AuxString string
	AuxInt    int64
}

//Add trending, newinwhooplist, WHYNOTADD, visiting, PROFILEPICTUREUPDATED,
// SCHOOLUPDATED
const (
	NfTrending = iota
	NfNewInWhooplist
	NfWelcome
	NfWhyNotAdd
	NfNewInUserList
	NfVisiting
	NfProfilePictureUpdated
	NfSchoolUpdated
	NfFriendAdded
)

var addNewsfeedItemStmt, getNewsfeedStmt,
	getNewsfeedEarlierStmt *sql.Stmt

func prepareNewsfeed() {
	stmt(&addNewsfeedItemStmt,
		"INSERT INTO wl.feed_item (user_id, latitude, longitude, "+
			"place_id, list_id, picture, type, aux_string, aux_int) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)")

	stmt(&getNewsfeedStmt,
		"SELECT id, timestamp, user_id, latitude, longitude, place_id, "+
			"list_id, picture, type, aux_string, aux_int "+
			"FROM wl.feed_item "+
			"WHERE ("+
			"(user_id = $1) "+
			"OR "+
			"(latitude + (($5 * 0.00000898314) / cos(latitude)) > $3 AND "+
			"latitude - (($5 * 0.00000898314) / cos(latitude)) < $3 AND "+
			"longitude + ($5 * 0.00000898314) > $4 AND "+
			"longitude - ($5 * 0.00000898314) < $4)) "+
			"AND (id > $2 OR $2 = -1) "+
			"LIMIT 30")

	stmt(&getNewsfeedEarlierStmt,
		"SELECT id, timestamp, user_id, latitude, longitude, place_id, "+
			"list_id, picture, type, aux_string, aux_int "+
			"FROM wl.feed_item "+
			"WHERE ("+
			"(user_id = $1) "+
			"OR "+
			"(latitude + (($5 * 0.00000898314) / cos(latitude)) > $3 AND "+
			"latitude - (($5 * 0.00000898314) / cos(latitude)) < $3 AND "+
			"longitude + ($5 * 0.00000898314) > $4 AND "+
			"longitude - ($5 * 0.00000898314) < $4)) "+
			"AND id < $2 LIMIT 30")
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

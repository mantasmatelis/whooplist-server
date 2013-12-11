package whooplist

import (
	"time"
)

type User struct {
	Id       int64
	Email    string
	Name     string
	Fname    *string
	Lname    *string
	Birthday *time.Time `json:",omitempty"`
	School   *string    `json:",omitempty"`
	Picture  *string    `json:",omitempty"`
	Gender   *int64     `json:",omitempty"` /* Follows ISO/IEC 5218 */

	/* Private (i.e. not passed to client ever) */
	PasswordHash string `json:"-"`
	Role         byte   `json:"-"`

	/* For the client to pass to user only */
	Password    string `json:",omitempty"`
	OldPassword string `json:",omitempty"`
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

/**/
type Place struct {
	Id int64

	Latitude  float64
	Longitude float64

	FactualId string

	Name     string
	Address  string
	Locality string
	Region   string
	Postcode string
	Country  string
	Tel      string
	Website  string
	Email    string

	Score int64
}

type Location struct {
	Id   int
	Name string
}

type List struct {
	Id       int64
	Name     string
	Icon     string
	Children []int64
}

type ListItem struct {
	Id      int64 `json:"-"`
	PlaceId int64
	ListId  int64 `json:",omitempty"`
	UserId  int64 `json:",omitempty"`
	Rank    int
}

type FeedItem struct {
	/* General Items */
	Id        int
	Timestamp time.Time

	/* The following group of items do not all
	   have to be present and may be 0 */
	UserId     int
	LocationId int
	PlaceId    int
	ListId     int
	Picture    string

	/* Type Describes what the Newsfeed event is
	* (No UserId) 1: Trending, 2: NewInWhooplist (AuxInt: Position)
	*    (UserId) 3: Welcome, 4: WhyNotAdd,
	          ... 5: NewInUserlist (AuxInt: Position),
	*         ... 6: Visiting 7: ProfilePictureUpdated,
	*         ... 8: School Updated (AuxString: Name),
	*         ... 9: FriendAdded (AuxInt: Id, AuxString: Name)
	*/
	Type      int64
	AuxString string
	AuxInt    int64
}

/* The following four structures are not database tables
   and are for convenience purposes only */

type UserList struct {
	UserId int64
	ListId int64
	Items  []ListItem
}

type Whooplist struct {
	ListId int64
	Items  []ListItem
}

type WhooplistCoordinate struct {
	Whooplist
	Latitude  int
	Longitude int
	Radius    int
}

type WhooplistLocation struct {
	Whooplist
	LocationId int
}

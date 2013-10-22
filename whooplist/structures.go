package whooplist

import (
	"time"
)

type User struct {
	Id       int64
	Email    string
	Name     string
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

/* Sessions offer long-term authenticated communication between server and client.

Identity is proved with a long key (that must be kept absolutely secret at all times).

Sessions are signed out after 28 days of inactivity (not making requests under the key) for security.

*/
type Session struct {
	Id       int64
	UserId   int64
	Key      string
	LastAuth time.Time
	LastUse  time.Time
}

/**/
type Place struct {
	Id        int64
	Latitude  float64
	Longitude float64
	TomTomId  string
	Name      string
	Address   string
	Phone     string
}

type Location struct {
	Id   int
	Name string
}

type List struct {
	Id   int
	Name string
	Icon *byte
}

type ListItem struct {
	Id      int `json:"-"`
	PlaceId int
	ListId  int `json:",omitempty"`
	UserId  int `json:",omitempty"`
	Rank    int
}

type FeedItem struct {
	Id         int
	UserId     int /* User-specific (e.g. your friend Jitesh id 3) or 0 for relevant to all */
	LocationId int
	PlaceId    int
	ListId     int
	Timestamp  time.Time
	Picture    string

	/* Has only one of the following groups */
	IsNewInList bool /* Can be friend-specific or afriend-specific */
	Position    int

	IsTrending bool /* Afriend-specific */

	IsVisiting bool /* Friend-specific */
}

/* The following four structures are not database tables and are for convenience purposes only */

type UserList struct {
	UserId int
	ListId int
	Items  []ListItem
}

type Whooplist struct {
	ListId int
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

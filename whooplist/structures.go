package whooplist

import (
//	_ "github.com/bmizerany/pq"
//	"database/sql"
	"time"
)

type User struct {
	Id int
	Email string
	Name text
	PasswordHash text
	FacebookId text
	GmailId text
	Role char
}

type Session struct {
	Id int
	UserId int
	Key string
	LastAuth time.Time
	LastUse time.Time
}

type Place struct {
	Id int
	Latitude float
	Longitude float
	GoogleId string
	GoogleReference string
	Name string
	Address string
	Phone string
}

type Location struct {
	Id int
	Name text
}

type List struct {
	Id int
	Name text
	Icon *byte
}

type ListItem struct {
	Id int
	PlaceId int
	ListId int
	UserId int
	Ranking int
}

/* The following four structures are not database tables and are for convenience purposes only */

type UserList struct {
	UserId int
	ListId int
	Items []ListItem
}

type Whooplist struct {
	ListId int
	Items []ListItem
}

type WhooplistCoordinate struct {
	Whooplist
	Latitude int
	Longitude int
	Radius int	
}

type WhooplistLocation struct {
	Whooplist
	LocationId int
}

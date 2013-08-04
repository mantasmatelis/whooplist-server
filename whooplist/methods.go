package whooplist

import (
//	_ "github.com/bmizerany/pq"
//	"database/sql"
//	"time"
)

func RetrieveUserData(id int) (user User, err error) {
}

func StoreUser(user User) (err error) {
}

func AuthUser(key string) (user User, session Session, err error) {
}

func LoginUserInternal(username string, password string) (user User, session Session, err error) {
}

func LoginUserFacebook(authToken string) (user User, session Session, err error) {
}

func LoginUserGoogle(authToken string) (user User, session Session, err error) {
}

func DeleteSession(key string) (err error) {
}

func RetrieveUserLists(userId int) (lists []List, err error) {
}

func RetrieveUserList(userId int, listId int) (list List, err error) {
}

func StoreUserList(userId int, list List) (err error) {
}

func DeleteUserList(userId int, listId int) (err error) {
}

func RetrieveListTypes() (err error) {
}

func RetrieveWhooplistCoordinate(userId int, listId int, page int, lat float, long float, radius float) (list List, err error) {
}

func RetrieveWhooplistLocation(userId int, listId int, page int, locationId int) (list List, err error) {
}

func RetrieveLocation(locationId int) (location Location, err error) {
}

func RetrieveLocationCoordinate(lat, long float) (locations []Location, err error) {
}

func RetrievePlace(placeId int) (place Place, err error) {
}

package whooplist

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
)

var UserError = errors.New("whooplist: user error")
var WeakPassword = errors.New("whooplist: weak password error")
var BadPassword = errors.New("whooplist: bad password error")

/* Application global database connection pool */
var db *sql.DB

func stmt(s **sql.Stmt, query string) {
	retS, err := db.Prepare(query)
	if err != nil {
		log.Print("error while preparing statement: " + err.Error() +
			" \n relevent statement: " + query)
	}
	*s = retS
}

func Initialize() (err error) {
	db, err = sql.Open("postgres",
		"user=whooplist dbname=whooplist password=moteifae0ohcaiCo "+
			"sslmode=disable search_path=wl host=localhost")

	if err != nil {
		return err
	}

	factualInitialize()

	prepareUser()
	prepareUserFriend()
	prepareList()
	prepareNewsfeed()
	preparePlace()

	return
}

func Disconnect() (err error) {
	err = db.Close()
	return
}

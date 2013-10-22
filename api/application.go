package main

import (
	"io"
	"os"
	"log"
	"net/http"
	"source.whooplist.com/route"
	"source.whooplist.com/whooplist"
)

type MultiLog struct {
	Out []io.Writer
}

func (ml MultiLog) Write(p []byte) (n int, err error) {
	for _, o := range ml.Out {
		n, err = o.Write(p)
		if err != nil {
			panic("Could not log.")
		}
	}
	return len(p), nil
}

func main() {
	logFile, err := os.OpenFile("api.log", os.O_CREATE | os.O_RDWR | os.O_APPEND, 0660)

	if err != nil {
		panic("Could not create text log.")
	}

	log.SetOutput(&MultiLog{[]io.Writer{os.Stdout, logFile}})

	err = whooplist.Connect()
	if err != nil {
		log.Fatal("Could not connect to database/prepare statements, dying: " + err.Error())
	}

	var router route.Router
	router.SetRoutes(
		/* User Routes */
		route.Route{"POST", "/users/login", LoginUser},
		route.Route{"POST", "/users/logout", LogoutUser},
		route.Route{"POST", "/users/:UserId", UpdateUser},
		route.Route{"GET", "/users/:UserId", GetUser},
		route.Route{"PUT", "/users", CreateUser},

		/* User List Routes */
		route.Route{"GET", "/users/:UserId/lists", GetUserLists},
		route.Route{"GET", "/users/:UserId/lists/:ListId", GetUserList},
		route.Route{"POST", "/users/:UserId/lists/:ListId", CreateUserList},
		route.Route{"DELETE", "/users/:UserId/lists/:ListId", DeleteUserList},

		/* Possible List Routes */
		route.Route{"GET", "/listTypes", GetListTypes},

		/* Whooplist Routes */
		route.Route{"GET", "/whooplist/:ListId/:Page/coordinate/:Lat/:Long/:Radius", GetWhooplistCoordinate},
		route.Route{"GET", "/whooplist/:ListId/:Page/location/:LocationId", GetWhooplistLocation},

		/* Newsfeed Routes */
		route.Route{"GET", "/newsfeed/:Location/:LatestId", GetNewsfeed},
		route.Route{"GET", "/newsfeed/:Location/older/:EarliestId/", GetNewsfeedOlder},

		/* Location Routes */
		route.Route{"GET", "/locations/:LocationId", GetLocation},
		route.Route{"GET", "/locations/:Latitude/:Longitude", GetLocationsCoordinate},

		/* Place Routes */
		route.Route{"GET", "/places/:PlaceId", GetPlace},
	)

	s := &Server{Router: router}
	hs := &http.Server{
		Addr: ":3000",
		Handler: logHandler(panicHandler(http.HandlerFunc(s.handleRequest))),
	}
	log.Fatal(hs.ListenAndServe())
}

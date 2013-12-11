package main

import (
	"../whooplist"
	"github.com/mantasmatelis/go-trie-url-route"
	"log"
	"net/http"
)

func main() {
	/* Let caller handle timestamping */
	log.SetFlags(0)

	log.Print("initializing...")

	/* Initialze data layer */
	err := whooplist.Initialize()

	if err != nil {
		log.Fatal("could not initialize data layer, " +
			"dying: " + err.Error())
	}

	/* Set up routes */
	var router route.Router
	router.SetRoutes(
		route.Route{"GET", "/ping", Ping},

		/* User Routes */
		route.Route{"POST", "/users/login", LoginUser},
		route.Route{"POST", "/users/logout", LogoutUser},
		route.Route{"GET", "/user/exists/*Email", ExistsUser},
		route.Route{"POST", "/users", UpdateUser},
		route.Route{"GET", "/users/:UserId", GetUser},
		route.Route{"PUT", "/users", CreateUser},

		/* User List Routes */
		route.Route{"GET", "/users/:UserId/lists", GetUserLists},
		route.Route{"GET", "/users/:UserId/lists/:ListId", GetUserList},
		route.Route{"POST", "/users/lists/:ListId", CreateUserList},
		route.Route{"DELETE", "/users/lists/:ListId", DeleteUserList},

		/* User Friend Routes */
		//route.Route{"GET", "/friends", GetUserFriends},
		//route.Route{"PUT", "/friends/:OtherId", AddUserFriend},
		//route.Route{"DELETE", "/friends/:OtherId", DeleteUserFriend},

		/* Possible List Routes */
		route.Route{"GET", "/listTypes", GetListTypes},

		/* Whooplist Routes */
		route.Route{"GET",
			"/whooplist/:ListId/coordinate/:Lat/:Long/:Radius/:Page",
			GetWlCoordinate},
		//route.Route{
		//	"GET", "/whooplist/:ListId/location/:LocationId/:Page",
		//	GetWhooplistLocation},

		/* Newsfeed Routes */
		//route.Route{"GET", "/newsfeed/:Location/:LatestId", GetNewsfeed},
		//route.Route{"GET", "/newsfeed/:Location/older/:EarliestId/", GetNewsfeedOlder},

		/* Location Routes */
		//route.Route{"GET", "/locations/:LocationId", GetLocation},
		//route.Route{"GET", "/locations/:Latitude/:Longitude", GetLocationsCoordinate},

		/* Place Routes */
		route.Route{"GET",
			"/places/search/:ListId/:Lat/:Long/:Radius/:Page/*SearchString",
			SearchPlace},
		route.Route{"GET", "/places/:PlaceId", GetPlace},
	)

	/* Define the server, run it */
	s := &Server{Router: router}
	hs := &http.Server{
		Addr:    ":3000",
		Handler: logHandler(panicHandler(http.HandlerFunc(s.handleRequest))),
	}

	log.Print("done! listening.")
	log.Fatal(hs.ListenAndServe())
}

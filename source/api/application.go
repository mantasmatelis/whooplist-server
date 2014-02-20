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
		route.Route{"POST", "/users/lists/:ListId/append/:PlaceId",
			AppendUserList},

		/* User Friend Routes */
		route.Route{"GET", "/users/:UserId/friends", GetUserFriends},
		route.Route{"PUT", "/users/friends/:OtherId", AddUserFriend},
		route.Route{"DELETE", "/users/friends/:OtherId", DeleteUserFriend},

		/* User Friend Suggestion Routes */
		route.Route{"GET", "/users/friends/suggestions",
			SuggestUserFriends},
		route.Route{"GET", "/users/friends/suggestions/contacts",
			ContactsUserFriends},
		route.Route{"GET", "/users/friends/suggestions/network",
			NetworkUserFriends},

		/* Possible List Routes */
		route.Route{"GET", "/listTypes", GetListTypes},

		/* Whooplist Routes */
		route.Route{"GET",
			"/whooplist/:ListId/coordinate/:Lat/:Long/:Radius/:Page",
			GetWlCoordinate},

		/* Newsfeed Routes */
		route.Route{"GET", "/feed/older/:Lat/:Long/:Radius/:EarliestId",
			GetNewsfeedOlder},
		route.Route{"GET", "/feed/:Lat/:Long/:Radius/:LatestId",
			GetNewsfeedUpdate},
		route.Route{"GET", "/feed/:Lat/:Long/:Radius",
			GetNewsfeedNew},

		/* Place Routes */
		route.Route{"GET",
			"/places/search/:ListId/:Lat/:Long/:Radius/:Page/*SearchString",
			SearchPlace},
		route.Route{"GET", "/places/search/:ListId/:Lat/:Long/:Radius/:Page",
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

package main

import (
	"log"
	"net/http"
	"source.whooplist.com/whooplist"
	"source.whooplist.com/route"
	"github.com/gorilla/context"
)

func main() {
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
		route.Route{"GET", "/whooplist/:Day/:Time", GetWhooplists},	 
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

	err := whooplist.Connect()

	if(err != nil) {
		log.Fatal("Could not connect to database/prepare statements, dying: " + err.Error())
	}

	s := &Server{ Router: router}

	hs := &http.Server{
		Addr:	":3000",
		Handler: context.ClearHandler(logHandler(panicHandler(parseRequest(authenticate(errorHandler(s.handleRequest)))))),
	}

	log.Fatal(hs.ListenAndServe())
}

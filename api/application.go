package main

import (
	"log"
	"net/http"
	"github.com/ant0ine/go-json-rest"
	"source.whooplist.com/whooplist"
	"github.com/gorilla/context"
)

func main() {
	handler := rest.ResourceHandler{}
	handler.EnableRelaxedContentType = true

	handler.SetRoutes(
		/* User Routes */
		rest.Route{"POST", "/users/login", LoginUser},
		rest.Route{"POST", "/users/logout", LogoutUser},
		rest.Route{"GET", "/users/:UserId", GetUser},
		rest.Route{"POST", "/users/:UserId", UpdateUser},
		rest.Route{"PUT", "/users", CreateUser},
		
		/* User List Routes */
		rest.Route{"GET", "/users/:UserId/lists", GetUserLists},
		rest.Route{"GET", "/users/:UserId/lists/:ListId", GetUserList},
		rest.Route{"POST", "/users/:UserId/lists/:ListId", CreateUserList},
		rest.Route{"DELETE", "/users/:UserId/lists/:ListId", DeleteUserList},

		/* Possible List Routes */
		rest.Route{"GET", "/listTypes", GetListTypes},

		/* Whooplist Routes */
		rest.Route{"GET", "/whooplist/:Day/:Time", GetWhooplists},	 
		rest.Route{"GET", "/whooplist/:ListId/:Page/coordinate/:Lat/:Long/:Radius", GetWhooplistCoordinate},
		rest.Route{"GET", "/whooplist/:ListId/:Page/location/:LocationId", GetWhooplistLocation},
		
		/* Newsfeed Routes */
		rest.Route{"GET", "/newsfeed/:Location/:LatestId", GetNewsfeed},
		rest.Route{"GET", "/newsfeed/:Location/older/:EarliestId/", GetNewsfeedOlder},

		/* Location Routes */
		rest.Route{"GET", "/locations/:LocationId", GetLocation},	
		rest.Route{"GET", "/locations/:Latitude/:Longitude", GetLocationsCoordinate},		

		/* Place Routes */
		rest.Route{"GET", "/places/:PlaceId", GetPlace},
	)

	err := whooplist.Connect()

	if(err != nil) {
		log.Fatal("Could not connect to database/prepare statements, dying: " + err.Error())
	}

	log.Fatal(http.ListenAndServe(":3000", context.ClearHandler(
		parseRequest(
		Authenticate(
		errorHandler(&handler)))))
}

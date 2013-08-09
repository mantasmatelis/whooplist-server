package main

import (
	"net/http"
	"github.com/ant0ine/go-json-rest"
)

//Main sets up the routes and handles connections.
func main() {
	handler := rest.ResourceHandler{}

	handler.SetRoutes(
		/* User Routes */
		rest.Route{"GET", "/users/:UserId", GetUser},
		rest.Route{"POST", "/users/:UserId", UpdateUser},
		rest.Route{"POST", "/users/login", LoginUser},
		rest.Route{"POST", "/users/logout", LogoutUser},
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
		rest.Route{"GET", "/whooplist/:ListId/:Page/city/:LocationId", GetWhooplistLocation},
		
		/* Newsfeed Routes */
		rest.Route{"GET", "/newsfeed", GetNewsfeed},
		rest.Route{"GET", "/newsfeed/refresh/:LatestId/", GetNewsfeedRefresh},
		rest.Route{"GET", "/newsfeed/older/:EarliestId/", GetNewsfeedOlder},

		/* Location Routes */
		rest.Route{"GET", "/locations/:LocationId", GetwLocation},	
		rest.Route{"GET", "/locations/:Latitude/:Longitude", GetLocationsCoordinate},		

		/* Place Routes */
		rest.Route{"GET", "/places/:PlaceId", GetPlace},
	)

	http.ListenAndServe("/tmp/whooplist" + ""  + ".socket", &handler)
}

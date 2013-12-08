package main

//TODO: change all Show in controller to Get

import (
	"../whooplist"
	"log"
	"net/http"
)

func Ping(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {

	ensure(context.Session != nil, 403)
	return 200, nil
}

func GetUser(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {

	userId := parseInt64(context.Params["UserId"])
	user, err := whooplist.GetUserData(userId, "")

	if_error(err)
	ensure(user != nil, 404)

	writeObject(user, w)
	return 0, nil
}

func UpdateUser(w http.ResponseWriter, req *http.Request, context Context) (code int, err error) {

	ensure(context.Body != nil && context.Session != nil, 400)
	user := context.Body.User

	user.Id = context.Session.UserId

	var oldUser *whooplist.User
	if user.Password != "" {
		oldUser, err = whooplist.CheckUpdateUser(user.Email, user.OldPassword)
	} else {
		oldUser, err = whooplist.GetUserData(context.Session.UserId, "")
	}
	if_error(err)
	ensure(oldUser != nil, 403)
	ensure(user.Email == oldUser.Email && user.Id == oldUser.Id, 400)

	user.Role = oldUser.Role
	user.PasswordHash = oldUser.PasswordHash

	whooplist.UpdateUser(user)
	return 0, nil
}

func CreateUser(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {

	ensure(context.Body != nil, 400)

	user := context.Body.User

	ensure(user.Email != "" && user.Name != "" && user.Password != "", 400)
	exists, err := whooplist.UserExists(user.Email)
	if_error(err)
	ensure(!exists, 409)

	if_error(whooplist.CreateUser(&user))
	return 0, nil
}

func LoginUser(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {

	ensure(context.Body != nil, 400)

	user, session, err := whooplist.LoginUser(context.Body.User.Email,
		context.Body.User.Password)
	if_error(err)
	ensure(user != nil && session != nil, 403)

	writeObject(&session, w)
	return
}

func LogoutUser(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {

	exist, err := whooplist.DeleteSession(context.Body.Key)
	if_error(err)
	ensure(exist, 403)
	return 200, nil
}

func ExistsUser(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {

	ensure(context.Params["Email"] != "", 400)
	exist, err := whooplist.UserExists(context.Params["Email"])
	if_error(err)
	ensure(exist, 404)
	return 200, nil
}

func GetUserLists(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {

	userId := parseInt64(context.Params["UserId"])

	lists, err := whooplist.GetUserLists(userId)
	if_error(err)

	writeObject(&lists, w)
	return
}

func GetUserList(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {
	userId := parseInt64(context.Params["UserId"])
	listId := parseInt64(context.Params["ListId"])

	list, err := whooplist.GetUserList(userId, listId)
	if_error(err)
	ensure(list != nil, 404)

	writeObject(&list, w)
	return
}

func CreateUserList(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {

	ensure(context.User != nil, 403)
	ensure(context.Body != nil, 400)
	listId := parseInt64(context.Params["ListId"])

	list := whooplist.UserList{}
	list.Items = context.Body.Items
	list.UserId = context.User.Id
	list.ListId = listId

	err = whooplist.PutUserList(list)
	if_error(err)

	return 200, nil
}

func DeleteUserList(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {

	ensure(context.User != nil, 403)

	listId := parseInt64(context.Params["ListId"])
	err = whooplist.DeleteUserList(context.User.Id, listId)
	if_error(err)

	return 200, nil
}

func GetListTypes(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {

	lists, err := whooplist.GetListTypes()
	if_error(err)

	writeObject(&lists, w)
	return
}

func GetWhooplistCoordinate(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {
	return
}

/*func GetWhooplistLocation(w http.ResponseWriter, req *http.Request) {
}*/

func GetNewsfeed(w http.ResponseWriter, req *http.Request) {
}

func GetNewsfeedOlder(w http.ResponseWriter, req *http.Request) {
}

/*func GetLocation(w http.ResponseWriter, req *http.Request) {
}

func GetLocationsCoordinate(w http.ResponseWriter, req *http.Request) {
}*/

func GetPlace(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {

	placeId := parseInt64(context.Params["PlaceId"])
	place, err := whooplist.GetPlace(placeId)
	if_error(err)

	ensure(place != nil, 404)

	writeObject(&place, w)
	return
}

func SearchPlace(w http.ResponseWriter, req *http.Request,
	context Context) (code int, err error) {

	log.Print(context.Params)

	listId := parseInt64(context.Params["ListId"])
	lat := parseFloat64(context.Params["Lat"])
	long := parseFloat64(context.Params["Long"])
	radius := parseFloat64(context.Params["Radius"])
	page := parseInt32(context.Params["Page"])

	places, err := whooplist.SearchPlace(
		context.Params["SearchString"], listId, page, lat, long, radius)
	if_error(err)

	writeObject(&places, w)
	return
}

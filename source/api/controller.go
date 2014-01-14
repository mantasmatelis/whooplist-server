package main

//TODO: change all Show in controller to Get

import (
	"../whooplist"
	"net/http"
)

func Ping(w http.ResponseWriter, req *http.Request, ctx Context) {
	ensure(ctx.Session != nil, 403)
}

func GetUser(w http.ResponseWriter, req *http.Request, ctx Context) {
	userId := parseInt64(ctx.Params["UserId"])
	user, err := whooplist.GetUserData(userId, "")

	if_error(err)
	ensure(user != nil, 404)

	writeObject(user, w)
}

func UpdateUser(w http.ResponseWriter, req *http.Request, ctx Context) {
	ensure(ctx.Body != nil && ctx.Session != nil, 400)

	err := whooplist.UpdateUser(ctx.User, &ctx.Body.User)
	if_error(err)
}

func CreateUser(w http.ResponseWriter, req *http.Request, ctx Context) {
	ensure(ctx.Body != nil, 400)

	user := ctx.Body.User

	ensure(user.Email != nil && user.Name != nil && user.Password != nil, 400)
	ensure(*user.Email != "" && *user.Name != "" && *user.Password != "", 400)
	exists, err := whooplist.UserExists(*user.Email)
	if_error(err)
	ensure(!exists, 409)

	if_error(whooplist.CreateUser(&user))
}

func LoginUser(w http.ResponseWriter, req *http.Request, ctx Context) {
	ensure(ctx.Body != nil, 400)

	user, session, err := whooplist.LoginUser(*ctx.Body.User.Email,
		*ctx.Body.User.Password)
	if_error(err)
	ensure(user != nil && session != nil, 403)

	writeObject(&session, w)
}

func LogoutUser(w http.ResponseWriter, req *http.Request, ctx Context) {
	exist, err := whooplist.DeleteSession(ctx.Body.Key)
	if_error(err)
	ensure(exist, 403)
}

func ExistsUser(w http.ResponseWriter, req *http.Request, ctx Context) {
	ensure(ctx.Params["Email"] != "", 400)
	exist, err := whooplist.UserExists(ctx.Params["Email"])
	if_error(err)
	ensure(exist, 404)
}

func GetUserLists(w http.ResponseWriter, req *http.Request, ctx Context) {
	userId := parseInt64(ctx.Params["UserId"])

	lists, err := whooplist.GetUserLists(userId)
	if_error(err)

	writeObject(&lists, w)
}

func GetUserList(w http.ResponseWriter, req *http.Request, ctx Context) {
	userId := parseInt64(ctx.Params["UserId"])
	listId := parseInt64(ctx.Params["ListId"])

	list, err := whooplist.GetUserList(userId, listId)
	if_error(err)
	ensure(list != nil, 404)

	writeObject(&list, w)
}

func CreateUserList(w http.ResponseWriter, req *http.Request, ctx Context) {
	ensure(ctx.User != nil, 403)
	ensure(ctx.Body != nil, 400)
	listId := parseInt64(ctx.Params["ListId"])

	if_error(whooplist.PutUserList(*ctx.User.Id, listId, ctx.Body.Places))
}

func DeleteUserList(w http.ResponseWriter, req *http.Request, ctx Context) {
	ensure(ctx.User != nil, 403)

	listId := parseInt64(ctx.Params["ListId"])
	if_error(whooplist.DeleteUserList(*ctx.User.Id, listId))
}

type UserFriendsResponse struct {
	Followers []whooplist.User
	Following []whooplist.User
	Both      []whooplist.User
}

func GetUserFriends(w http.ResponseWriter, req *http.Request, ctx Context) {
	userId := parseInt64(ctx.Params["UserId"])

	followers, following, both, err := whooplist.GetUserFriends(userId)
	if_error(err)

	writeObject(UserFriendsResponse{followers, following, both}, w)
}

func AddUserFriend(w http.ResponseWriter, req *http.Request, ctx Context) {
	ensure(ctx.User != nil, 403)
	userId := parseInt64(ctx.Params["UserId"])

	if_error(whooplist.AddUserFriend(*ctx.User.Id, userId))
}

func DeleteUserFriend(w http.ResponseWriter, req *http.Request, ctx Context) {
	ensure(ctx.User != nil, 403)
	userId := parseInt64(ctx.Params["UserId"])

	if_error(whooplist.DeleteUserFriend(*ctx.User.Id, userId))
}

func NetworkUserFriends(w http.ResponseWriter, req *http.Request,
	ctx Context) {

	ensure(ctx.User != nil, 403)

	friends, err := whooplist.NetworkUserFriends(*ctx.User.Id)

	if_error(err)

	writeObject(&friends, w)
}

func ContactsUserFriends(w http.ResponseWriter, req *http.Request,
	ctx Context) {

	ensure(ctx.User != nil, 403)

	friends, err := whooplist.ContactsUserFriends(*ctx.User.Id,
		ctx.Body.Contacts)
	if_error(err)

	writeObject(&friends, w)
}

func SuggestUserFriends(w http.ResponseWriter, req *http.Request,
	ctx Context) {

	ensure(ctx.User != nil, 403)

	friends, err := whooplist.SuggestUserFriends(*ctx.User.Id,
		ctx.Body.Contacts)
	if_error(err)

	writeObject(&friends, w)
}

func GetListTypes(w http.ResponseWriter, req *http.Request, ctx Context) {
	lists, err := whooplist.GetListTypes()
	if_error(err)

	writeObject(&lists, w)
}

func GetWlCoordinate(w http.ResponseWriter, req *http.Request, ctx Context) {
	userId := int64(0)
	if ctx.User != nil {
		userId = *ctx.User.Id
	}

	listId := parseInt64(ctx.Params["ListId"])
	lat := parseFloat64(ctx.Params["Lat"])
	long := parseFloat64(ctx.Params["Long"])
	radius := parseFloat64(ctx.Params["Radius"])
	page := parseInt32(ctx.Params["Page"])

	list, err := whooplist.GetWhooplistCoordinate(
		userId, listId, page, lat, long, radius)
	if_error(err)

	writeObject(&list, w)
}

func GetNewsfeedNew(w http.ResponseWriter, req *http.Request, ctx Context) {
	ctx.Params["LatestId"] = "-1"
	GetNewsfeedUpdate(w, req, ctx)
}

func GetNewsfeedUpdate(w http.ResponseWriter, req *http.Request, ctx Context) {
	ensure(ctx.User != nil, 403)

	lat := parseFloat64(ctx.Params["Lat"])
	long := parseFloat64(ctx.Params["Long"])
	radius := parseFloat64(ctx.Params["Radius"])
	latestId := parseInt64(ctx.Params["LatestId"])

	items, err := whooplist.GetNewsfeed(
		*ctx.User.Id, latestId, lat, long, radius)
	if_error(err)

	writeObject(&items, w)
}

func GetNewsfeedOlder(w http.ResponseWriter, req *http.Request, ctx Context) {
	ensure(ctx.User != nil, 403)

	lat := parseFloat64(ctx.Params["Lat"])
	long := parseFloat64(ctx.Params["Long"])
	radius := parseFloat64(ctx.Params["Radius"])
	earliestId := parseInt64(ctx.Params["LatestId"])

	items, err := whooplist.GetNewsfeedEarlier(
		*ctx.User.Id, earliestId, lat, long, radius)
	if_error(err)

	writeObject(&items, w)
}

func GetPlace(w http.ResponseWriter, req *http.Request, ctx Context) {
	placeId := parseInt64(ctx.Params["PlaceId"])
	place, err := whooplist.GetPlace(placeId)
	if_error(err)

	ensure(place != nil, 404)

	writeObject(&place, w)
}

func SearchPlace(w http.ResponseWriter, req *http.Request, ctx Context) {
	listId := parseInt64(ctx.Params["ListId"])
	lat := parseFloat64(ctx.Params["Lat"])
	long := parseFloat64(ctx.Params["Long"])
	radius := parseFloat64(ctx.Params["Radius"])
	page := parseInt32(ctx.Params["Page"])

	places, err := whooplist.SearchPlace(
		ctx.Params["SearchString"], listId, page, lat, long, radius)
	if_error(err)

	writeObject(&places, w)
}

package main

//TODO: change all Show in controller to Get

import (
	"log"
	"strconv"
	"net/http"
	"github.com/gorilla/context"
	"source.whooplist.com/whooplist"
)

func GetUser(w http.ResponseWriter, req *http.Request) (code int, err error) {

	userId, err := strconv.ParseInt(req.Form.Get("UserId"), 10, 64)

	if err != nil {
		return 400, err
	}

	user, err := whooplist.GetUserData(userId, "")

	if err != nil {
		return 500, err
	}

	if user == nil {
		return 404, nil
	}

	writeObject(user, w)
	return 0, nil
}

func UpdateUser(w http.ResponseWriter, req *http.Request) (code int, err error) {
	body, _ := context.Get(req, Body).(*RequestBody)
	if body == nil {
		return 400, nil
	}

	user := body.User

	var oldUser *whooplist.User
	if user.Password != "" {
		oldUser, err = whooplist.CheckUpdateUser(user.Email, user.OldPassword)
	} else {
		id, err := strconv.ParseInt(req.Form.Get("UserId"), 10, 64)
		if err != nil {
			return 400, err
		}
		oldUser, err = whooplist.GetUserData(id, "")
	}
	if err != nil {
		return 500, err
	}
	if oldUser == nil {
		return 403, nil
	}
	if user.Email != oldUser.Email {
		return 400, nil
	}

	user.Role = oldUser.Role
	user.PasswordHash = oldUser.PasswordHash

	whooplist.UpdateUser(user)
	return 0, nil
}

func CreateUser(w http.ResponseWriter, req *http.Request) (code int, err error) {
	log.Print("creating user")
	body, _ := context.Get(req, Body).(*RequestBody)
	if body == nil {
		return 400, nil
	}
	user := body.User

	if user.Email == "" || user.Name == "" || user.Password == "" {
		return 400, nil
	}

	//TODO: Check case where e-mail already exists.
	//Include password strength requirements.
	//409 conflict, 406 bad password
	if err = whooplist.CreateUser(user); err != nil {
		return 500, err
	}
	return 0, nil
}

func LoginUser(w http.ResponseWriter, req *http.Request) (code int, err error) {
	body, _ := context.Get(req, Body).(*RequestBody)
	user, session, err := whooplist.LoginUser(body.User.Email, body.User.Password)
	if err != nil {
		return 500, err
	}
	if user == nil || session == nil {
		return 403, nil
	}

	writeObject(&session, w)
	return
}

func LogoutUser(w http.ResponseWriter, req *http.Request) (code int, err error) {
	body, _ := context.Get(req, Body).(*RequestBody)

	exist, err := whooplist.DeleteSession(body.Key)
	if err != nil {
		return 500, err
	}
	if !exist {
		return 403, err
	}
	w.WriteHeader(200)
	return
}

func GetUserLists(w http.ResponseWriter, req *http.Request) (code int, err error) {

	userId, err := strconv.Atoi(req.Form.Get("UserId"))

	if err != nil {
		return 400, err
	}

	lists, err := whooplist.GetUserLists(userId)
	if err != nil {
		return 500, err
	}

	writeObject(&lists, w)
	return
}

func GetUserList(w http.ResponseWriter, req *http.Request) (code int, err error) {
	userId, err := strconv.Atoi(req.Form.Get("UserId"))

	if err != nil {
		return 400, err
	}

	listId, err := strconv.Atoi(req.Form.Get("ListId"))

	if err != nil {
		return 400, nil
	}

	list, err := whooplist.GetUserList(userId, listId)
	if err != nil {
		return 500, err
	}
	if list == nil {
		return 404, nil
	}

	writeObject(&list, w)
	return
}

func CreateUserList(w http.ResponseWriter, req *http.Request) {
}

func DeleteUserList(w http.ResponseWriter, req *http.Request) {
}

func GetListTypes(w http.ResponseWriter, req *http.Request) {
}

func GetWhooplists(w http.ResponseWriter, req *http.Request) {
}

func GetWhooplistCoordinate(w http.ResponseWriter, req *http.Request) {
}

func GetWhooplistLocation(w http.ResponseWriter, req *http.Request) {
}

func GetNewsfeed(w http.ResponseWriter, req *http.Request) {
}

func GetNewsfeedOlder(w http.ResponseWriter, req *http.Request) {
}

func GetLocation(w http.ResponseWriter, req *http.Request) {
}

func GetLocationsCoordinate(w http.ResponseWriter, req *http.Request) {
}

func GetPlace(w http.ResponseWriter, req *http.Request) {
}

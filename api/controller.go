package main

//TODO: change all Show in controller to Get

import (
	"net/http"
	"github.com/ant0ine/go-json-rest"
	"whooplist.com/whooplist"
)

func GetUser(w *rest.ResponseWriter, req *rest.Request) {
	user, err := whooplist.GetUserData(req.PathParam("UserId"))
	
	if err != nil {
		http.Error(w, "", 500)
		return
	}
	
	if user == nil {
		http.Error(w, "", 404)
		return
	}

	outUser user := User{Id: user.Id, Email: user.Email, Name: user.Name}	
	w.WriteJson(&outUser)
}

func UpdateUser(w *rest.ResponseWriter, req *rest.Request) {
}

func CreateUser(w *rest.ResponseWriter, req *rest.Request) {
	
}

func LoginUser(w *rest.ResponseWriter, req *rest.Request) {
}

func LogoutUser(w *rest.ResponseWriter, req *rest.Request) {
	user, session, err := AuthUser(req.FormValue("key"))
	if err != nil {
		http.Error(w, "", 500)
	}
	if user == nil || session == nil {
		http.Error(w, "", 403)
		return
	}

	err = DeleteSession(req.FormValue("key"))
	if err != nil {
		http.error(w, "", 500)
		return
	}
}

func GetUserLists(w *rest.ResponseWriter, req *rest.Request) {
	lists, err = RetrieveUserLists(req.FormValue("UserId")
	if err != nil {
		http.Error(w, "", 500)
		return
	}
	
	w.WriteJson(&lists)
}

func GetUserList(w *rest.ResponseWriter, req *rest.Request) {
	list, err = RetrieveUserList(req.FormValue("UserId", req.FormValue("ListId"))
	if err != nil {
		http.Error(w, "", 500)
		return
	}
	
	w.WriteJson(&list)
}

func CreateUserList(w *rest.ResponseWriter, req *rest.Request) {
	
}

func DeleteUserList(w *rest.ResponseWriter, req *rest.Request) {
}

func GetListTypes(w *rest.ResponseWriter, req *rest.Request) {
}

func GetWhooplists(w *rest.ResponseWriter, req *rest.Request) {	
}

func GetWhooplistCoordinate(w *rest.ResponseWriter, req *rest.Request) {
}

func GetWhooplistLocation(w *rest.ResponseWriter, req *rest.Request) {
}

func GetLocation(w *rest.ResponseWriter, req *rest.Request) {
}

func GetLocationsCoordinate(w *rest.ResponseWriter, req *rest.Request) {
}

func GetPlace(w *rest.ResponseWriter, req *rest.Request) {
}

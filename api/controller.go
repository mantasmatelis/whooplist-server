package main

//TODO: change all Show in controller to Get

import (
	"strconv"
	"log"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"github.com/ant0ine/go-json-rest"
	"github.com/gorilla/context"
	"source.whooplist.com/whooplist"
)



type RequestBody struct {
	Key string
	User whooplist.User
	Place whooplist.Place
	UserList whooplist.UserList
}

type key int
const Body key = 0
const User key = 1
const Session key = 2


func parseRequest(handler http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    body, err := readRequest(r)
    if err == nil {
      context.Set(r, Body, body)
    } else {
      log.Print("Error parsing request: " + err.Error())
    }
    handler.ServeHTTP(w, r)
  })
}

func readRequest(req *http.Request) (body *RequestBody, err error) {
	content, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		return
	}

	err = json.Unmarshal(content, &body)
	
	log.Print(body)

	return	
}

func authenticate(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r* http.Request) {
		auth(r)
		handler.ServeHTTP(w, r)
	})
}

func auth(r *http.Request) {
	body,_ := context.Get(r, Body).(*RequestBody)
	if body != nil && body.Key != "" {
		user, session, _ := whooplist.AuthUser(body.Key)
		context.Set(r, User, user)
		context.Set(r, Session, session)
 	}
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, err := f(w, r)
		if err != nil {
			http.Error(w, "", code)
			if(code == 500) {
				log.Print(err.Error())
			}
		}
	}
}

func GetUser(w *rest.ResponseWriter, req *rest.Request) {

	userId, err := strconv.ParseInt(req.PathParam("UserId"), 10, 64)

	if err != nil {
		http.Error(w, "", 400)
		return
	}

	user, err := whooplist.GetUserData(userId)
	
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "", 500)
		return
	}
	
	if user == nil {
		http.Error(w, "", 404)
		return
	}

	w.WriteJson(&user)
}

func UpdateUser(w *rest.ResponseWriter, req *rest.Request) {
	body, _ := context.Get(req.Request, Body).(*RequestBody)
	if body == nil {
		http.Error(w, "", 400)
		return
	}

	user := body.User
}

func CreateUser(w *rest.ResponseWriter, req *rest.Request) {
	body, _ := context.Get(req.Request, Body).(*RequestBody)
	if body == nil {
		http.Error(w, "", 400)
		return
	}
	user := body.User

	if user.Email == "" || user.Name == "" || user.Password == "" {
		http.Error(w, "", 400)		
		return
	}
	err = whooplist.CreateUser(user)
	
	if err != nil {
		http.Error(w, "", 500)
		log.Print("Could not create user: " + err.Error())
	}
}

func LoginUser(w *rest.ResponseWriter, req *rest.Request) {
	body, _ := context.Get(req.Request, Body).(*RequestBody)
	user, session, err := whooplist.LoginUser(body.User.Email, body.User.Password)
	log.Print(user,session)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "", 500)
		return
	}
	if user == nil || session == nil {
		http.Error(w, "", 403)
		return
	}

	w.WriteJson(&session)
}

func LogoutUser(w *rest.ResponseWriter, req *rest.Request) {
	body, _ := context.Get(req.Request, Body).(*RequestBody)

	exist, err := whooplist.DeleteSession(body.Key)
	if err != nil {
                http.Error(w, "", 500)
                return
        }
	if !exist {
		http.Error(w, "", 403)
		return
	}
	w.WriteHeader(200)
}

func GetUserLists(w *rest.ResponseWriter, req *rest.Request) {

	userId, err := strconv.Atoi(req.PathParam("UserId"))

        if err != nil {
                http.Error(w, "", 400)
        }

	lists, err := whooplist.GetUserLists(userId)
	if err != nil {
		http.Error(w, "", 500)
		return
	}
	
	w.WriteJson(&lists)
}

func GetUserList(w *rest.ResponseWriter, req *rest.Request) {
	userId, err := strconv.Atoi(req.PathParam("UserId"))

        if err != nil {
                http.Error(w, "", 400)
        }

	listId, err := strconv.Atoi(req.PathParam("ListId"))

        if err != nil {
                http.Error(w, "", 400)
        }

	list, err := whooplist.GetUserList(userId, listId)
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

func GetNewsfeed(w *rest.ResponseWriter, req *rest.Request) {
}

func GetNewsfeedOlder(w *rest.ResponseWriter, req *rest.Request) {
}

func GetLocation(w *rest.ResponseWriter, req *rest.Request) {
}

func GetLocationsCoordinate(w *rest.ResponseWriter, req *rest.Request) {
}

func GetPlace(w *rest.ResponseWriter, req *rest.Request) {
}

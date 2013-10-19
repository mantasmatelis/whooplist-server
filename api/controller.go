package main

//TODO: change all Show in controller to Get

import (
	"strconv"
	"log"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"net/url"
	"runtime/debug"
	"github.com/gorilla/context"
	"source.whooplist.com/route"
	"source.whooplist.com/whooplist"
)

func logHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("request: ", r.URL.Path)
		handler.ServeHTTP(w, r)
		log.Print("end request")
	})
}

func panicHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Print("\n***PANIC***\n", rec, "\n\n", string(debug.Stack()), "\n***END PANIC***")
				http.Error(w, "", 500)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}

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
	if err != nil || string(content) == "" {
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

func errorHandler(f func(http.ResponseWriter, *http.Request)(int, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code, err := f(w, r)
		if code != 0 {
			http.Error(w, "", code)
			if(err != nil) {
				log.Print("code: ", code, " handling request: ", err.Error())
			}
		}
	})
}

type Server struct {
	Router route.Router	
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) (code int, err error) {
	route, params, pathMatched := s.Router.FindRouteFromURL(r.Method, r.URL)
	if route == nil && pathMatched {
		http.Error(w, "", 405)
		return
	}
	if route == nil {
		http.Error(w, "", 400)
		return
	}

	r.Form = url.Values{}
	for key, value := range params {
		
		r.Form.Set(key, value)
	}

	return route.Func.(func(http.ResponseWriter, *http.Request)(int,error))(w,r)
}

func writeObject(obj interface{}, w http.ResponseWriter) (err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)	
	return
}

func GetUser(w http.ResponseWriter, req *http.Request) (code int, err error) {

	userId, err := strconv.ParseInt(req.Form.Get("UserId"), 10, 64)

	if err != nil {
		return 400, err
	}

	user, err := whooplist.GetUserData(userId)
	
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
		oldUser, err = whooplist.GetUserData(id)
	}
	if err != nil {
		return 500, err
	}
	if oldUser == nil {
		return 403, nil
	}
	if(user.Email != oldUser.Email) {
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
	
	if err := whooplist.CreateUser(user); err != nil {
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

func GetUserList(w http.ResponseWriter, req *http.Request) {
	userId, err := strconv.Atoi(req.Form.Get("UserId"))

        if err != nil {
                http.Error(w, "", 400)
        }

	listId, err := strconv.Atoi(req.Form.Get("ListId"))

        if err != nil {
                http.Error(w, "", 400)
        }

	list, err := whooplist.GetUserList(userId, listId)
	if err != nil {
		http.Error(w, "", 500)
		return
	}
	
	writeObject(&list, w)
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

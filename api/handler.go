package main

import (

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


package main

import (
	"../whooplist"
	"encoding/json"
	"github.com/mantasmatelis/go-trie-url-route"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	code *int /* This is an *int because fuck go interfaces. */
}

func (w loggingResponseWriter) WriteHeader(code int) {
	/* By that, I mean I have to "w l..." and not " w *l..." */
	*w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func logHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		baseCode := 200
		loggingW := loggingResponseWriter{code: &baseCode, ResponseWriter: w}
		log.Print("Request: ", r.Method, " ", r.URL.Path)
		handler.ServeHTTP(loggingW, r)
		log.Print("Response: ", *loggingW.code)
	})
}

func panicHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				/* If we're dealing with an internal API error,
				   return the proper code and drop out */
				if aE, ok := rec.(apiError); ok {
					http.Error(w, "", aE.code)
					if aE.err != nil {
						log.Print("Error handling request: ",
							aE.err.Error())
					}
					return
				}
				/* Otherwise, we're dealing with a real error */
				log.Print("\n***PANIC***\n", rec, "\n\n",
					string(debug.Stack()), "\n***END PANIC***")
				http.Error(w, "", 500)
			}
		}()
		handler.ServeHTTP(w, r)
	})
}

type RequestBody struct {
	Key    string
	User   whooplist.User
	Places []int64
}

type Context struct {
	Params  map[string]string
	Body    *RequestBody
	User    *whooplist.User
	Session *whooplist.Session
}

type Server struct {
	Router route.Router
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	var context Context
	var err error
	context.Body, err = readRequest(r)
	if err != nil {
		log.Print("Error parsing request: " + err.Error())
		http.Error(w, "", 400)
		return
	}

	if context.Body != nil && context.Body.Key != "" {
		context.User, context.Session, err = whooplist.AuthUser(context.Body.Key)
	}
	if err != nil {
		log.Print("Error authenticating user: " + err.Error())
		http.Error(w, "", 500)
		return
	}

	route, params, pathMatched := s.Router.FindRouteFromURL(r.Method, r.URL)

	if route == nil && pathMatched {
		http.Error(w, "", 405)
		return
	}
	if route == nil {
		http.Error(w, "", 400)
		return
	}

	context.Params = params

	route.Func.(func(http.ResponseWriter, *http.Request,
		Context))(w, r, context)

	w.Write([]byte{}) /* Force a 200 status if none sent */
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

func writeObject(obj interface{}, w http.ResponseWriter) (err error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	return
}

type apiError struct {
	code int
	err  error
}

func ensure(condition bool, code int) {
	if !condition {
		panic(apiError{code: code})
	}
}

func if_error(err error) {
	if err != nil {
		panic(apiError{code: 500, err: err})
	}

}

func parseInt32(str string) int32 {
	num, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		log.Print(err)
	}
	ensure(err == nil, 400)
	return int32(num)
}

func parseInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Print(err)
	}
	ensure(err == nil, 400)
	return num
}

func parseFloat64(str string) float64 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Print(err)
	}
	ensure(err == nil, 400)
	return num
}

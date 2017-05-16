package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Route struct {
	URL        string
	Method     string
	Handler    http.HandlerFunc
	Middleware []Middleware
}

func CreateRoute(router *httprouter.Router, r *Route) {
	router.Handle(r.Method, r.URL, func(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// 244 ns/op vs 476 ns/op when adding context
		params := map[string]string{}
		// add url paramters to context if they exist
		if len(ps) > 0 {
			for _, b := range ps {
				params[b.Key] = b.Value
			}
			ctx := SetContext(req.Context(), Params, params)
			req = req.WithContext(ctx)
		}

		adapt(r.Handler, r.Middleware...).ServeHTTP(res, req)
	})
}

func GroupRoutes(router *httprouter.Router, middleware []Middleware, routes []*Route) {
	for _, r := range routes {
		r.Middleware = append(middleware, r.Middleware...)
		CreateRoute(router, r)
	}
}

type Middleware func(http.Handler) http.Handler

func adapt(h http.Handler, adapters ...Middleware) http.Handler {
	l := len(adapters) - 1
	for i := l; i >= 0; i-- {
		adapter := adapters[i]
		h = adapter(h)
	}
	return h
}

// JSON
func WriteJSON(res http.ResponseWriter, status int, obj interface{}) error {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(status)

	return json.NewEncoder(res).Encode(obj)
}

package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ROUTE EXAMPLES
func initRoutes(router *httprouter.Router, dep *Dependencies) {
	/*
		single route
	*/
	indexRoute := &Route{
		URL:     "/",
		Method:  "GET",
		Handler: DisplaySuccess(dep.SomeConfig),
	}
	CreateRoute(router, indexRoute)

	/*
		single route with dynamic segment via context
	*/
	dynamicSegmentRoute := &Route{
		URL:     "/params-example/:commentid/:userid",
		Method:  "GET",
		Handler: DisplayID(dep.SomeConfig),
	}
	CreateRoute(router, dynamicSegmentRoute)

	/*
		single route with middleware + context
	*/
	middleWareExample := &Route{
		URL:     "/profile",
		Method:  "GET",
		Handler: DisplayProfile(),
		Middleware: []Middleware{
			LogSomething(dep.LoggerSettings),
			AddUserProfileContext(dep.SomeService),
		},
	}
	CreateRoute(router, middleWareExample)

	/*
		multiple routes with common middleware
	*/
	routesToGroup := []*Route{
		&Route{
			URL:     "/route-a",
			Method:  "GET",
			Handler: DisplaySuccess(dep.SomeConfig),
		},
		&Route{
			URL:     "/route-b",
			Method:  "GET",
			Handler: DisplaySuccess(dep.SomeConfig),
		},
		&Route{
			URL:     "/route-c",
			Method:  "GET",
			Handler: DisplaySuccess(dep.SomeConfig),
		},
	}
	GroupRoutes(router, []Middleware{
		LogSomething(dep.LoggerSettings),
	}, routesToGroup)

	/*
		multiple routes with same AND different middle ware
	*/
	moreGroupedRoutes := []*Route{
		&Route{
			URL:     "/route-x",
			Method:  "GET",
			Handler: DisplaySuccess(dep.SomeConfig),
		},
		// this route has another extra middleware
		&Route{
			URL:     "/route-y",
			Method:  "GET",
			Handler: DisplaySuccess(dep.SomeConfig),
			Middleware: []Middleware{
				RandomMiddleware(),
			},
		},
	}
	GroupRoutes(router, []Middleware{
		LogSomething(dep.LoggerSettings),
	}, moreGroupedRoutes)
}

/*

	HANDLERS

*/
func DisplaySuccess(conf *SomeConfig) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// could do something with config passed in

		res.Write([]byte("Success! " + req.URL.String()))
	})
}

func DisplayID(conf *SomeConfig) http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// get route parameters
		params, _ := GetURLParams(req.Context())
		commentid := params["commentid"]
		if commentid == "" {
			res.WriteHeader(500)
			res.Write([]byte("Failed to get commentid"))
			return
		}
		userid := params["userid"]
		if userid == "" {
			res.WriteHeader(500)
			res.Write([]byte("Failed to get userid"))
			return
		}

		res.Write([]byte(fmt.Sprintf("Comment id is %v and user id is %v", commentid, userid)))
	})
}

func DisplayProfile() http.HandlerFunc {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// retrieve profile from context
		profile, err := GetProfileFromContext(req.Context())
		if err != nil {
			res.WriteHeader(500)
			res.Write([]byte("Failed to get profile"))
			return
		}

		// do whatever else
		// ...

		// can return profile
		WriteJSON(res, 200, profile)
	})
}

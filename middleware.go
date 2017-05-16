package main

import (
	"fmt"
	"net/http"
)

func LogSomething(logger *LoggerSettings) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			// log whatever
			fmt.Println(logger.LogMsgPrefix, req.URL)

			// calls next middleware
			next.ServeHTTP(res, req)
		})
	}
}

func AddUserProfileContext(someAuthService *SomeService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

			// check for some token
			authorization := req.Header.Get("Authorization")
			if authorization == "" {
				res.WriteHeader(400)
				res.Write([]byte("Must provide authorization"))
				return
			}

			// could grab profile with some service
			profile := someAuthService.getProfile(authorization)

			// add to req context
			ctx := SetContext(req.Context(), ProfileContext, profile)

			// call next middleware
			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}

func RandomMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			fmt.Println("Random middleware is called")

			// calls next middleware
			next.ServeHTTP(res, req)
		})
	}
}

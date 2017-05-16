package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// EXAMPLE of dependency
type Dependencies struct {
	*SomeConfig
	*LoggerSettings
	*SomeService
	Port string
}
type SomeConfig struct {
	Val1 string
	Val2 string
}
type LoggerSettings struct {
	LogMsgPrefix string
}
type SomeService struct{}

func (s *SomeService) getProfile(token string) *Profile {
	fmt.Println("Would fetch profile with token", token)
	return &Profile{
		Username: "Timmy",
		ID:       "0t97afd09uvan094eutn",
		Age:      22,
	}
}

// PROFILE EXAMPLE
type Profile struct {
	Username string
	ID       string
	Age      int
}

func main() {
	// some example dependency is created here
	dep := Dependencies{
		Port: ":5555",
		SomeConfig: &SomeConfig{
			Val1: "a value",
			Val2: "another value",
		},
		LoggerSettings: &LoggerSettings{
			LogMsgPrefix: "LOGGING ROUTE:",
		},
		SomeService: &SomeService{},
	}

	// create router
	router := httprouter.New()

	// init routes
	initRoutes(router, &dep)

	log.Fatal(http.ListenAndServe(dep.Port, router))
}

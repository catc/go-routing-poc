package main

import (
	"context"
	"errors"
)

type key string

// String method for info about context when debuggings
func (c key) String() string {
	return "context key " + string(c)
}

const (
	ProfileContext = key("profile")

	// Params are url params from httprouter
	Params = key("urlParams")
)

// SetContext sets context user data
func SetContext(ctx context.Context, k key, val interface{}) context.Context {
	return context.WithValue(ctx, k, val)
}

// retrieve profilel from context
func GetProfileFromContext(ctx context.Context) (*Profile, error) {
	profile, ok := ctx.Value(ProfileContext).(*Profile)

	if !ok {
		return profile, errors.New("failed to get profile from context")
	}
	return profile, nil
}

// GetURLParams attempts to retrieve url parameter data from context
func GetURLParams(ctx context.Context) (map[string]string, bool) {
	params, ok := ctx.Value(Params).(map[string]string)
	return params, ok
}

package main

import (
	"log"
	"os"
)

type AuthzUser struct {
	context *RequestContext
}

// NewAuthzUser reuturns new AuthzUser struct
func NewAuthzUser(req *AuthorizationRequest) *AuthzUser {
	context := NewRequestContext(req)

	return &AuthzUser{
		context: context,
	}
}

// IsAllowed checks if service account can access resource
// returns true on success, false otherwise
func (r *AuthzUser) IsAllowed() bool {
	for _, entry := range config.Rules {
		accessMode := entry.GetAccessMode(r.context)
		if accessMode == ALLOW {
			debug("user %s (ns %s): matched ALLOW entry %+v", r.Username(), r.Request().Namespace(), entry)
			return true
		} else if accessMode == DENY {
			debug("user %s (ns %s): matched DENY entry %+v", r.Username(), r.Request().Namespace(), entry)
			return false
		}
	}
	debug("user %s (ns %s): no matches, default DENY", r.Username(), r.Request().Namespace())
	return false
}

// Username returns request's spec.user
func (r *AuthzUser) Username() string {
	return r.context.Username
}

// Request returns full request struct
func (r *AuthzUser) Request() *AuthorizationRequest {
	return r.context.Request
}

func debug(msg string, args ...interface{}) {
	debug := os.Getenv("DEBUG")
	if debug != "" {
		log.Printf(msg, args...)
	}
}

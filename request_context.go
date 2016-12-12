package main

// RequestContext holds request struct
type RequestContext struct {
	Request        *AuthorizationRequest
	ServiceAccount *ServiceAccount
	Username       string
	UserGroups     []string
}

// NewRequestContext builds request context out of req object
func NewRequestContext(req *AuthorizationRequest) *RequestContext {
	username := req.Spec.User
	usergroups := req.Spec.Group
	serviceAccount := NewServiceAccount(username)

	return &RequestContext{
		Request:        req,
		ServiceAccount: serviceAccount,
		Username:       username,
		UserGroups:     usergroups,
	}
}

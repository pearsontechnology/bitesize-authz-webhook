package main

import "encoding/json"
import "io"

// ResourceAttributesSpec defines resource attributes
type ResourceAttributesSpec struct {
	Namespace string `json:"namespace,omitempty"`
	Verb      string `json:"verb"`
	Group     string `json:"group,omitempty"`
	Resource  string `json:"resource"`
}

// NonResourceAttributesSpec defines non-resource attributes (like /api)
type NonResourceAttributesSpec struct {
	Path string `json:"path"`
	Verb string `json:"verb"`
}

// AuthorizationRequestSpec represents auth request structure found in HTTP call
type AuthorizationRequestSpec struct {
	NonResourceAttributes *NonResourceAttributesSpec `json:"nonResourceAttributes,omitempty"`
	ResourceAttributes    *ResourceAttributesSpec    `json:"resourceAttributes,omitempty"`
	User                  string                     `json:"user"`
	Group                 []string                   `json:"group,omitempty"`
}

// AuthorizationRequest represents incoming HTTP request body
type AuthorizationRequest struct {
	ApiVersion string                   `json:"apiVersion"`
	Kind       string                   `json:"kind"`
	Spec       AuthorizationRequestSpec `json:"spec"`
}

// NewAuthorizationRequest returns AuthorizationRequest struct based on
// HTTP request body
func NewAuthorizationRequest(body io.Reader) (*AuthorizationRequest, error) {
	var req *AuthorizationRequest

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&req)

	return req, err
}

// Namespace returns namespace from request
func (r *AuthorizationRequest) Namespace() string {
	if !r.IsResourceRequest() {
		return "__MAGIC__NOMATCH_*_KEY__"
	}
	return r.Spec.ResourceAttributes.Namespace
}

func (r *AuthorizationRequest) IsResourceRequest() bool {
	return r.Spec.ResourceAttributes != nil
}

func (r *AuthorizationRequest) Action() string {
	if !r.IsResourceRequest() {
		return r.Spec.NonResourceAttributes.Verb
	}
	return r.Spec.ResourceAttributes.Verb
}

// Path is in NonResourceAttributes only
func (r *AuthorizationRequest) Path() string {
	if r.IsResourceRequest() {
		return "__MAGIC__NOMATCH_*_KEY__"
	}
	return r.Spec.NonResourceAttributes.Path
}

func (r *AuthorizationRequest) Group() string {
	if !r.IsResourceRequest() {
		return "__MAGIC__NOMATCH_*_KEY__"
	}
	return r.Spec.ResourceAttributes.Group
}

func (r *AuthorizationRequest) Resource() string {
	if !r.IsResourceRequest() {
		return "__MAGIC__NOMATCH_*_KEY__"
	}
	return r.Spec.ResourceAttributes.Resource
}

func (r *AuthorizationRequest) ServiceAccount() *ServiceAccount {
	return NewServiceAccount(r.Spec.User)
}

func (r *AuthorizationRequest) UserGroups() []string {
	if !r.IsResourceRequest() {
		return nil
	}
	return r.Spec.Group
}

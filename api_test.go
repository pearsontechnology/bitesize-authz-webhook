package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	server *httptest.Server
	reader io.Reader
	index  string
)

func init() {
	server = httptest.NewServer(Handlers())
	index = fmt.Sprintf("%s/", server.URL)

	LoadConfigFromFile("examples/example1.hcl")
}

func postIndex(data string) (*http.Response, error) {
	reader = strings.NewReader(data)

	request, err := http.NewRequest("POST", index, reader)
	res, err := http.DefaultClient.Do(request)
	return res, err

}

func parseResponse(r *http.Response) (*AuthorizationResponse, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	respJSON := NewAuthorizationResponse(true)
	err = json.Unmarshal(body, &respJSON)
	return respJSON, err
}

func TestNewAuthorizationByNamespace(t *testing.T) {
	var saTests = []struct {
		namespace string
		username  string
		status    int
		allowed   bool
	}{
		{"namespace-dev", "system:serviceaccount:namespace-sss:default", 200, true},
		{"namespace", "system:serviceaccount:namespace-sss:default", 200, true},
		{"namespace-any", "badguy", 403, false},
		{"namespace-dev", "someuser", 200, true},
		{"namespace-dev", "system:serviceaccount:default-ns:default", 403, false},
		{"sample-app", "system:serviceaccount:sample-app:default", 200, true},
		{"sample-app", "system:serviceaccount:sample-app:default", 200, true},
		{"sample-app-dev", "system:serviceaccount:kube-system:default", 200, true},
	}

	for _, tst := range saTests {

		reqJSON := fmt.Sprintf(`
    {
      "spec":{
        "resourceAttributes": {
          "namespace":"%s",
          "verb": "get"
        },
        "user":"%s"
      }
    }`, tst.namespace, tst.username)

		result, err := postIndex(reqJSON)

		if err != nil {
			t.Error(err)
		}

		// if result.StatusCode != tst.status {
		//  	t.Errorf("Expected status %d, got: %d (%s)", tst.status, result.StatusCode, reqJSON)
		// }

		respJSON, err := parseResponse(result)
		if err != nil {
			t.Error(err)
		}

		if respJSON.Status.Allowed != tst.allowed {
			t.Errorf("Bad response status, expected %t got: %t (%s)", tst.allowed, respJSON.Status.Allowed, reqJSON)
		}
	}
}

func TestNewAuthorizationBadRequest(t *testing.T) {
	reqJSON := `asd`
	result, err := postIndex(reqJSON)

	if err != nil {
		t.Error(err)
	}

	if result.StatusCode != 400 {
		t.Errorf("Bad request expected: %d", result.StatusCode)
	}
}

func TestNewAuthorizationByPath(t *testing.T) {
	var saTests = []struct {
		path     string
		verb     string
		username string
		allowed  bool
	}{
		{"/apis", "get", "system:serviceaccount:random:default", true},
		{"/api", "get", "system:serviceaccount:random:default", true},
		{"/version", "get", "system:serviceaccount:random:default", true},
		{"/swaggerapi/apis/extensions/v1beta1", "get", "system:serviceaccount:random:default", true},
		{"/api/v1", "get", "system:serviceaccount:random:default", false},
	}

	for _, tst := range saTests {
		reqJSON := fmt.Sprintf(`
    {
      "spec":{
        "nonResourceAttributes":{
          "path": "%s",
          "verb": "%s"
        },
        "user": "%s"
      }
    }`, tst.path, tst.verb, tst.username)
		result, err := postIndex(reqJSON)

		if err != nil {
			t.Error(err)
		}

		respJSON, err := parseResponse(result)
		if err != nil {
			t.Error(err)
		}

		if respJSON.Status.Allowed != tst.allowed {
			t.Errorf("Bad response status, expected %t got: %t (%s)", tst.allowed, respJSON.Status.Allowed, reqJSON)
		}
	}
}

func TestNewAuthorizationByVerb(t *testing.T) {

	var saTests = []struct {
		namespace string
		verb      string
		username  string
		allowed   bool
	}{
		{"default", "watch", "system:serviceaccount:kube-system:default", true},
		{"default", "watch", "system:serviceaccount:random:default", false},
		{"default", "list", "system:serviceaccount:kube-system:default", true},
		{"default", "list", "system:serviceaccount:random:default", false},
	}

	for _, tst := range saTests {

		reqJSON := fmt.Sprintf(`
    {
      "spec":{
        "resourceAttributes":{
          "namespace":"%s",
          "verb":"%s",
          "resource": "services"
        },
        "user":"%s"
      }
    }`, tst.namespace, tst.verb, tst.username)
		result, err := postIndex(reqJSON)

		if err != nil {
			t.Error(err)
		}
		respJSON, err := parseResponse(result)
		if err != nil {
			t.Error(err)
		}
		if respJSON.Status.Allowed != tst.allowed {
			t.Errorf("Expected status %t, got: %t (%s)", tst.allowed, respJSON.Status.Allowed, reqJSON)
		}
	}
}

func TestCustomVerb(t *testing.T) {
	var saTests = []struct {
		verb     string
		group    string
		resource string
		username string
		allowed  bool
	}{
		{"create", "extensions", "thirdpartyresources", "system:serviceaccount:random:default", true},
	}

	for _, tst := range saTests {
		reqJSON := fmt.Sprintf(`
    {
      "spec":{
        "resourceAttributes":{
          "group":"%s",
          "verb":"%s",
          "resource": "%s"
        },
        "user":"%s"
      }
    }`, tst.group, tst.verb, tst.resource, tst.username)
		result, err := postIndex(reqJSON)

		if err != nil {
			t.Error(err)
		}

		respJSON, err := parseResponse(result)
		if err != nil {
			t.Error(err)
		}
		if respJSON.Status.Allowed != tst.allowed {
			t.Errorf("Expected status %t, got: %t (%s)", tst.allowed, respJSON.Status.Allowed, reqJSON)
		}
	}

}

func TestMultipleGroups(t *testing.T) {
	var saTests = []struct {
		username string
		group    []string
		allowed  bool
	}{
		{"system:serviceaccount:random:default", []string{"badgang", "hackers"}, false},
		{"system:serviceaccount:random:default", []string{"operations", "admins"}, true},
	}

	for _, tst := range saTests {
		j, _ := json.Marshal(tst.group)
		render := string(j)
		reqJSON := fmt.Sprintf(`
    {
      "spec":{
        "resourceAttributes":{
          "verb": "create",
          "resource": "thirdpartyresources"
        },
        "user":"%s",
        "group":  %v
      }
    }`, tst.username, render)
		result, err := postIndex(reqJSON)

		if err != nil {
			t.Error(err)
		}

		respJSON, err := parseResponse(result)
		if err != nil {
			t.Error(err)
		}
		if respJSON.Status.Allowed != tst.allowed {
			t.Errorf("Expected status %d, got: %d (%s)", tst.allowed, respJSON.Status.Allowed, reqJSON)
		}
	}

}

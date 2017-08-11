# kube-system default user can access any namespace
access "allow" {
    username = "system:serviceaccount:kube-system:default"
    verb = "(list|watch|get)"
}

access "deny" {
    username = "badguy"
}

# manage access for users in a single or multiple groups
access "allow" {
   group = "(operations|engineering)"
}

access "deny" {
   group = "(badgang|badgroup)"
}

access "deny" {
   group = "banned"
}

# God mode for regular non-serviceaccount users 
access "allow" {
    username = "[a-z]+"
}

access "allow" {
    group = "{{ .Request.Namespace }}"
}

# default service accounts can create thirdpartyresources
access "allow" {
    username = "system:serviceaccount:.*:default"
    verb = "create"
    resource = "thirdpartyresources"
}

# anyone can access /api and /apis endpoints
access "allow" {
    path = "/api(s*)"
}

# anyone can access /version
access "allow" {
    path = "/version"
}

# read-only for everyone
access "allow" {
    path = "/swaggerapi/api/v1"
    verb = "get"
}

# read-only for everyone
access "allow" {
    path = "/swaggerapi/apis/extensions/v1beta1"
    verb = "get"
}

# This magic allows serviceaccounts to access namespaces
# which have the same prefix as service accounts' originating
# namespace. E.g. namespace-dev can access namespace-prd, namespace-tst  etc.
access "allow" {
    user = "system:serviceaccount:.*:default"
    namespace = "{{ replace .ServiceAccount.Namespace \"-[a-z]{3}\" \"\" }}(-.{3})*"
}

# service account in god mode for it's own namespace 
access "allow" {
    user = "system:serviceaccount:.*:default",
    namespace = "{{ .ServiceAccount.Namespace }}",
}
	

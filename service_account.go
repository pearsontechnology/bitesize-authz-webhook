package main

import "strings"

type ServiceAccount struct {
	User      string
	Namespace string
}

func NewServiceAccount(userString string) *ServiceAccount {
	data := strings.Split(userString, ":")

	if len(data) == 4 &&
		data[0] == "system" &&
		data[1] == "serviceaccount" {
		return &ServiceAccount{
			User:      data[3],
			Namespace: data[2],
		}
	}
	return &ServiceAccount{}
}

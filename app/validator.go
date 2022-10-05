package main

import (
	"net"
	"strings"
)

func (app *application) validateIP(ip string) bool {

	checkedIP := net.ParseIP(ip)

	splittedIP := strings.FieldsFunc(ip, Split)

	key := splittedIP[0] + splittedIP[1] + splittedIP[2]

	for i := range app.IPList[key] {

		if app.IPList[key][i].IP != nil && app.IPList[key][i].IP.Equal(checkedIP) {
			return true
		}

		if app.IPList[key][i].IPNet != nil && app.IPList[key][i].IPNet.Contains(checkedIP) {
			return true
		}
	}
	return false
}

func Split(r rune) bool {
	return r == ':' || r == '.'
}

package utils

import (
	"net"
)

func IsValidIP(ip string) bool {
	res := net.ParseIP(ip)
	return res != nil
}

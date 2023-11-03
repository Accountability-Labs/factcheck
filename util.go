package main

import "strconv"

func isValidPort(portStr string) bool {
	maybePort, err := strconv.ParseInt(portStr, 10, 32)
	if err != nil {
		return false
	}
	return maybePort > 0 && maybePort <= 65535
}

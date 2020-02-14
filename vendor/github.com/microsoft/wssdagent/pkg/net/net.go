// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package net

import (
	"net"
)

func GetIPAddress() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return conn.LocalAddr().(*net.UDPAddr).IP.String(), nil
}

func StringToNetIPAddress(ipString string) net.IP {
	return net.ParseIP(ipString)
}

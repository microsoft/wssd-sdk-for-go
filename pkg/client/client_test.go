package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getServerEndpoint(t *testing.T) {
	tests := []struct {
		name                   string
		serverName             string
		expectedServerName     string
		expectedAuthServerName string
	}{
		{
			name:                   "server with no port",
			serverName:             "myhost",
			expectedServerName:     "myhost:45000",
			expectedAuthServerName: "myhost:45001",
		},
		{
			name:                   "server with port",
			serverName:             "myhost:1234",
			expectedServerName:     "myhost:1234",
			expectedAuthServerName: "myhost:1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualServerName := getServerEndpoint(&tt.serverName)
			assert.Equal(t, tt.expectedServerName, actualServerName)

			actualAuthServerName := getAuthServerEndpoint(&tt.serverName)
			assert.Equal(t, tt.expectedAuthServerName, actualAuthServerName)
		})
	}
}

func Test_getServerAddressOnly(t *testing.T) {
	tests := []struct {
		name               string
		serverName         string
		expectedServerName string
	}{
		{
			name:               "server with no port",
			serverName:         "myhost",
			expectedServerName: "myhost",
		},
		{
			name:               "server with port",
			serverName:         "myhost:1234",
			expectedServerName: "myhost",
		},
		{
			name:               "server with port and whitespace",
			serverName:         " myhost : 1234 ",
			expectedServerName: "myhost",
		},
		{
			name:               "empty string",
			serverName:         "",
			expectedServerName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualServerName := getServerAddressOnly(tt.serverName)
			assert.Equal(t, tt.expectedServerName, actualServerName)
		})
	}
}

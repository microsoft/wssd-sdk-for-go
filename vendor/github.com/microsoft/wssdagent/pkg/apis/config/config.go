// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package config

import (
	"bytes"
	"context"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"strings"

	"github.com/microsoft/wssdagent/pkg/trace"
)

const (
	ListenAddress       string = "0.0.0.0"
	ServerPort          int    = 45000
	DefaultProviderSpec string = "vmms"
)

const environmentProviderSpecStr = "WSSD_PROVIDER_SPEC"

// Returns nil if debug mode is on; err if it is not
func environmentProviderSpec(span *trace.LogSpan) string {
	environmentProvSpec := strings.ToLower(os.Getenv(environmentProviderSpecStr))

	if len(environmentProvSpec) > 0 {
		if environmentProvSpec != "vmms" && environmentProvSpec != "hcs" {
			span.Log("Invalid WSSD_PROVIDER_SPEC: '%s', the only supported values are 'vmms' and 'hcs'", environmentProvSpec)
			return DefaultProviderSpec
		}

		return environmentProvSpec
	}

	return DefaultProviderSpec
}

// DefaultAgentCofiguration
func DefaultAgentConfiguration() *WSSDAgentConfiguration {
	_, span := trace.NewSpan(context.Background(), "Wssdagent DefaultAgentConfiguration Span")

	provSpec := environmentProviderSpec(span)
	basePath := viper.GetString("BaseDir")
	ac := WSSDAgentConfiguration{
		Port:    ServerPort,
		Address: ListenAddress,
		BaseConfiguration: BaseConfiguration{
			DataStorePath:   path.Join(basePath),
			ConfigStorePath: path.Join(basePath),
			LogPath:         path.Join(basePath, "log"),
		},
		ProviderConfigurations: map[string]ChildAgentConfiguration{
			"virtualmachine":          newChildAgentConfiguration(path.Join(basePath), "virtualmachine", provSpec),
			"virtualnetwork":          newChildAgentConfiguration(path.Join(basePath), "virtualnetwork", provSpec),
			"virtualnetworkinterface": newChildAgentConfiguration(path.Join(basePath), "virtualnetworkinterface", provSpec),
			"loadbalancer":            newChildAgentConfiguration(path.Join(basePath), "loadbalancer", provSpec),
		},
		ImageStorePath: "c:/wssdimagestore",
	}
	// Load Default configuration
	agentConfig, err := yaml.Marshal(ac)
	if err != nil {
		panic("Default configuration generation failed")
	}
	err = viper.ReadConfig(bytes.NewBuffer(agentConfig))
	if err != nil {
		panic("Default configuration generation failed")
	}
	return &ac
}

func GetAgentConfiguration() *WSSDAgentConfiguration {
	return viper.Get("wssdagent").(*WSSDAgentConfiguration)
}

func GetChildAgentConfiguration(childAgentName string) *ChildAgentConfiguration {
	_, span := trace.NewSpan(context.Background(), "Wssdagent ChildAgentConfiguration Span")

	// Get Agents configuration
	wssdAgentConfig := GetAgentConfiguration()
	if val, ok := wssdAgentConfig.ProviderConfigurations[childAgentName]; ok {
		return &val
	}

	childAgentConfig := newChildAgentConfiguration(wssdAgentConfig.DataStorePath, childAgentName, environmentProviderSpec(span))

	wssdAgentConfig.ProviderConfigurations[childAgentName] = childAgentConfig

	return &childAgentConfig
}

func GetPublicKeyConfiguration() string {
	basePath := viper.GetString("BaseDir")
	publicKeyName := viper.GetString("PublicKeyName")
	return path.Join(basePath, publicKeyName)
}

func GetTLSServerCertConfiguration() string {
	return viper.GetString("TLSCertPath")
}
func GetTLSServerKeyConfiguration() string {
	return viper.GetString("TLSKeyPath")
}

func newChildAgentConfiguration(dataStorePath string, childAgentName string, providerSpec string) ChildAgentConfiguration {
	return ChildAgentConfiguration{
		BaseConfiguration: BaseConfiguration{
			DataStorePath:   path.Join(dataStorePath, childAgentName),
			ConfigStorePath: path.Join(dataStorePath, childAgentName),
			LogPath:         path.Join(dataStorePath, childAgentName, "log"),
		},
		ProviderSpec: providerSpec,
	}
}

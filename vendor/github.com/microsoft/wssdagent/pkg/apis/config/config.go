// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package config

import (
	"bytes"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"path"
)

const (
	ListenAddress string = "0.0.0.0"
	ServerPort    int    = 45000
)

// DefaultAgentCofiguration
func DefaultAgentConfiguration() *WSSDAgentConfiguration {
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
				"virtualmachine": ChildAgentConfiguration{ProviderSpec: "hcs",},
		},
		ImageStorePath:         "c:/wssdimagestore",
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
	// Get Agents configuration
	wssdAgentConfig := GetAgentConfiguration()
	if val, ok := wssdAgentConfig.ProviderConfigurations[childAgentName]; ok {
		return &val
	}

	childAgentConfig := ChildAgentConfiguration{
		BaseConfiguration: BaseConfiguration{
			DataStorePath:   path.Join(wssdAgentConfig.DataStorePath, childAgentName),
			ConfigStorePath: path.Join(wssdAgentConfig.DataStorePath, childAgentName),
			LogPath:         path.Join(wssdAgentConfig.DataStorePath, childAgentName, "log"),
		},
	}

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
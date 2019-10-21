// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package server

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"

	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/marshal"
	aserver "github.com/microsoft/wssdagent/pkg/server"
)

const (
	WSSDAgent  string = "wssdagent"
	ConfigPath string = "C:/ProgramData/WssdAgent"
)

type Flags struct {
	// ServerName which hosts this virtual machine
	ServerName string
	// LogLevel which hosts this virtual machine
	LogLevel int
	// Verbose mode for debugging
	Verbose bool

	ConfigFilePath string
}

func NewCommand() *cobra.Command {
	flags := &Flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "wssdagent",
		Short: "Wssd node agent",
		Long: ` Wssdagent is a node agent that runs on each node. 
		It provides basic services like compute, network, storage, etc to the users.
		`,
		SilenceUsage: true,
		Version:      "0.01",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	viper.AddConfigPath(ConfigPath)
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	viper.SetDefault("BaseDir", ConfigPath)
	viper.SetDefault("PublicKeyName", "jwt.pem")
	viper.SetDefault("TLSCertPath", path.Join(wd, "wssd.pem"))
	viper.SetDefault("TLSKeyPath", path.Join(wd, "wssd-key.pem"))

	cmd.PersistentFlags().BoolVar(&flags.Verbose, "verbose", false, "Verbose Output")
	cmd.PersistentFlags().IntVar(&flags.LogLevel, "loglevel", 1, "Logging level")
	cmd.Flags().StringVar(&flags.ConfigFilePath, "config", path.Join(ConfigPath, "config.yaml"), "configuration file path")

	return cmd

}

func runE(flags *Flags) error {
	// Make sure the base path exists
	os.MkdirAll(ConfigPath, os.ModeDir)

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			agentConfig := config.DefaultAgentConfiguration()
			// Create the file, if not exists
			viper.Set("wssdagent", agentConfig)
			err = viper.WriteConfigAs(flags.ConfigFilePath)
			if err != nil {
				fmt.Printf("Unable to write configuration [%+v]", err)
				return err
			}

		} else {
			panic("Unable to read configuration")
		}
	} else {
		agentConfigInterface := viper.Get("wssdagent")
		yamlBytes, err := marshal.ToYAMLBytes(agentConfigInterface)
		if err != nil {
			return err
		}

		agentConfig := new(config.WSSDAgentConfiguration)
		err = marshal.FromYAMLBytes(yamlBytes, agentConfig)
		if err != nil {
			return err
		}
		viper.Set("wssdagent", agentConfig)
	}

	return aserver.NewWssdAgentServer()
}

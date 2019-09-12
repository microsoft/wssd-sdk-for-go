// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package server

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssdagent/pkg/wssdagent/apis/config"
	aserver "github.com/microsoft/wssdagent/pkg/wssdagent/server"
)

const (
	WSSDAgent string = "wssdagent"
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
	viper.AddConfigPath("C:/ProgramData/WssdAgent/")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetDefault("BaseDir", "C:/ProgramData/WssdAgent/")

	cmd.PersistentFlags().BoolVar(&flags.Verbose, "verbose", false, "Verbose Output")
	cmd.PersistentFlags().IntVar(&flags.LogLevel, "loglevel", 1, "Logging level")
	cmd.Flags().StringVar(&flags.ConfigFilePath, "config", "", "configuration file path")

	return cmd

}

func runE(flags *Flags) error {
	var err error
	if len(flags.ConfigFilePath) != 0 {
		contents, err := ioutil.ReadFile(flags.ConfigFilePath)
		if err != nil {
			return err
		}
		// a config file is passed in
		err = viper.ReadConfig(bytes.NewBuffer(contents))
	} else {
		err = viper.ReadInConfig()
	}

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			agentConfig := config.DefaultAgentConfiguration()
			viper.WriteConfig()
			viper.Set("wssdagent", agentConfig)
		} else {
			panic(fmt.Sprintf("Failed to load Configuration %v", err))
		}
	}

	return aserver.NewWssdAgentServer()
}

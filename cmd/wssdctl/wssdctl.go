// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/network"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/security"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/storage"
)

type Flags struct {
	// ServerName which hosts this virtual machine
	ServerName string
	// LogLevel which hosts this virtual machine
	LogLevel int
	// Verbose mode for debugging
	Verbose bool
	// Debug mode to disable TLS
	Debug bool
	// Group
	Group string
}

func NewCommand() *cobra.Command {
	flags := &Flags{}
	cmd := &cobra.Command{
		Args:         cobra.NoArgs,
		Use:          "wssdctl",
		Short:        "",
		Long:         "",
		SilenceUsage: true,
		Version:      "0.01",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.PersistentFlags().String("server", "127.0.0.1", "server to which the request has to be sent to")
	viper.BindPFlag("server", cmd.PersistentFlags().Lookup("server"))

	cmd.PersistentFlags().BoolVar(&flags.Verbose, "verbose", false, "Verbose Output")
	cmd.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "Debug Mode to disable TLS")
	viper.BindPFlag("debug", cmd.PersistentFlags().Lookup("debug"))

	cmd.PersistentFlags().IntVar(&flags.LogLevel, "loglevel", 1, "Logging level")

	cmd.PersistentFlags().String("group", "dummpGroup", "Group Name")
	viper.BindPFlag("group", cmd.PersistentFlags().Lookup("group"))

	cmd.PersistentFlags().String("container", "", "Storage Container Name")
	viper.BindPFlag("container", cmd.PersistentFlags().Lookup("container"))

	cmd.AddCommand(compute.NewCommand())
	cmd.AddCommand(network.NewCommand())
	cmd.AddCommand(storage.NewCommand())
	cmd.AddCommand(security.NewCommand())

	return cmd

}

func runE(flags *Flags) error {
	viper.SetDefault("Debug", flags.Debug)
	return nil
}

func Run() error {
	return NewCommand().Execute()
}

func main() {
	// klog.InitFlags(nil)
	//	_ = flag.Set("logtostderr", "false")
	//	_ = flag.Set("logtostderr", "false")
	// flag.Parse()

	if err := Run(); err != nil {
		os.Exit(1)
	}

}

// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package main

import (
	"flag"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"

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
	// Group
	Group string
}

func NewCommand() *cobra.Command {
	flags := Flags{}
	cmd := &cobra.Command{
		Args:         cobra.NoArgs,
		Use:          "wssdctl",
		Short:        "",
		Long:         "",
		SilenceUsage: true,
		Version:      "0.01",
	}

	cmd.PersistentFlags().String("server", "127.0.0.1", "server to which the request has to be sent to")
	viper.BindPFlag("server", cmd.PersistentFlags().Lookup("server"))

	cmd.PersistentFlags().BoolVar(&flags.Verbose, "verbose", false, "Verbose Output")
	cmd.PersistentFlags().IntVar(&flags.LogLevel, "loglevel", 1, "Logging level")

	cmd.PersistentFlags().String("group", "dummpGroup", "Group Name")
	viper.BindPFlag("group", cmd.PersistentFlags().Lookup("group"))

	cmd.AddCommand(compute.NewCommand())
	cmd.AddCommand(network.NewCommand())
	cmd.AddCommand(storage.NewCommand())
	cmd.AddCommand(security.NewCommand())

	return cmd

}

func Run() error {
	return NewCommand().Execute()
}

func main() {
	klog.InitFlags(nil)
	//	_ = flag.Set("logtostderr", "false")
	//	_ = flag.Set("logtostderr", "false")
	flag.Parse()

	if err := Run(); err != nil {
		os.Exit(1)
	}

}

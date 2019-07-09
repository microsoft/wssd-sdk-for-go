// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package get

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/get/all"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/get/virtualmachine"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/get/virtualmachinescaleset"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/get/virtualnetwork"
)

type Format string

const (
	Yaml Format = "yaml"
	Json Format = "json"
	None Format = "none"
)

type GetFlags struct {
	// OutputFormat to display the output yaml/json
	OutputFormat string
}

func NewCommand() *cobra.Command {
	flags := GetFlags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "get",
		Short: "Get a resource",
		Long:  "Get a resource",
	}
	cmd.PersistentFlags().StringVar(&flags.OutputFormat, "output", "yaml", "format to print the output")
	viper.BindPFlag("output", cmd.PersistentFlags().Lookup("output"))

	cmd.AddCommand(all.NewCommand())
	cmd.AddCommand(virtualmachine.NewCommand())
	cmd.AddCommand(virtualmachinescaleset.NewCommand())
	cmd.AddCommand(virtualnetwork.NewCommand())

	return cmd
}

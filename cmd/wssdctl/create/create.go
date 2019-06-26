// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package create

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/create/virtualmachine"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/create/virtualmachinescaleset"
)

type CreateFlags struct {
	FilePath string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "create",
		Short: "Create a resource",
		Long:  "Create a resource",
	}

	cmd.PersistentFlags().StringP("config", "c", "", "configuration file path")
	cmd.MarkFlagRequired("config")
	viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))

	cmd.AddCommand(virtualmachine.NewCommand())
	cmd.AddCommand(virtualmachinescaleset.NewCommand())
	return cmd
}

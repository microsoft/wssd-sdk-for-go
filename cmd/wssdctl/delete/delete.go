// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package delete

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/delete/virtualmachine"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/delete/virtualmachinescaleset"
)

// DeleteFlags
type DeleteFlags struct {
	FilePath string
}

func NewCommand() *cobra.Command {
	flags := DeleteFlags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "delete",
		Short: "delete a resource",
		Long:  "delete a resource",
	}
	cmd.AddCommand(virtualmachine.NewCommand())
	cmd.AddCommand(virtualmachinescaleset.NewCommand())

	cmd.PersistentFlags().StringVar(&flags.FilePath, "config", "", "configuration file path")
	viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config"))

	return cmd
}

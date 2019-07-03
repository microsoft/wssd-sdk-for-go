// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package start

import (
	"github.com/spf13/cobra"
)

type flags struct {
	Name     string
	FilePath string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "start",
		Aliases: []string{"vm"},
		Short:   "Create a Container",
		Long:    "Create a Container",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.FilePath, "config", "", "configuration file path")
	cmd.MarkFlagRequired("config")

	return cmd
}

func runE(flags *flags) error {
	panic("container start not implemented")

	return nil
}

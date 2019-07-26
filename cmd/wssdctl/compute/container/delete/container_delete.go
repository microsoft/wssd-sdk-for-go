// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package delete

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
		Use:     "delete",
		Aliases: []string{"vm"},
		Short:   "delete a Container",
		Long:    "delete a Container",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.FilePath, "config", "", "configuration file path")
	cmd.MarkFlagRequired("config")

	return cmd
}

func runE(flags *flags) error {
	panic("container delete not implemented")

	return nil
}

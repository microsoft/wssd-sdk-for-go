// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package show

import (
	"github.com/spf13/cobra"
)

type flags struct {
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "show",
		Short: "show a specific vhd",
		Long:  "show a specific vhd ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the vhd")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {
	panic("vhd show not implemented")

	return nil
}

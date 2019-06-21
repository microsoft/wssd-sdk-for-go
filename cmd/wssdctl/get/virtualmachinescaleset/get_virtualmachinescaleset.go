// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualmachinescaleset

import (
	"github.com/spf13/cobra"
)

type flags struct {
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualmachinescaleset",
		Aliases: []string{"vmss"},
		Short:   "Get a specific/all Virtual Machine Scale Set(s)",
		Long:    "Get a specific/all Virtual Machine Scale Set(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	return cmd
}

func runE(*flags) error {
	return nil
}

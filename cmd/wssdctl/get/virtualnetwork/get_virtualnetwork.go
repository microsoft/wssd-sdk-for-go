// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualnetwork

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
		Use:     "virtualnetwork",
		Aliases: []string{"vnet"},
		Short:   "Get a specific/all Virtual Network(s)",
		Long:    "Get a specific/all Virtual Network(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	return cmd
}

func runE(*flags) error {
	return nil
}

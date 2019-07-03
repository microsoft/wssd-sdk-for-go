// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list

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
		Use:   "list",
		Short: "list loadbalancer",
		Long:  "list loadbalancer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	return cmd
}

func runE(flags *flags) error {
	panic("lb list not implemented")

	return nil
}

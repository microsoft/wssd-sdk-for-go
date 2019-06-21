// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package loadbalancer

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
		Use:     "loadbalancer",
		Aliases: []string{"elb"},
		Short:   "Get a specific/all LoadBalancer(s)",
		Long:    "Get a specific/all LoadBalancer(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	return cmd
}

func runE(*flags) error {
	return nil
}

// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package all

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
		Use:   "all",
		Short: "Get all Resources)",
		Long:  "Get all Resources managed by Wssd Agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	return cmd
}

func runE(*flags) error {
	return nil
}

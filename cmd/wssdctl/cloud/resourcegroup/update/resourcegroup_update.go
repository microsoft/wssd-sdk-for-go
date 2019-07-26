// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package update

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
		Use:   "update",
		Short: "update a specific resourecegroup",
		Long:  "update a specific resourecegroup ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the resourecegroup")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {
	panic("resourecegroup update not implemented")

	return nil
}

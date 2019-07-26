// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list_ip_addresses

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
		Use:     "list-ip-addresses",
		Aliases: []string{"listips"},
		Short:   "List all IPs of a Virtual Machine",
		Long:    "List all IPs of a Virtual Machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual machine")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {

	panic("vm list ips not implemented")
	return nil
}

// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualmachine

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine"
)

type flags struct {
	// Name of the Virtual Machine to get
	Name string
	// ServerName which hosts this virtual machine
	ServerName string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualmachine",
		Aliases: []string{"vm"},
		Short:   "Delete a specific/all Virtual Machine(s)",
		Long:    "Delete a specific/all Virtual Machine(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name(s) of the virtual machine, comma separated")
	cmd.Flags().StringVar(&flags.ServerName, "server", "", "server to which the request has to be sent to")

	return cmd
}

func runE(flags *flags) error {
	vmclient, err := virtualmachine.NewVirtualMachineClient(flags.ServerName)
	if err != nil {
		return err
	}

	err = vmclient.Delete(nil, "", flags.Name)
	if err != nil {
		return err
	}

	return nil
}

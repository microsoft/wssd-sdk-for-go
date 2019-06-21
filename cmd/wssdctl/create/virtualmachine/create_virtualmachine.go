// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package create

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine"
)

type flags struct {
	ServerName string
	FilePath   string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualmachine",
		Aliases: []string{"vm"},
		Short:   "Create a Virtual Machine",
		Long:    "Create a Virtual Machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.FilePath, "f", "", "configuration file path")
	cmd.Flags().StringVar(&flags.ServerName, "server", "", "server to which the request has to be sent to")

	return cmd
}

func runE(flags *flags) error {
	vmconfig, err := virtualmachine.LoadConfig(flags.FilePath)
	if err != nil {
		return err
	}

	vmclient, err := virtualmachine.NewVirtualMachineClient(flags.ServerName)
	if err != nil {
		return err
	}

	_, err = vmclient.CreateOrUpdate(nil, "", "", vmconfig)
	if err != nil {
		return err
	}

	return nil
}

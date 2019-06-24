// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualmachine

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine"
)

type flags struct {
	// Name of the Virtual Machine to get
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "virtualmachine",
		Aliases: []string{"vm"},
		Short:   "Delete a specific VirtualMachine",
		Long:    "Delete a specific VirtualMachine",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name(s) of the virtual machine, comma separated")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	vmclient, err := virtualmachine.NewVirtualMachineClient(server)
	if err != nil {
		return err
	}

	vmName := flags.Name
	vmId := ""
	if len(vmName) == 0 {
		config := viper.GetString("config")
		vmconfig, err := virtualmachine.LoadConfig(config)
		if err != nil {
			return err
		}
		vmName = *(vmconfig.Name)
		vmId = *(vmconfig.ID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = vmclient.Delete(ctx, vmId, vmId)
	if err != nil {
		return err
	}

	return nil
}

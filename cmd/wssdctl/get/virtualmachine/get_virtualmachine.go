// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package virtualmachine

import (
	"context"
	"fmt"
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
		Short:   "Get a specific/all Virtual Machine(s)",
		Long:    "Get a specific/all Virtual Machine(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual machine")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	vmclient, err := virtualmachine.NewVirtualMachineClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(flags.Name) == 0 {
		vms, err := vmclient.List(ctx)
		if err != nil {
			return err
		}
		if vms == nil || len(*vms) == 0 {
			fmt.Println("No VirtualMachine Resources")
			return nil
		}
		virtualmachine.PrintList(vms)
	} else {
		vm, err := vmclient.Get(ctx, flags.Name)
		if err != nil {
			return err
		}
		virtualmachine.Print(vm)
	}

	return nil
}

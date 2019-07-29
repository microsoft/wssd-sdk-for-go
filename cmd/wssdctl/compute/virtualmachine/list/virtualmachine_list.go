// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachine"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "list",
		Short: "Get a specific/all Virtual Machine(s)",
		Long:  "Get a specific/all Virtual Machine(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	vmclient, err := virtualmachine.NewVirtualMachineClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vms, err := vmclient.Get(ctx, "")
	if err != nil {
		return err
	}
	if vms == nil || len(*vms) == 0 {
		fmt.Println("No VirtualMachine Resources")
		// Not an error
		return nil
	}

	virtualmachine.PrintList(vms)

	return nil
}

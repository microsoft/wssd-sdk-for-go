// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package scale

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachinescaleset"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name     string
	Replicas string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "scale",
		Short: "Reimage Virtual machines within the specified Virtual Machine Scale Set",
		Long:  "Reimage Virtual machines within the specified Virtual Machine Scale Set",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual machine scale set resource")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&flags.Replicas, "count", "", "number of replicas to scale the virtual machine scale set resource")
	cmd.MarkFlagRequired("count")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	group := viper.GetString("group")
	client, err := virtualmachinescaleset.NewVirtualMachineScaleSetClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vmss, err := client.Get(ctx, group, flags.Name)
	if err != nil {
		return err
	}
	if vmss == nil || len(*vmss) == 0 {
		return fmt.Errorf("Unable to find Virtual Machine Scale Set [%s]", flags.Name)
	}

	panic("vmss scale not implemented")

	return nil

}

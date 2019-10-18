// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package create

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/wssd-sdk-for-go/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/compute"
	"github.com/microsoft/wssd-sdk-for-go/services/compute/virtualmachinescaleset"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name     string
	FilePath string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "create",
		Short: "Create a Virtual Machine Scale Set",
		Long:  "Create a Virtual Machine Scale Set",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.FilePath, "config", "", "configuration file path")
	cmd.MarkFlagRequired("config")

	return cmd
}

func runE(flags *flags) error {
	server := viper.GetString("server")
	group := viper.GetString("group")
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}
	client, err := virtualmachinescaleset.NewVirtualMachineScaleSetClient(server, authorizer)
	if err != nil {
		return err
	}

	vmconfig := compute.VirtualMachineScaleSet{}
	err = config.LoadYAMLFile(flags.FilePath, &vmconfig)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	_, err = client.CreateOrUpdate(ctx, group, *(vmconfig.Name), &vmconfig)
	if err != nil {
		return err
	}

	return nil
}

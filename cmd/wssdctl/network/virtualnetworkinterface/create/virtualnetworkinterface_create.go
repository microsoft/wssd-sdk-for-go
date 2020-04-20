// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package create

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/config"
	"github.com/microsoft/moc/pkg/errors"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetworkinterface"

	wssdcommon "github.com/microsoft/moc/common"
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
		Short: "Create a Virtual Network interface",
		Long:  "Create a Virtual Network interface",
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

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	vnicclient, err := virtualnetworkinterface.NewVirtualNetworkInterfaceClient(server, authorizer)
	if err != nil {
		return err
	}

	vnicconfig := network.VirtualNetworkInterface{}
	err = config.LoadYAMLFile(flags.FilePath, &vnicconfig)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	if vnicconfig.Name == nil {
		return errors.Wrapf(errors.InvalidInput, "The YAML is missing the 'Name' element")
	}

	_, err = vnicclient.CreateOrUpdate(ctx, group, *(vnicconfig.Name), &vnicconfig)
	if err != nil {
		return err
	}

	return nil
}

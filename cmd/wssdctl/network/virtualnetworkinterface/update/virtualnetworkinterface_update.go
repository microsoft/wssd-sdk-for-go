// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package update

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/network/virtualnetworkinterface"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
		Use:   "update",
		Short: "update a Virtual Network interface",
		Long:  "update a Virtual Network interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual network interface resource(s), comma separated")
	cmd.MarkFlagRequired("name")

	return cmd
}

func runE(flags *flags) error {

	group := viper.GetString("group")
	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	vnicclient, err := virtualnetworkinterface.NewVirtualNetworkInterfaceClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()
	panic("vnic update not implemented")

	_, err = vnicclient.CreateOrUpdate(ctx, group, flags.Name, nil)
	if err != nil {
		return err
	}

	return nil
}

// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package delete

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/services/network"
	"github.com/microsoft/wssd-sdk-for-go/services/network/loadbalancer"

	wssdcommon "github.com/microsoft/moc/common"
)

type flags struct {
	Name string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "delete",
		Short: "Delete a specific loadbalancer",
		Long:  "Delete a specific loadbalancer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the virtual network resource(s), comma separated")
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

	lbclient, err := loadbalancer.NewLoadBalancerClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	lbName := flags.Name
	if len(lbName) == 0 {
		configPath := viper.GetString("config")

		lbconfig := network.VirtualNetworkInterface{}
		err = config.LoadYAMLFile(configPath, &lbconfig)
		if err != nil {
			return err
		}
		lbName = *(lbconfig.Name)
	}

	err = lbclient.Delete(ctx, group, lbName)
	if err != nil {
		return err
	}

	return nil
}

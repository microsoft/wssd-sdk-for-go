// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package create

import (
	"context"
	//"time"
	"fmt"

	wssdcommon "github.com/microsoft/moc/common"
	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/services/security"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Short: "Create a keyvault",
		Long:  "Create a keyvault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the keyvault resource(s), comma separated")
	cmd.Flags().StringVar(&flags.FilePath, "config", "", "configuration file path")

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")

	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	vaultClient, err := keyvault.NewKeyVaultClient(server, authorizer)
	if err != nil {
		return err
	}

	var vaultName string
	var kvConfig *security.KeyVault

	if flags.FilePath != "" {
		kvConfig = &security.KeyVault{}
		err = config.LoadYAMLFile(flags.FilePath, kvConfig)
		if err != nil {
			return err
		}
		vaultName = *kvConfig.Name
	} else {
		if flags.Name == "" {
			return fmt.Errorf("Error: must specify --name or --config")
		}
		kvConfig = nil
		vaultName = flags.Name
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	_, err = vaultClient.CreateOrUpdate(ctx, group, vaultName, kvConfig)
	if err != nil {
		return err
	}

	return nil
}

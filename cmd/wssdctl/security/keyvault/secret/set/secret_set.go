// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package set

import (
	"context"
	//"time"
	"fmt"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
	"github.com/microsoft/wssd-sdk-for-go/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault/secret"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type flags struct {
	Name           string
	FilePathConfig string
	FilePathValue  string
	Value          string
	VaultName      string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "set",
		Short: "set a secret",
		Long:  "set a secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the secret, comma separated")
	cmd.Flags().StringVar(&flags.Value, "value", "", "name of the secret, comma separated")
	cmd.Flags().StringVar(&flags.VaultName, "vault-name", "", "name of the secret, comma separated")
	cmd.Flags().StringVar(&flags.FilePathConfig, "config", "", "configuration file path")
	cmd.Flags().StringVar(&flags.FilePathValue, "file", "", "value file path")

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")
	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return err
	}

	secretClient, err := secret.NewSecretClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	var secretName string
	var srtConfig *keyvault.Secret
	if flags.FilePathConfig != "" {
		srtConfig = &keyvault.Secret{}
		err = config.LoadYAMLFile(flags.FilePathConfig, srtConfig)
		if err != nil {
			return err
		}

		secretName = *srtConfig.Name
	} else {
		if flags.Name == "" || flags.VaultName == "" {
			return fmt.Errorf("Error: must specify --config or --name, --vault-name")
		}
		var value *string
		if flags.FilePathValue != "" {
			value, err = config.LoadValueFile(flags.FilePathValue)
			if err != nil {
				return err
			}
		} else if flags.Value != "" {
			value = &flags.Value
		} else {
			return fmt.Errorf("Error: must specify --value or --file-path")
		}

		srtConfig = &keyvault.Secret{
			Name:  &flags.Name,
			Value: value,
			SecretProperties: &keyvault.SecretProperties{
				VaultName: &flags.VaultName,
			},
		}
		secretName = flags.Name
	}

	_, err = secretClient.CreateOrUpdate(ctx, group, secretName, srtConfig)
	if err != nil {
		return err
	}

	return nil
}

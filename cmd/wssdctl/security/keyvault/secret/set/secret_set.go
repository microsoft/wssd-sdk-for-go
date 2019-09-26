// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package set

import (
	"context"
	//"time"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/microsoft/wssd-sdk-for-go/services/security"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault/secret"
	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
)

type flags struct {
	Name      string
	FilePath  string
	Value     string
	VaultName string
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
	cmd.Flags().StringVar(&flags.FilePath, "config", "", "configuration file path")

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")

	server := viper.GetString("server")
	secretClient, err := secret.NewSecretClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	var secretName string
	var srtConfig *keyvault.Secret
	if flags.FilePath != "" {
		config := flags.FilePath
		srtConfig, err := secret.LoadConfig(config)
		if err != nil {
			return err
		}

		secretName = *srtConfig.Name 
	} else {
		if flags.Name == "" || flags.Value == "" || flags.VaultName == "" {
			return fmt.Errorf("Error: must specify --config or --name, --vault-name, --value")
		}

		srtConfig = &keyvault.Secret{
			BaseProperties: security.BaseProperties{
				Name: &flags.Name,
			},
			Value : &flags.Value,
			VaultName: &flags.VaultName,
		}
		secretName = flags.Name
	}
 
	_, err = secretClient.CreateOrUpdate(ctx, group, secretName, srtConfig)
	if err != nil {
		return err
	}

	return nil
}

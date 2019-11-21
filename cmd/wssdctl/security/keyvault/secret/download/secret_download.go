// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package download

import (
	"context"
	//"time"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	wssdcommon "github.com/microsoft/wssd-sdk-for-go/common"
	"github.com/microsoft/wssd-sdk-for-go/pkg/auth"
	"github.com/microsoft/wssd-sdk-for-go/pkg/config"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault/secret"
)

type flags struct {
	Name      string
	FilePath  string
	VaultName string
	Output    string
	Query     string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "download",
		Short: "download a secret",
		Long:  "download a secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the secret, comma separated")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&flags.VaultName, "vault-name", "", "name of the secret, comma separated")
	cmd.MarkFlagRequired("vault-name")
	cmd.Flags().StringVar(&flags.FilePath, "file-path", "", "name of the secret, comma separated")
	cmd.MarkFlagRequired("file-path")

	cmd.Flags().StringVar(&flags.Output, "output", "yaml", "Output Format")
	cmd.Flags().StringVar(&flags.Query, "query", "", "Output Format")

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

	secrets, err := secretClient.Get(ctx, group, flags.Name, flags.VaultName)
	if err != nil {
		return err
	}

	if secrets == nil || len(*secrets) == 0 {
		fmt.Println("No Secrets")
		// Not an error
		return nil
	}

	err = config.ExportFormatList(secrets, flags.FilePath, flags.Query, flags.Output)

	return err
}

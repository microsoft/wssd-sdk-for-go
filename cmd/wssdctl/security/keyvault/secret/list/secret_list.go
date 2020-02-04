// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list

import (
	"context"
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
	VaultName string
	FilePath  string
	Output    string
	Query     string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "list",
		Short: "list a secret",
		Long:  "list a secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}
	cmd.Flags().StringVar(&flags.Output, "output", "yaml", "Output Format [yaml, json, csv, tsv]")
	cmd.Flags().StringVar(&flags.Query, "query", "", "Output Format")
	cmd.Flags().StringVar(&flags.VaultName, "vault-name", "", "name of the vault")
	cmd.MarkFlagRequired("vault-name")

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")
	server := viper.GetString("server")

	authorizer, err := auth.NewAuthorizerFromEnvironment(server)
	if err != nil {
		return err
	}

	vaultClient, err := secret.NewSecretClient(server, authorizer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	secrets, err := vaultClient.Get(ctx, group, flags.Name, flags.VaultName)
	if err != nil {
		return err
	}

	if secrets == nil || len(*secrets) == 0 {
		fmt.Println("No Key Vaults")
		// Not an error
		return nil
	}

	config.PrintFormatList(*secrets, flags.Query, flags.Output)
	return nil
}

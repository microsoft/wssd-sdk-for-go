// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package list

import (
	"context"
	//"time"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//"github.com/microsoft/wssd-sdk-for-go/services/keyvault"
	"github.com/microsoft/wssd-sdk-for-go/services/keyvault/simplevault"
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
		Use:   "list",
		Short: "list a simplevault",
		Long:  "list a simplevault",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	return cmd
}

func runE(flags *flags) error {
	group := viper.GetString("group")

	server := viper.GetString("server")
	vaultClient, err := simplevault.NewSimpleVaultClient(server)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wssdcommon.DefaultServerContextTimeout)
	defer cancel()

	vaultName := flags.Name

	keyvaults, err := vaultClient.Get(ctx, group, vaultName)
	if err != nil {
		return err
	}
	
	if keyvaults == nil || len(*keyvaults) == 0 {
		fmt.Println("No Key Vaults")
		// Not an error
		return nil           
	}

	simplevault.PrintList(keyvaults)

	return nil
}
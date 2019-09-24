// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package keyvault

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/keyvault/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/keyvault/delete"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/keyvault/list"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/keyvault/secret"
)

type GetFlags struct {
	// OutputFormat to display the output yaml/json
	OutputFormat string
	Group        string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "keyvault",
		Short: "keyvault resource",
		Long:  "keyvault resource",
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(delete.NewCommand())
	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(secret.NewCommand())

	return cmd
}
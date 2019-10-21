// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package security

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/security/keyvault"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/security/identity"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/security/login"
)

type GetFlags struct {
	// OutputFormat to display the output yaml/json
	OutputFormat string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "security",
		Short: "security resource",
		Long:  "security resource",
	}

	cmd.AddCommand(keyvault.NewCommand())
	cmd.AddCommand(identity.NewCommand())
	cmd.AddCommand(login.NewCommand())
	
	return cmd
}

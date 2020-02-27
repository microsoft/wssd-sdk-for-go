// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package certificate

import (
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/security/certificate/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/security/certificate/list"
	"github.com/spf13/cobra"
)

type GetFlags struct {
	// OutputFormat to display the output yaml/json
	OutputFormat string
	Group        string
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "certificate",
		Short: "certificate resource",
		Long:  "certificate resource",
	}

	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(create.NewCommand())

	return cmd
}

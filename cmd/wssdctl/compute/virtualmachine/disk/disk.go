// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package disk

import (
	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachine/disk/attach"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/compute/virtualmachine/disk/detach"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:    cobra.NoArgs,
		Use:     "disk",
		Aliases: []string{"disk"},
		Short:   "data disk resource",
		Long:    "data disk resource",
	}

	cmd.AddCommand(attach.NewCommand())
	cmd.AddCommand(detach.NewCommand())

	return cmd
}

// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/create"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/delete"
	"github.com/microsoft/wssd-sdk-for-go/cmd/wssdctl/get"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Args:         cobra.NoArgs,
		Use:          "wssdctl",
		Short:        "",
		Long:         "",
		SilenceUsage: true,
		Version:      "0.01",
	}

	cmd.AddCommand(create.NewCommand())
	cmd.AddCommand(get.NewCommand())
	cmd.AddCommand(delete.NewCommand())

	return cmd

}

func Run() error {
	return NewCommand().Execute()
}

func main() {
	if err := Run(); err != nil {
		os.Exit(1)
	}

}

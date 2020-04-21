// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package show

import (
	"github.com/spf13/cobra"
)

type flags struct {
	Name   string
	Output string
	Query  string
}

func NewCommand() *cobra.Command {
	flags := &flags{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "show",
		Short: "show a specific container",
		Long:  "show a specific container ",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(flags)
		},
	}

	cmd.Flags().StringVar(&flags.Name, "name", "", "name of the container")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&flags.Output, "output", "o", "yaml", "Output Format [yaml, json, csv, tsv]")
	cmd.Flags().StringVarP(&flags.Query, "query", "q", "", "Output Format")

	return cmd
}

func runE(flags *flags) error {
	panic("container show not implemented")

	return nil
}

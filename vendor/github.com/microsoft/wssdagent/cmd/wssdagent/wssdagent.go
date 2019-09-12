// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package main

import (
	"flag"
	"fmt"
	"os"

	"k8s.io/klog"

	"github.com/microsoft/wssdagent/cmd/wssdagent/server"
)

func Run() error {
	return server.NewCommand().Execute()
}

func main() {
	klog.InitFlags(nil)
	//	_ = flag.Set("logtostderr", "false")
	//	_ = flag.Set("logtostderr", "false")
	flag.Parse()

	if err := Run(); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

}

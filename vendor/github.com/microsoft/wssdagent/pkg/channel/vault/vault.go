// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package vault

import (
	"github.com/microsoft/wssdagent/pkg/channel"
)

type Notification struct {
	Name      string
	Operation channel.OperationType
}

// Group Channel
type Channel struct {
	// Notify Send channel to send vault Notification
	Notify chan Notification
	// Result channel to get the Result
	Result chan int
}

func MakeChannel() Channel {
	return Channel{Notify: make(chan Notification), Result: make(chan int)}
}

func MakeNotificationData(vaultName string, operation channel.OperationType) Notification {
	return Notification{Name: vaultName, Operation: operation}
}

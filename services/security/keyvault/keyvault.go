// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package keyvault

import (
	"github.com/microsoft/wssd-sdk-for-go/services/security"
)

// Secret defines the structure of a secret
type Secret struct {
	security.BaseProperties
	// KeyValues
	VaultName *string `json:"vaultname"`
	Value     *string `json:"value"`
	FileName  *string `json:"filename"`
}

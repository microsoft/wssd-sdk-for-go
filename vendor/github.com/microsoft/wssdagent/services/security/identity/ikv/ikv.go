// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package ikv
// ikv stands for InternalKeyVault
import (
	"context"
	"time"
	"github.com/microsoft/wssdagent/pkg/auth"
	pb "github.com/microsoft/wssdagent/rpc/security"
	"github.com/microsoft/wssdagent/services/security/identity/internal"
	"github.com/microsoft/wssdagent/services/security/keyvault/secret"
	"github.com/microsoft/wssdagent/services/security/keyvault"

)

const IdentityVaultName = "INTERNAL_IDENTITY_VAULT"

type Client struct {
	secretProvider *secret.SecretProvider
}

func NewClient() *Client {
	identityKeyVaultInit(keyvault.GetKeyVaultProvider())
	return &Client{
		secretProvider: secret.GetSecretProvider(),
	}
}

// Create a Identity
func (c *Client) CreateIdentity(ctx context.Context, identityInternal *internal.IdentityInternal) (err error) {
	ident := identityInternal.Entity

	secretValue, err := auth.GeneratePrivateKey()
	if err != nil {
		return 
	}
	secretListForAPICall := []*pb.Secret{&pb.Secret{
		Name:      ident.Name,
		VaultName: IdentityVaultName,
		Value:     secretValue,
	}}
	_, err = c.secretProvider.CreateOrUpdate(ctx, secretListForAPICall)

	if err != nil {
		return 
	}

	return
}

// Delete a Identity
func (c *Client) CleanupIdentity(ctx context.Context, identityInternal *internal.IdentityInternal) (err error) {
	identityToBeDeleted := identityInternal.Entity

	secretListForAPICall := []*pb.Secret{&pb.Secret{
		Name:      identityToBeDeleted.Name,
		VaultName: IdentityVaultName,
	}}

	err = c.secretProvider.Delete(ctx, secretListForAPICall)
	if err != nil {
		return 
	}
	return
}


func identityKeyVaultInit(keyvaultProvider *keyvault.KeyVaultProvider) {
	var keyvaultListForAPICall []*pb.KeyVault
	keyvaultListForAPICall = append(keyvaultListForAPICall, &pb.KeyVault{
		Name: IdentityVaultName,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()


	vaults, err := keyvaultProvider.Get(ctx, keyvaultListForAPICall)
	if err == nil && len(vaults) != 0 {
		return
	}

	// Identity Vault does not exist .. lets create it
	_, err = keyvaultProvider.CreateOrUpdate(ctx, keyvaultListForAPICall)
	if err != nil {
		panic("failed to create the identity vault")
	}
}
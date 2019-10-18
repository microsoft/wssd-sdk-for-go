package ikv
// ikv stands for InternalKeyVault
import (
	"fmt"
	"context"
	"github.com/microsoft/wssdagent/pkg/auth"
	pb "github.com/microsoft/wssdagent/rpc/security"
	"github.com/microsoft/wssdagent/services/security/keyvault/secret"
	"github.com/microsoft/wssdagent/pkg/apis/config"
)

const IdentityVaultName = "INTERNAL_IDENTITY_VAULT"

type Client struct {
	secretProvider *secret.SecretProvider
}

func NewClient() *Client {
	return &Client{
		secretProvider: secret.GetSecretProvider(),
	}
}

// Create an Identity
func (c *Client) Login(ctx context.Context, id *pb.Identity, authorizer *auth.JwtAuthorizer) (*string, error) {

	secretListForAPICall := []*pb.Secret{&pb.Secret{
		Name:      id.Name,
		VaultName: IdentityVaultName,
	}}
	
	identitySecret, err := c.secretProvider.Get(ctx, secretListForAPICall)
	if err != nil {
		return nil, fmt.Errorf("Error Retrieving Identity: %v", err)
	}

	privateKey := identitySecret[0].Value

	jwtSigner, err := auth.NewJwtSigner(privateKey)
	if err != nil {
		return nil, err
	}

	token, err := jwtSigner.IssueJWT(id.Name, id.Id)
	if err != nil {
		return nil, err
	}

	authorizer.UpdatePublicKey(jwtSigner.GetPublicKey())
	
	// So that the correct public key is picked up on reboot or crash
	err = authorizer.WritePublicKeyToPem(config.GetPublicKeyConfiguration())

	if err != nil {
		return nil, err
	}

	return &token, nil
}
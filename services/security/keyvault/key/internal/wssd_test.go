package internal

import (
	"testing"

	"github.com/microsoft/moc/rpc/common"
	"github.com/microsoft/wssd-sdk-for-go/services/security/keyvault"
	"github.com/stretchr/testify/assert"
)

func StringPtr(s string) *string {
	return &s
}

func Test_KeyOperation(t *testing.T) {
	testKeyType := keyvault.AES
	testKey := keyvault.Key{
		Name:      StringPtr("testKeyName"),
		ID:        StringPtr("testKeyID"),
		VaultName: StringPtr("testVault"),
		KeyType:   &testKeyType,
		KeySize:   256}

	alg := keyvault.A256KW

	t.Run("Test rotate marshaling", func(t *testing.T) {
		keyReq := keyvault.KeyOperationRequest{
			Key:       &testKey,
			Algorithm: &alg,
			Data:      StringPtr("test data")}
		res, err := getKeyOperationRequest(&keyReq, common.ProviderAccessOperation_Key_Rotate)
		assert.Nil(t, err)
		assert.Equal(t, res.Key.Name, "testKeyName")
		assert.Equal(t, res.Key.VaultName, "testVault")
		assert.Equal(t, res.Key.Type, common.JsonWebKeyType_AES)
		assert.Equal(t, res.Data, "") // No data during rotate
	})

	t.Run("Test wrap marshaling", func(t *testing.T) {
		keyReq := keyvault.KeyOperationRequest{
			Key:       &testKey,
			Algorithm: &alg,
			Data:      StringPtr("test data")}
		res, err := getKeyOperationRequest(&keyReq, common.ProviderAccessOperation_Key_WrapKey)
		assert.Nil(t, err)
		assert.Equal(t, res.Key.Name, "testKeyName")
		assert.Equal(t, res.Key.VaultName, "testVault")
		assert.Equal(t, res.Key.Type, common.JsonWebKeyType_AES)
		assert.Equal(t, res.Data, "test data")
	})

	t.Run("Test error", func(t *testing.T) {
		keyReq := keyvault.KeyOperationRequest{
			Key:       &testKey,
			Algorithm: &alg,
			Data:      StringPtr("test data")}
		res, err := getKeyOperationRequest(&keyReq, common.ProviderAccessOperation_Authentication_Login)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "Invalid Input")
	})
}

func Test_Key(t *testing.T) {
	// Make a keyvault.Key
	testKeyType := keyvault.AES
	testKey := keyvault.Key{
		Name:      StringPtr("testKeyName"),
		ID:        StringPtr("testKeyID"),
		VaultName: StringPtr("testVault"),
		KeyType:   &testKeyType,
		KeySize:   256}

	// Convert to wssd.key
	wssdKey := getWssdKey(&testKey)
	assert.Equal(t, *testKey.Name, wssdKey.Name)
	assert.Equal(t, *testKey.VaultName, wssdKey.VaultName)
	assert.Equal(t, common.KeySize__256, wssdKey.Size)

	// Convert back to keyvault.Key
	convertedKey := getKey(wssdKey)
	assert.Equal(t, testKey.Name, convertedKey.Name)
	assert.Equal(t, testKey.VaultName, convertedKey.VaultName)
	assert.Equal(t, testKey.KeySize, convertedKey.KeySize)
}

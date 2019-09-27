package fs

import (
	"encoding/json"
	"fmt"
	"github.com/microsoft/wssdagent/pkg/apis/config"
	"github.com/microsoft/wssdagent/pkg/errors"
	"github.com/microsoft/wssdagent/pkg/crypto"


	pb "github.com/microsoft/wssdagent/rpc/security"
	"io/ioutil"
	"path/filepath"
	"sync"
)

const VaultFileName = "vault.json"

type filesystemvaultmanager struct {
	FilePath string
	mux      sync.Mutex
}

type KeyVaultFileManager struct {
	Vaults map[string]*pb.KeyVault
}

var instance *filesystemvaultmanager
var once sync.Once

func GetInstance() *filesystemvaultmanager {
	once.Do(func() {
		wssdAgentConfig := config.GetAgentConfiguration()
		instance = &filesystemvaultmanager{FilePath: wssdAgentConfig.DataStorePath}

		_, err := instance.openVaultManager()
		if err == nil {
			return
		}

		// Create vault
		rootManager, err := json.Marshal(KeyVaultFileManager{Vaults: make(map[string]*pb.KeyVault)})
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(
			filepath.Join(instance.FilePath, VaultFileName),
			rootManager,
			0644)
		if err != nil {
			panic(err)
		}
	})
	return instance
}

// Must be under lock
func (fsv *filesystemvaultmanager) openVaultManager() (*KeyVaultFileManager, error) {
	data, err := ioutil.ReadFile(filepath.Join(fsv.FilePath, VaultFileName))
	if err != nil {
		return nil, err
	}

	var vaultManager KeyVaultFileManager
	err = json.Unmarshal(data, &vaultManager)
	if err != nil {
		return nil, err
	}

	return &vaultManager, err
}

// Must be under lock
func (fsv *filesystemvaultmanager) closeVaultManager(vaultManager KeyVaultFileManager) error {
	manager, err := json.Marshal(vaultManager)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(
		filepath.Join(fsv.FilePath, VaultFileName),
		manager,
		0644)
}

func (fsv *filesystemvaultmanager) AddVault(kv pb.KeyVault) error {
	fsv.mux.Lock()
	defer fsv.mux.Unlock()

	vaultManager, err := fsv.openVaultManager()
	if err != nil {
		return err
	}
	// Add the new vault
	vaultManager.Vaults[kv.Id] = &kv

	// Save back to file
	err = fsv.closeVaultManager(*vaultManager)
	if err != nil {
		return err
	}

	return nil
}

func (fsv *filesystemvaultmanager) RemoveVault(kv pb.KeyVault) error {
	fsv.mux.Lock()
	defer fsv.mux.Unlock()

	vaultManager, err := fsv.openVaultManager()
	if err != nil {
		return err
	}

	// Remove from Map
	delete(vaultManager.Vaults, kv.Id)

	// Save back to file
	err = fsv.closeVaultManager(*vaultManager)
	if err != nil {
		return err
	}

	return nil
}

func (fsv *filesystemvaultmanager) AddSecretToVault(sec pb.Secret) error {
	fsv.mux.Lock()
	defer fsv.mux.Unlock()

	vaultManager, err := fsv.openVaultManager()
	if err != nil {
		return fmt.Errorf("Failed to open Vault Manager, err: %v", err)
	}

	vault := vaultManager.Vaults[sec.VaultId]

	if vault == nil {
		return fmt.Errorf("Vault Not Found, vault Id: %v", sec.VaultId)
	}

	encryptedValue, err := crypto.EncryptSecret(sec.Value)
	if err != nil {
		return err
	}

	sec.Value = *encryptedValue

	vault.Secrets = append(vault.Secrets, &sec)
	vaultManager.Vaults[sec.VaultId] = vault

	err = fsv.closeVaultManager(*vaultManager)
	if err != nil {
		return err
	}

	return nil
}

func (fsv *filesystemvaultmanager) RemoveSecretFromVault(sec pb.Secret) error {
	fsv.mux.Lock()
	defer fsv.mux.Unlock()

	vaultManager, err := fsv.openVaultManager()
	if err != nil {
		return fmt.Errorf("Failed to open Vault Manager, err: %v", err)
	}

	vault := vaultManager.Vaults[sec.VaultId]

	for index, srt := range vault.Secrets {
		if srt.Id == sec.Id {
			// Remove Secret
			vault.Secrets = append(vault.Secrets[:index], vault.Secrets[index+1:]...)
			break
		}
	}
	vaultManager.Vaults[sec.VaultId] = vault

	err = fsv.closeVaultManager(*vaultManager)
	if err != nil {
		return err
	}

	return nil
}

func (fsv *filesystemvaultmanager) ShowSecretFromVault(sec pb.Secret) (*pb.Secret, error) {
	fsv.mux.Lock()
	defer fsv.mux.Unlock()
	vaultManager, err := instance.openVaultManager()
	if err != nil {
		return nil, fmt.Errorf("Failed to open Vault Manager, err: %v", err)
	}
	if len(sec.Name) == 0 {
		return nil, errors.Wrapf(errors.InvalidInput, "Missing Secret Name")
	}

	vault := vaultManager.Vaults[sec.VaultId]

	for _, srt := range vault.Secrets {
		if srt.Name == sec.Name {
			decryptedValue, err := crypto.DecryptSecret(srt.Value)
			if err != nil {
				return nil, err
			}

			srt.Value = *decryptedValue
			return srt, nil
		}
	}

	return nil, errors.NotFound
}

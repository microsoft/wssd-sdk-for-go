// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package cloudinit

import (
	"encoding/base64"
	"encoding/json"
	"github.com/kdomanski/iso9660"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type group struct {
	names []string
}
type SSHKeys struct {
	RSAPrivateKey string `yaml:"rsa_private,omitempty"`
	RSAPublicKey  string `yaml:"rsa_public,omitempty"`
	DSAPrivateKey string `yaml:"dsa_private,omitempty"`
	DSAPublicKey  string `yaml:"dsa_public,omitempty"`
}

type user struct {
	Name              string   `yaml:"name,omitempty"`
	Gecos             string   `yaml:"gecos,omitempty"`
	Groups            []string `yaml:"groups,omitempty,flow"`
	Password          string   `yaml:"passwd,omitempty"` // Base64 Encoded string
	Sudo              string   `yaml:"sudo,omitempty"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	Keys              SSHKeys  `yaml:"ssh_keys,omitempty"`
}

type WriteFile struct {
	Encoding    string `yaml:"encoding,omitempty"`
	Content     string `yaml:"content,omitempty"`
	Owner       string `yaml:"owner,omitempty"`
	Path        string `yaml:"path,omitempty"`
	Permissions string `yaml:"permission,omitempty"`
}

// Metadata
type Metadata struct {
	InstanceID string `json:"instance-id,omitempty"`
	Hostname   string `json:"local_hostname,omitempty"`
}

// Userdata
type Userdata struct {
	Hostname                  string       `yaml:"hostname,omitempty"`
	ResizeRootFS              bool         `yaml:"resize_rootfs,omitempty"`
	Users                     []*user      `yaml:"users,omitempty"`
	SSHPasswordAuthentication bool         `yaml:"ssh_pwauth,omitempty"`
	WriteFiles                []*WriteFile `yaml:"write_files,omitempty"`
	RunCommands               []string     `yaml:"runcmd,omitempty"`
	FinalMessage              string       `yaml:"final_message,omitempty"`
}

// CreateMetadata
func CreateMetadata(hostname string) *Metadata {
	return &Metadata{
		Hostname:   hostname,
		InstanceID: "00000000-0000-0000-0000-000000000001", // TODO: generate random guid
	}
}

// RenderJson
func (m *Metadata) RenderJson(path string) error {
	out, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, out, 0644)
}

func CreateSSHKeys(privateRsaKey, publicRsaKey string) *SSHKeys {
	return &SSHKeys{
		RSAPrivateKey: privateRsaKey,
		RSAPublicKey:  publicRsaKey,
	}
}

// CreateUserdata
func CreateUserdata(hostname string) *Userdata {
	return &Userdata{
		Hostname:                  hostname,
		ResizeRootFS:              false,
		SSHPasswordAuthentication: true,
		Users:                     []*user{},
		WriteFiles:                []*WriteFile{},
		RunCommands:               []string{},
		FinalMessage:              "Bootstrapped by WSSDAgent using Cloud-Init for " + hostname,
	}
}

// RenderYAML
func (u *Userdata) RenderYAML(path string) error {
	out, err := yaml.Marshal(u)

	if err != nil {
		return err
	}

	out = append([]byte("#cloud-config\n"), out...)
	return ioutil.WriteFile(path, out, 0644)
}

// AddUser
func (ud *Userdata) AddUser(name, gecos, passwd string, groups, ssh_authorized_keys []string, keys *SSHKeys) {
	u := &user{
		Name:              name,
		Gecos:             gecos,
		SSHAuthorizedKeys: ssh_authorized_keys,
		Groups:            groups,
		Sudo:              "ALL=(ALL) NOPASSWD:ALL", // TODO: Make this configurable. This is unrestricted SUDO access
	}
	if keys != nil {
		u.Keys = *keys
	}

	if len(passwd) > 0 {
		u.SetPassword(passwd)
	}
	ud.Users = append(ud.Users, u)
}

//
func (ud *Userdata) GenerateSeedIso(files []string, optionalFiles []string, path string) error {
	writer, err := iso9660.NewWriter()
	if err != nil {
		return err
	}
	defer writer.Cleanup()
	// Add mandatory files
	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		err = writer.AddFile(file, filepath.Base(filePath))
		if err != nil {
			return err
		}
	}
	// Add Optional files
	for _, filePath := range optionalFiles {
		file, err := os.Open(filePath)
		if err != nil { // Check for FileNotFound error only
			continue
		}
		defer file.Close()

		err = writer.AddFile(file, filepath.Base(filePath))
		if err != nil {
			return err
		}
	}

	isofile, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer isofile.Close()

	err = writer.WriteTo(isofile)
	if err != nil {
		return err
	}
	return nil
}

// AddWriteFile
func (ud *Userdata) AddWriteFile(encoding, content, owner, path, permissions string) {
	ud.WriteFiles = append(ud.WriteFiles, &WriteFile{
		Encoding:    encoding,
		Content:     content,
		Owner:       owner,
		Path:        path,
		Permissions: permissions,
	})
}

func (ud *Userdata) AddRunCommand(cmd []string) {
	ud.RunCommands = append(ud.RunCommands, strings.Join(cmd, " "))
}

func (u *user) SetPassword(passwd string) {
	// FIXME
	// u.Password = base64.StdEncoding.EncodeToString([]byte(passwd))
}

type Vendordata struct {
	b64encodedData string
}

// createVendorData
func CreateVendordata(data string) *Vendordata {
	return &Vendordata{b64encodedData: data}
}

// RenderYAML
func (v *Vendordata) RenderYAML(path string) error {
	data, err := base64.StdEncoding.DecodeString(v.b64encodedData)
	if err != nil {
		return err
	}

	//out = append([]byte("#cloud-config\n"), out...)
	return ioutil.WriteFile(path, data, 0644)
}

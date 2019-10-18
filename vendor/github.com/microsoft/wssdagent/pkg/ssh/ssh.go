// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package ssh

import (
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

func ExecuteCommand(userName, remoteHost, privateKeyFile, command string) error {
	keyData, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return err
	}
	key, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return err
	}
	config := &ssh.ClientConfig{
		User:            userName,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}
	client, err := ssh.Dial("tcp", remoteHost+":22", config)
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return session.Run(command)
}

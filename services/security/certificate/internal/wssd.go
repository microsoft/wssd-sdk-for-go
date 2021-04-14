// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license

package internal

import (
	"context"

	"github.com/microsoft/moc/pkg/auth"
	"github.com/microsoft/moc/pkg/errors"
	wssdsecurity "github.com/microsoft/moc/rpc/nodeagent/security"

	"github.com/microsoft/moc/pkg/status"
	wssdclient "github.com/microsoft/wssd-sdk-for-go/pkg/client"
	"github.com/microsoft/wssd-sdk-for-go/services/security"
)

type client struct {
	wssdsecurity.CertificateAgentClient
}

// NewCertificateClientN- creates a client session with the backend agent
func NewCertificateClient(subID string, authorizer auth.Authorizer) (*client, error) {
	c, err := wssdclient.GetCertificateClient(&subID, authorizer)
	if err != nil {
		return nil, err
	}
	return &client{c}, nil
}

// Get
func (c *client) Get(ctx context.Context, group, name string) (*[]security.Certificate, error) {
	request, err := getCertificateRequest(name, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.CertificateAgentClient.Get(ctx, request)
	if err != nil {
		return nil, err
	}
	return getCertificatesFromResponse(response), nil
}

// CreateOrUpdate
func (c *client) CreateOrUpdate(ctx context.Context, group, name string, sg *security.Certificate) (*security.Certificate, error) {
	request, err := getCertificateRequest(name, sg)
	if err != nil {
		return nil, err
	}
	response, err := c.CertificateAgentClient.CreateOrUpdate(ctx, request)
	if err != nil {
		err = errors.Wrapf(err, "[Certificate] Create failed with error %v", err)
		return nil, err
	}

	cert := getCertificatesFromResponse(response)

	if len(*cert) == 0 {
		return nil, errors.New("[Certificate][Create] Unexpected error: Creating a security returned no result")
	}

	return &((*cert)[0]), err
}

// Delete methods invokes create or update on the client
func (c *client) Delete(ctx context.Context, group, name string) error {
	cert, err := c.Get(ctx, group, name)
	if err != nil {
		return err
	}
	if len(*cert) == 0 {
		return errors.Wrapf(errors.NotFound, "Certificate [%s] not found", name)
	}

	request, err := getCertificateRequest(name, &(*cert)[0])
	if err != nil {
		return err
	}
	_, err = c.CertificateAgentClient.Delete(ctx, request)
	return err
}

func getCertificatesFromResponse(response *wssdsecurity.CertificateResponse) *[]security.Certificate {
	certs := []security.Certificate{}
	for _, certificates := range response.GetCertificates() {
		certs = append(certs, *(getCertificate(certificates)))
	}

	return &certs
}

func getCertificateRequest(name string, cert *security.Certificate) (*wssdsecurity.CertificateRequest, error) {
	request := &wssdsecurity.CertificateRequest{
		Certificates: []*wssdsecurity.Certificate{},
	}
	wssdcertificate := &wssdsecurity.Certificate{
		Name: name,
	}

	var err error
	if cert != nil {
		wssdcertificate, err = getWssdCertificate(cert)
		if err != nil {
			return nil, err
		}
	}
	request.Certificates = append(request.Certificates, wssdcertificate)
	return request, nil
}

func getCertificate(cert *wssdsecurity.Certificate) *security.Certificate {

	return &security.Certificate{
		ID:   &cert.Id,
		Name: &cert.Name,
		Cer:  &cert.Certificate,
		Attributes: &security.CertificateAttributes{
			NotBefore:         &cert.NotBefore,
			Expires:           &cert.NotAfter,
			ProvisioningState: status.GetProvisioningState(cert.GetStatus().GetProvisioningStatus()),
			Statuses:          status.GetStatuses(cert.GetStatus()),
		},
	}
}

func getWssdCertificate(cert *security.Certificate) (*wssdsecurity.Certificate, error) {
	if cert.Name == nil {
		return nil, errors.Wrapf(errors.InvalidInput, "Certificate name is missing")
	}
	return &wssdsecurity.Certificate{
		Name: *cert.Name,
	}, nil
}

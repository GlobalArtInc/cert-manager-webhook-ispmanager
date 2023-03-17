package ispmanager

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	panelUrl = "https://k8s-prod-ispmanager.globalart.dev"
	username = ""
	password = ""
)

func TestNewIspClient(t *testing.T) {
	p := NewIspClient(panelUrl, username, password)
	assert.NotNil(t, p.config)
	assert.NotNil(t, p.client)
	assert.Equal(t, p.config.panelUrl, panelUrl)
	assert.Equal(t, p.config.username, username)
	assert.Equal(t, p.config.password, password)
}

func TestIspClient_Present(t *testing.T) {
	var (
		zone = "example.com"
		fqdn = "_acme-challenge.example.com."
		key  = "test_key"
	)
	p := NewIspClient(panelUrl, username, password)
	err := p.Present(zone, fqdn, key)
	assert.NoError(t, err)
	assert.Equal(t, err, nil)
}

func TestIspClient_CleanUp(t *testing.T) {
	var (
		zone = "example.com"
		fqdn = "_acme-challenge.example.com."
		key  = "test_key"
	)
	p := NewIspClient(panelUrl, username, password)
	err := p.CleanUp(zone, fqdn, key)
	assert.NoError(t, err)
	assert.Equal(t, err, nil)
}

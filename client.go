// package ispmanager contains a self-contained ispmanager of a webhook that passes the cert-manager
// DNS conformance tests
package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

type IspClient struct {
	panelUrl string
	username string
	password string
}

func NewIspClient(panelUrl string, username string, password string) *IspClient {
	return &IspClient{
		panelUrl: panelUrl,
		username: username,
		password: password,
	}
}

func (c *IspClient) createTXT(plid string, name string, value string) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	var (
		authInfo     = fmt.Sprintf("%s:%s", c.username, c.password)
		functionName = "domain.record.edit"
		data         = url.Values{
			"authinfo": {authInfo},
			"out":      {"json"},
			"func":     {functionName},
			"plid":     {plid},
			"rtype":    {"txt"},
			"name":     {name},
			"sok":      {"ok"},
			"value":    {value},
		}
	)
	_, err := http.PostForm(c.panelUrl, data)

	if err != nil {
		return fmt.Errorf("failed to make POST request: %v", err)
	}
	log.Info().
		Str("event", "created_txt").
		Str("domain", plid).
		Str("zone", name).Send()

	return nil
}

func (c *IspClient) deleteTXT(plid string, name string, value string) error {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	var (
		authInfo     = fmt.Sprintf("%s:%s", c.username, c.password)
		functionName = "domain.record.delete"
		elid         = fmt.Sprintf("%s TXT  %s", name, value)
		data         = url.Values{
			"authinfo": {authInfo},
			"out":      {"json"},
			"func":     {functionName},
			"plid":     {plid},
			"elid":     {elid},
			"elname":   {elid},
		}
	)
	_, err := http.PostForm(c.panelUrl, data)
	if err != nil {
		return fmt.Errorf("failed to make POST request: %v", err)
	}
	log.Info().
		Str("event", "deleted_txt").
		Str("domain", plid).
		Str("zone", name).Send()

	return nil
}

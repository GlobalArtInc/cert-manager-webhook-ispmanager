// package ispmanager contains a self-contained ispmanager of a webhook that passes the cert-manager
// DNS conformance tests
package main

import (
	"crypto/tls"
	"encoding/json"
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

type APIError struct {
	Doc struct {
		Error struct {
			Type   string `json:"$type"`
			Object string `json:"$object"`
			Lang   string `json:"$lang"`
			Detail struct {
				Text string `json:"$"`
			} `json:"detail"`
			Message struct {
				Text string `json:"$"`
			} `json:"msg"`
		} `json:"error"`
	} `json:"doc"`
}

type APIResponse struct {
	Doc APIError `json:"doc"`
}

func NewIspClient(panelUrl string, username string, password string) *IspClient {
	return &IspClient{
		panelUrl: panelUrl,
		username: username,
		password: password,
	}
}

func checkResponse(res *http.Response) error {
	if res.Body == nil {
		return fmt.Errorf("request failed with status code %v and empty body", res.StatusCode)
	}

	decoder := json.NewDecoder(res.Body)

	var apiError APIError
	err := decoder.Decode(&apiError)
	if err != nil {
		return fmt.Errorf("failed to decode: %s", err)
	}
	if apiError != (APIError{}) {
		//fmt.Printf("dev: %s", apiError.Doc.Error.Message.Text)
		return fmt.Errorf("ISPManager Error: %s", apiError.Doc.Error.Message.Text)
	}

	return nil
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
	res, err := http.PostForm(c.panelUrl, data)
	if err != nil {
		return fmt.Errorf("failed to make POST request: %v", err)
	}
	err = checkResponse(res)
	if err != nil {
		log.Err(err).Send()
		return err
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
	res, err := http.PostForm(c.panelUrl, data)
	if err != nil {
		return fmt.Errorf("failed to make POST request: %v", err)
	}
	err = checkResponse(res)
	if err != nil {
		log.Err(err).Send()
		return err
	}

	log.Info().
		Str("event", "deleted_txt").
		Str("domain", plid).
		Str("zone", name).Send()

	return nil
}

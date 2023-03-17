package ispmanager

import (
	"fmt"
	"github.com/GlobalArtInc/cert-manager-webhook-ispmanager/ispmanager/internal"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type IspClientConfig struct {
	panelUrl string
	username string
	password string
}

type IspClient struct {
	config IspClientConfig
	client *internal.Client
}

func NewIspClient(panelUrl, username, password string) *IspClient {
	return &IspClient{
		config: IspClientConfig{
			panelUrl: panelUrl,
			username: username,
			password: password,
		},
		client: &internal.Client{
			PanelUrl:   panelUrl,
			Username:   username,
			Password:   password,
			HttpClient: &http.Client{},
		},
	}
}

// Present function
func (c IspClient) Present(zone, fqdn, key string) error {
	txtRecord := internal.CreateTxtRecord{
		AuthInfo: fmt.Sprintf("%s:%s", c.config.username, c.config.password),
		Out:      "json",
		Func:     "domain.record.edit",
		Plid:     zone,
		Rtype:    "txt",
		Name:     fqdn,
		Sok:      "ok",
		Value:    key,
	}
	_, err := c.client.CreateTXT(txtRecord)
	if err != nil {
		return fmt.Errorf("isp_manager: %v", err)
	}
	log.WithFields(log.Fields{
		"event":    "created_txt",
		"domain":   zone,
		"zone":     fqdn,
		"panelUrl": c.config.panelUrl,
	}).Info("Created txt record")

	return nil
}

func (c IspClient) CleanUp(zone, fqdn, key string) error {
	elid := fmt.Sprintf("%s TXT  %s", fqdn, key)
	txtRecord := internal.DeleteTxtRecord{
		AuthInfo: fmt.Sprintf("%s:%s", c.config.username, c.config.password),
		Out:      "json",
		Func:     "domain.record.delete",
		Plid:     zone,
		Elid:     elid,
		Elname:   elid,
	}
	_, err := c.client.DeleteTXT(txtRecord)
	if err != nil {
		return fmt.Errorf("isp_manager: %v", err)
	}
	log.WithFields(log.Fields{
		"event":    "deleted_txt",
		"domain":   zone,
		"zone":     fqdn,
		"panelUrl": c.config.panelUrl,
	}).Info("Deleted txt record")

	return nil
}

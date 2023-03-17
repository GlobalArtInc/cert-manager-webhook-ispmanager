package main

import (
	"encoding/json"
	"fmt"
	"github.com/GlobalArtInc/cert-manager-webhook-ispmanager/ispmanager"
	"os"

	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	log "github.com/sirupsen/logrus"
)

var (
	ProviderName = "ispmanager-provider"
	GroupName    = os.Getenv("GROUP_NAME")
)

func main() {
	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}

	cmd.RunWebhookServer(GroupName,
		&ispmanagerDNSProviderSolver{},
	)
}

type ispmanagerDNSProviderSolver struct {
	client *kubernetes.Clientset
}

type ispmanagerDNSProviderConfig struct {
	PanelUrl string `json:"panelUrl"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// Name is used as the name for this DNS solver when referencing it on the ACME
// Issuer resource.
// This should be unique **within the group name**, i.e. you can have two
// solvers configured with the same Name() **so long as they do not co-exist
// within a single webhook deployment**.
// For example, `cloudflare` may be used as the name of a solver.
func (c *ispmanagerDNSProviderSolver) Name() string {
	return ProviderName
}

// Present is responsible for actually presenting the DNS record with the
// DNS provider.
// This method should tolerate being called multiple times with the same value.
// cert-manager itself will later perform a self check to ensure that the
// solver has correctly configured the DNS provider.
func (c *ispmanagerDNSProviderSolver) Present(challengeRequest *v1alpha1.ChallengeRequest) error {
	cfg, err := loadConfig(challengeRequest.Config)
	if err != nil {
		return err
	}
	provider, err := c.provider(&cfg)
	if err != nil {
		return err
	}

	return provider.Present(getDomainFromZone(challengeRequest.ResolvedZone), challengeRequest.ResolvedFQDN, challengeRequest.Key)
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
// If multiple TXT records exist with the same record name (e.g.
// _acme-challenge.example.com) then **only** the record with the same `key`
// value provided on the ChallengeRequest should be cleaned up.
// This is in order to facilitate multiple DNS validations for the same domain
// concurrently.
func (c *ispmanagerDNSProviderSolver) CleanUp(challengeRequest *v1alpha1.ChallengeRequest) error {
	cfg, err := loadConfig(challengeRequest.Config)
	if err != nil {
		return err
	}
	provider, err := c.provider(&cfg)
	if err != nil {
		return err
	}

	return provider.CleanUp(getDomainFromZone(challengeRequest.ResolvedZone), challengeRequest.ResolvedFQDN, challengeRequest.Key)
}

// Initialize will be called when the webhook first starts.
// This method can be used to instantiate the webhook, i.e. initialising
// connections or warming up caches.
// Typically, the kubeClientConfig parameter is used to build a Kubernetes
// client that can be used to fetch resources from the Kubernetes API, e.g.
// Secret resources containing credentials used to authenticate with DNS
// provider accounts.
// The stopCh can be used to handle early termination of the webhook, in cases
// where a SIGTERM or similar signal is sent to the webhook process.
func (c *ispmanagerDNSProviderSolver) Initialize(kubeClientConfig *rest.Config, _ <-chan struct{}) error {
	klog.Infof("call function Initialize")
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}
	c.client = cl
	log.SetFormatter(&log.JSONFormatter{})

	return nil
}

func (c *ispmanagerDNSProviderSolver) provider(cfg *ispmanagerDNSProviderConfig) (*ispmanager.IspClient, error) {
	IspClient := ispmanager.NewIspClient(cfg.PanelUrl, cfg.User, cfg.Password)

	return IspClient, nil
}

// loadConfig is a small helper function that decodes JSON configuration into
// the typed config struct.
func loadConfig(cfgJSON *extapi.JSON) (ispmanagerDNSProviderConfig, error) {
	cfg := ispmanagerDNSProviderConfig{}
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %v", err)
	}

	return cfg, nil
}

func getDomainFromZone(zone string) string {
	return zone[0 : len(zone)-1]
}

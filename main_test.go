package main

import (
	"os"
	"testing"

	"github.com/cert-manager/cert-manager/test/acme/dns"
)

var (
	zone = os.Getenv("TEST_ZONE_NAME")
)

func TestRunsSuite(t *testing.T) {
	fixture := dns.NewFixture(&ispmanagerDNSProviderSolver{},
		dns.SetAllowAmbientCredentials(false),
		dns.SetResolvedZone(zone),
		dns.SetDNSServer("62.109.29.39:53"),
		dns.SetManifestPath("testdata/ispmanager"),
		dns.SetStrict(false),
	)
	// fixture.RunConformance(t)
	fixture.RunBasic(t)
}

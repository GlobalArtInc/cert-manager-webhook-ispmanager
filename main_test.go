package main

import (
	"os"
	"testing"

	"github.com/cert-manager/cert-manager/test/acme/dns"
)

var (
	zone         = os.Getenv("TEST_ZONE_NAME")
	dnsIpAddress = os.Getenv("DNS_SERVER")
)

func TestRunsSuite(t *testing.T) {
	fixture := dns.NewFixture(&ispmanagerDNSProviderSolver{},
		dns.SetAllowAmbientCredentials(false),
		dns.SetResolvedZone(zone),
		dns.SetDNSServer(dnsIpAddress),
		dns.SetManifestPath("testdata/ispmanager"),
		dns.SetStrict(false),
	)
	// fixture.RunConformance(t)
	fixture.RunBasic(t)
}

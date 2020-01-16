package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	whapi "github.com/jetstack/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/jetstack/cert-manager/pkg/acme/webhook/cmd"
	certmgrv1 "github.com/jetstack/cert-manager/pkg/apis/meta/v1"

	"github.com/okteto/civogo/pkg/client"
	"github.com/okteto/civogo/pkg/dns"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"

	"github.com/jetstack/cert-manager/pkg/issuer/acme/dns/util"
)

func main() {
	groupName := os.Getenv("GROUP_NAME")
	if groupName == "" {
		panic("GROUP_NAME must be specified")
	}

	cmd.RunWebhookServer(groupName,
		&civoDNSProviderSolver{},
	)
}

type civoDNSProviderSolver struct {
	client *kubernetes.Clientset
}

type civoDNSProviderConfig struct {
	APIKeySecretRef certmgrv1.SecretKeySelector `json:"apiKeySecretRef"`
}

func (c *civoDNSProviderSolver) Initialize(kubeClientConfig *restclient.Config, stopCh <-chan struct{}) error {
	cl, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}

	c.client = cl
	return nil
}

func (c *civoDNSProviderSolver) Present(ch *whapi.ChallengeRequest) error {
	client, err := c.newClientFromConfig(ch)
	if err != nil {
		return err
	}

	fmt.Printf("Decoded client")

	dnsZone := util.UnFqdn(ch.ResolvedZone)
	d, err := dns.GetDomain(client, dnsZone)
	if err != nil {
		fmt.Printf("Get domain %s error: %s", dnsZone, err)
		return err
	}

	rn := extractRecordName(ch.ResolvedFQDN, dnsZone)

	_, err = d.NewRecord(client, dns.TXT, rn, ch.Key, 10, 600)
	if err != nil {
		return err
	}

	fmt.Printf("Presented txt record %s", ch.ResolvedFQDN)
	return nil
}

func (c *civoDNSProviderSolver) CleanUp(ch *whapi.ChallengeRequest) error {
	client, err := c.newClientFromConfig(ch)
	if err != nil {
		return err
	}

	fmt.Printf("Decoded client")
	dnsZone := util.UnFqdn(ch.ResolvedZone)
	d, err := dns.GetDomain(client, dnsZone)
	if err != nil {
		fmt.Printf("Get domain %s error: %s", dnsZone, err)
		return err
	}

	rn := extractRecordName(ch.ResolvedFQDN, dnsZone)
	r, err := d.GetRecord(client, rn)
	if err != nil {
		fmt.Printf("Get record %s error: %s", rn, err)
		return err
	}

	if r.Value != ch.Key {
		fmt.Printf("Records value does not match: %v", ch.ResolvedFQDN)
		return errors.New("record value does not match")
	}

	if err := r.Delete(client); err != nil {
		fmt.Printf("Delete record %s error: %s", r.Name, err)
	}

	return nil
}

func (c *civoDNSProviderSolver) Name() string {
	return "civo"
}

func (c *civoDNSProviderSolver) newClientFromConfig(ch *whapi.ChallengeRequest) (*client.Client, error) {
	cfg, err := c.loadConfig(ch)
	if err != nil {
		return nil, err
	}

	token, err := c.getSecretData(cfg.APIKeySecretRef, ch.ResourceNamespace)
	if err != nil {
		return nil, err
	}

	return client.New(token), nil
}

func (c *civoDNSProviderSolver) loadConfig(ch *whapi.ChallengeRequest) (*civoDNSProviderConfig, error) {
	cfg := &civoDNSProviderConfig{}
	if ch.Config == nil {
		return cfg, nil
	}

	if err := json.Unmarshal(ch.Config.Raw, cfg); err != nil {
		return cfg, fmt.Errorf("error decoding solver config: %v", err)
	}

	fmt.Printf("Deleted txt record %s", ch.ResolvedFQDN)
	return cfg, nil
}

func (c *civoDNSProviderSolver) getSecretData(selector certmgrv1.SecretKeySelector, ns string) (string, error) {
	secret, err := c.client.CoreV1().Secrets(ns).Get(selector.Name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to load secret %s/%s: %w", ns, selector.Name, err)
	}

	if data, ok := secret.Data[selector.Key]; ok {
		return string(data), nil
	}

	return "", fmt.Errorf("no key %s in secret %s/%s", selector.Key, ns, selector.Name)
}

func extractRecordName(fqdn, domain string) string {
	name := util.UnFqdn(fqdn)
	if idx := strings.Index(name, "."+domain); idx != -1 {
		return name[:idx]
	}
	return name
}

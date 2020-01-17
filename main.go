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
	fmt.Printf("Presenting challenge for fqdn=%s zone=%s\n", ch.ResolvedFQDN, ch.ResolvedZone)
	client, err := c.newClientFromConfig(ch)
	if err != nil {
		return err
	}

	rn, domain := extractNames(ch.ResolvedFQDN)
	d, err := dns.GetDomain(client, domain)
	if err != nil {
		fmt.Printf("Get domain %s error: %s\n", domain, err)
		return err
	}

	_, err = d.NewRecord(client, dns.TXT, rn, ch.Key, 10, 600)
	if err != nil {
		return err
	}

	fmt.Printf("Presented txt record for fqdn=%s zone=%s\n", ch.ResolvedFQDN, ch.ResolvedZone)
	return nil
}

func (c *civoDNSProviderSolver) CleanUp(ch *whapi.ChallengeRequest) error {
	fmt.Printf("Cleaning up for fqdn=%s\n", ch.ResolvedFQDN)
	client, err := c.newClientFromConfig(ch)
	if err != nil {
		return err
	}

	rn, domain := extractNames(ch.ResolvedFQDN)
	d, err := dns.GetDomain(client, domain)
	if err != nil {
		fmt.Printf("Get domain %s error: %s\n", domain, err)
		return err
	}

	r, err := d.GetRecord(client, rn)
	if err != nil {
		fmt.Printf("Get record %s error: %s\n", rn, err)
		return err
	}

	if r.Value != ch.Key {
		fmt.Printf("Records value does not match: %v\n", ch.ResolvedFQDN)
		return errors.New("record value does not match")
	}

	if err := r.Delete(client); err != nil {
		fmt.Printf("Deleted record %s error: %s\n", r.Name, err)
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

func extractNames(fqdn string) (string, string) {
	p := strings.Split(fqdn, ".")
	record := p[0]
	zone := strings.Join(p[1:], ".")
	zone = strings.TrimSuffix(zone, ".")
	return record, zone
}

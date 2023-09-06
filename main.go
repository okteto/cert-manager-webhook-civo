package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	whapi "github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"

	"github.com/civo/civogo"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	ctx := context.Background()

	region := os.Getenv("REGION")
	if region == "" {
		region = "NYC1"
	}

	groupName := os.Getenv("GROUP_NAME")
	if groupName == "" {
		panic("GROUP_NAME must be specified")
	}

	cmd.RunWebhookServer(groupName,
		&civoDNSProviderSolver{ctx: ctx},
	)
}

type civoDNSProviderSolver struct {
	client *kubernetes.Clientset
	ctx    context.Context
	region string
}

type civoDNSProviderConfig struct {
	SecretNameRef string `json:"secretName"`
	SecretKeyRef  string `json:"secretKey"`
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
	log.Infof("Presenting challenge for fqdn=%s zone=%s", ch.ResolvedFQDN, ch.ResolvedZone)
	client, err := c.newClientFromConfig(ch)
	if err != nil {
		log.Errorf("failed to get client from ChallengeRequest: %s", err)
		return err
	}

	rn := strings.TrimSuffix(ch.ResolvedFQDN, "."+ch.ResolvedZone)
	d, err := client.GetDNSDomain(strings.TrimSuffix(ch.ResolvedZone, "."))
	if err != nil {
		log.Errorf("failed to get DNS domain '%s' from civo: %s", ch.ResolvedZone, err)
		return err
	}

	r := &civogo.DNSRecordConfig{
		Name:     rn,
		Value:    ch.Key,
		Type:     civogo.DNSRecordTypeTXT,
		Priority: 10,
		TTL:      600}

	log.Infof("creating DNS record %s/%s", d.ID, r.Name)
	_, err = client.CreateDNSRecord(d.ID, r)
	if err != nil {
		log.Errorf("failed to create DNS Record '%s': %s", r.Name, err)
		return err
	}

	log.Infof("Successfully created txt record for fqdn=%s zone=%s", ch.ResolvedFQDN, ch.ResolvedZone)
	return nil
}

func (c *civoDNSProviderSolver) CleanUp(ch *whapi.ChallengeRequest) error {
	log.Infof("Cleaning up entry for fqdn=%s", ch.ResolvedFQDN)
	client, err := c.newClientFromConfig(ch)
	if err != nil {
		log.Errorf("failed to get client from ChallengeRequest: %s", err)
		return fmt.Errorf("failed to get client from ChallengeRequest: %w", err)
	}
	rn := strings.TrimSuffix(ch.ResolvedFQDN, "."+ch.ResolvedZone)
	r, err := getDNSRecord(client, rn, strings.TrimSuffix(ch.ResolvedZone, "."), ch.Key)
	if err != nil {
		return err
	}

	resp, err := client.DeleteDNSRecord(r)
	if err != nil {
		log.Errorf("failed to delete DNS record '%s': %s", r.Name, err)
		return fmt.Errorf("failed to delete DNS record '%s': %s", r.Name, err)
	}

	if resp.Result == "success" {
		return nil
	}

	log.Errorf("failed to delete DNS record '%s': %s", r.Name, resp)
	return fmt.Errorf("failed to delete DNS record '%s': %s", r.Name, resp)
}

func (c *civoDNSProviderSolver) Name() string {
	return "civo"
}

func (c *civoDNSProviderSolver) newClientFromConfig(ch *whapi.ChallengeRequest) (*civogo.Client, error) {
	cfg, err := c.loadConfig(ch)
	if err != nil {
		return nil, err
	}

	apiKey, err := c.getSecretData(cfg.SecretNameRef, ch.ResourceNamespace, cfg.SecretKeyRef)
	if err != nil {
		return nil, err
	}

	return civogo.NewClient(apiKey, c.region)

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

func (c *civoDNSProviderSolver) getSecretData(secretName string, ns string, secretKey string) (string, error) {
	secret, err := c.client.CoreV1().Secrets(ns).Get(c.ctx, secretName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to load secret %s/%s: %w", ns, secretName, err)
	}

	if data, ok := secret.Data[secretKey]; ok {
		return string(data), nil
	}

	return "", fmt.Errorf("no key %s in secret %s/%s", secretKey, ns, secretName)
}

func getDNSRecord(client *civogo.Client, rn, domain, key string) (*civogo.DNSRecord, error) {
	log.Infof("getting domain %s from civo", domain)
	d, err := client.GetDNSDomain(domain)
	if err != nil {
		log.Errorf("failed to get DNS host '%s' in domain '%s' from civo: %s", rn, domain, err)
		return nil, err
	}

	log.Infof("getting DNS record %s/%s from civo", d.ID, rn)
	rs, err := client.ListDNSRecords(d.ID)
	if err != nil {
		log.Errorf("failed to get DNS Records for '%s': %s", d.ID, err)
		return nil, err
	}

	for _, r := range rs {
		if r.Name == rn {
			if r.Value == key {
				return &r, nil
			}

			log.Infof("Records value for %s does not match %s", r.Name, key)
		}
	}

	return nil, fmt.Errorf("record not found")
}

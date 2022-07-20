package vaultgrafanacloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/api"
)

func GrafanaCloudSecretBackendResource() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        grafanaCloudSecretBackendCreate,
		Delete:        grafanaCloudSecretBackendDelete,
		Read:          grafanaCloudSecretBackendRead,
		Update:        grafanaCloudSecretBackendUpdate,

		Schema: map[string]*schema.Schema{
			"backend": {
				Type:        schema.TypeString,
				Default:     "grafana-cloud",
				ForceNew:    true,
				Optional:    true,
				Description: `The mount path for a backend, for example, the path given in "$ vault secrets enable -path=grafana-cloud grafana-cloud-plugin".`,
				StateFunc: func(v interface{}) string {
					return strings.Trim(v.(string), "/")
				},
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "API key with Admin role to create user keys",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL for the Grafana Cloud API",
			},
			"organisation": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Organisation slug for the Grafana Cloud API",
			},
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The User that is needed to interact with prometheus, if set this is returned alongside every issued credential",
			},
		},
	}
}

func grafanaCloudSecretBackendCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	backend := d.Get("backend").(string)

	log.Printf("[DEBUG] Mounting grafana-cloud-plugin backend at %q", backend)
	err := client.Sys().Mount(backend, &api.MountInput{
		Type: "vault-plugin-secrets-grafanacloud",
	})
	if err != nil {
		return fmt.Errorf("error mounting to %q: %s", backend, err)
	}

	log.Printf("[DEBUG] Mounted vault grafana cloud backend at %q", backend)
	d.SetId(backend)

	data := map[string]interface{}{}
	if v, ok := d.GetOkExists("key"); ok {
		data["key"] = v
	}
	if v, ok := d.GetOkExists("url"); ok {
		data["url"] = v
	}
	if v, ok := d.GetOkExists("organisation"); ok {
		data["organisation"] = v
	}
	if v, ok := d.GetOkExists("user"); ok {
		data["user"] = v
	}

	configPath := fmt.Sprintf("%s/config", backend)
	log.Printf("[DEBUG] Writing %q", configPath)
	if _, err := client.Logical().Write(configPath, data); err != nil {
		return fmt.Errorf("error writing %q: %s", configPath, err)
	}
	log.Printf("[DEBUG] Wrote %q", configPath)
	return grafanaCloudSecretBackendRead(d, meta)
}

func grafanaCloudSecretBackendDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	vaultPath := d.Id()
	log.Printf("[DEBUG] Unmounting vault grafana cloud backend %q", vaultPath)

	err := client.Sys().Unmount(vaultPath)
	if err != nil && strings.Contains(err.Error(), "Code: 404") {
		log.Printf("[WARN] %q not found, removing from state", vaultPath)
		d.SetId("")
		return fmt.Errorf("error unmounting vault grafana cloud backend from %q: %s", vaultPath, err)
	} else if err != nil {
		return fmt.Errorf("error unmounting vault grafana cloud backend from %q: %s", vaultPath, err)
	}
	log.Printf("[DEBUG] Unmounted vault grafana cloud backend %q", vaultPath)
	return nil
}

func grafanaCloudSecretBackendRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)

	if err := d.Set("backend", d.Id()); err != nil {
		return fmt.Errorf("error setting backend: %s", err)
	}

	configPath := fmt.Sprintf("%s/config", d.Id())
	log.Printf("[DEBUG] Reading %q", configPath)

	resp, err := client.Logical().Read(configPath)
	if err != nil {
		return fmt.Errorf("error reading %q: %s", configPath, err)
	}
	log.Printf("[DEBUG] Read %q", configPath)
	if resp == nil {
		log.Printf("[WARN] %q not found, removing from state", configPath)
		d.SetId("")
		return nil
	}

	if val, ok := resp.Data["key"]; ok {
		if err := d.Set("key", val); err != nil {
			return fmt.Errorf("error setting state key 'key': %s", err)
		}
	}
	if val, ok := resp.Data["url"]; ok {
		if err := d.Set("url", val); err != nil {
			return fmt.Errorf("error setting state key 'url': %s", err)
		}
	}
	if val, ok := resp.Data["organisation"]; ok {
		if err := d.Set("organisation", val); err != nil {
			return fmt.Errorf("error setting state key 'organisation': %s", err)
		}
	}
	if val, ok := resp.Data["user"]; ok {
		if err := d.Set("user", val); err != nil {
			return fmt.Errorf("error setting state key 'user': %s", err)
		}
	}
	return nil
}

func grafanaCloudSecretBackendUpdate(d *schema.ResourceData, meta interface{}) error {
	backend := d.Id()

	client := meta.(*api.Client)
	data := map[string]interface{}{}

	vaultPath := fmt.Sprintf("%s/config", backend)
	log.Printf("[DEBUG] Updating %q", vaultPath)

	if raw, ok := d.GetOk("key"); ok {
		data["key"] = raw
	}
	if raw, ok := d.GetOk("url"); ok {
		data["url"] = raw
	}
	if raw, ok := d.GetOk("organisation"); ok {
		data["organisation"] = raw
	}
	if raw, ok := d.GetOk("user"); ok {
		data["user"] = raw
	}
	if _, err := client.Logical().Write(vaultPath, data); err != nil {
		return fmt.Errorf("error updating template secrets backend role %q: %s", vaultPath, err)
	}
	log.Printf("[DEBUG] Updated %q", vaultPath)
	return nil
}

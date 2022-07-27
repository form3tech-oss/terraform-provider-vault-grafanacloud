package vaultgrafanacloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/vault/api"
)

var (
	gcSecretFromPathRegex         = regexp.MustCompile("^(.+)/roles/.+$")
	gcSecretRoleNameFromPathRegex = regexp.MustCompile("^.+/roles/(.+$)")
)

func GrafanaCloudSecretRoleResource() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,
		Create:        grafanaCloudSecretRoleCreate,
		Delete:        grafanaCloudSecretRoleDelete,
		Read:          grafanaCloudSecretRoleRead,
		Update:        grafanaCloudSecretRoleUpdate,

		Schema: map[string]*schema.Schema{
			"backend": {
				Type:        schema.TypeString,
				Default:     "grafana-cloud",
				Optional:    true,
				ForceNew:    true,
				Description: "The mount path of the Grafana Cloud backend.",
				StateFunc: func(v interface{}) string {
					return strings.Trim(v.(string), "/")
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name for the role",
			},
			"gc_role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Grafana Cloud role, i.e. the key authorization level",
			},
			"ttl_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Default lease for generated credentials in seconds",
			},
			"max_ttl_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     300,
				Description: "Maximum time for role in seconds",
			},
		},
	}
}

func grafanaCloudSecretRoleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	backend := d.Get("backend").(string)
	name := d.Get("name").(string)
	rolePath := fmt.Sprintf("%s/roles/%s", backend, name)

	log.Printf("[DEBUG] Creating %q", rolePath)

	data := map[string]interface{}{}
	if v, ok := d.GetOkExists("gc_role"); ok {
		data["gc_role"] = v
	}
	if v, ok := d.GetOkExists("ttl_seconds"); ok {
		data["ttl_seconds"] = v
	}
	if v, ok := d.GetOkExists("max_ttl_seconds"); ok {
		data["max_ttl_seconds"] = v
	}

	log.Printf("[DEBUG] Writing %q", rolePath)
	if _, err := client.Logical().Write(rolePath, data); err != nil {
		return fmt.Errorf("error writing %q: %s", rolePath, err)
	}
	d.SetId(rolePath)
	log.Printf("[DEBUG] Wrote %q", rolePath)
	return grafanaCloudSecretRoleRead(d, meta)
}

func grafanaCloudSecretRoleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	rolePath := d.Id()
	log.Printf("[DEBUG] Deleting %q", rolePath)

	if _, err := client.Logical().Delete(rolePath); err != nil && !strings.Contains(err.Error(), "Code: 404") {
		return fmt.Errorf("error deleting %q: %s", rolePath, err)
	} else if err != nil {
		log.Printf("[DEBUG] %q not found, removing from state", rolePath)
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] Deleted template auth backend role %q", rolePath)
	return nil
}

func grafanaCloudSecretRoleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	rolePath := d.Id()
	log.Printf("[DEBUG] Reading %q", rolePath)

	roleName, err := gcSecretRoleNameFromPath(rolePath)
	if err != nil {
		return fmt.Errorf("invalid role ID %q: %s", rolePath, err)
	}
	if err := d.Set("name", roleName); err != nil {
		return fmt.Errorf("error setting name: %s", err)
	}

	backend, err := gcSecretFromPath(rolePath)
	if err != nil {
		return fmt.Errorf("invalid role ID %q: %s", rolePath, err)
	}
	if err := d.Set("backend", backend); err != nil {
		return fmt.Errorf("error setting backend: %s", err)
	}

	resp, err := client.Logical().Read(rolePath)
	if err != nil {
		return fmt.Errorf("error reading %q: %s", rolePath, err)
	}
	log.Printf("[DEBUG] Read %q", rolePath)

	if resp == nil {
		log.Printf("[WARN] %q not found, removing from state", rolePath)
		d.SetId("")
		return nil
	}

	if val, ok := resp.Data["gc_role"]; ok {
		if err := d.Set("gc_role", val); err != nil {
			return fmt.Errorf("error setting state key 'gc_role': %s", err)
		}
	}

	if val, ok := resp.Data["ttl_seconds"]; ok {
		if err := d.Set("ttl_seconds", val); err != nil {
			return fmt.Errorf("error setting state key 'ttl_seconds': %s", err)
		}
	}

	if val, ok := resp.Data["max_ttl_seconds"]; ok {
		if err := d.Set("max_ttl_seconds", val); err != nil {
			return fmt.Errorf("error setting state key 'max_ttl_seconds': %s", err)
		}
	}
	return nil
}

func grafanaCloudSecretRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.Client)
	rolePath := d.Id()
	log.Printf("[DEBUG] Updating %q", rolePath)

	data := map[string]interface{}{}
	if raw, ok := d.GetOk("gc_role"); ok {
		data["gc_role"] = raw
	}
	if raw, ok := d.GetOk("ttl_seconds"); ok {
		data["ttl_seconds"] = raw
	}
	if raw, ok := d.GetOk("max_ttl_seconds"); ok {
		data["max_ttl_seconds"] = raw
	}
	if _, err := client.Logical().Write(rolePath, data); err != nil {
		return fmt.Errorf("error updating template auth backend role %q: %s", rolePath, err)
	}
	log.Printf("[DEBUG] Updated %q", rolePath)
	return grafanaCloudSecretRoleRead(d, meta)
}

func gcSecretRoleNameFromPath(path string) (string, error) {
	if !gcSecretRoleNameFromPathRegex.MatchString(path) {
		return "", fmt.Errorf("no name found")
	}
	res := gcSecretRoleNameFromPathRegex.FindStringSubmatch(path)
	if len(res) != 2 {
		return "", fmt.Errorf("unexpected number of matches (%d) for name", len(res))
	}
	return res[1], nil
}

func gcSecretFromPath(path string) (string, error) {
	if !gcSecretFromPathRegex.MatchString(path) {
		return "", fmt.Errorf("no backend found")
	}
	res := gcSecretFromPathRegex.FindStringSubmatch(path)
	if len(res) != 2 {
		return "", fmt.Errorf("unexpected number of matches (%d) for backend", len(res))
	}
	return res[1], nil
}

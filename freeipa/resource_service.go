package freeipa

import (
	"context"
	"strings"

	api "terraform-provider-freeipa/freeipa/api"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func schemaService() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"krbcanonicalname": {
			Description: "Service canonical name (in the form of service/host_fqdn)",
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			ValidateDiagFunc: func(val interface{}, p cty.Path) diag.Diagnostics {
				var diags diag.Diagnostics

				if !strings.Contains(val.(string), "/") {
					diag := diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Invalid service name",
						Detail:   "Service name must be in the form of service/host_fqdn",
					}
					diags = append(diags, diag)
				}

				return diags
			},
		},
	}
}

func resourceService() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage FreeIPA services",
		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		// UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,
		Schema:        schemaService(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func flattenService(service *api.Service) JSON {
	flat := JSON{
		"krbcanonicalname": service.KrbCanonicalName[0],
	}

	return flat
}

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*api.APIClient)
	service, err := client.ServiceAdd(d.Get("krbcanonicalname").(string), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(service.KrbCanonicalName[0])

	return diags
}

func resourceServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*api.APIClient)
	service, err := client.ServiceShow(d.Id(), nil)
	if err != nil {
		if err.(*api.APIError).Code == 4001 { // Service not found
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	for k, v := range flattenService(service) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

//lint:ignore U1000 Function is unused at the moment, but it's here for potential future use
func resourceServiceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*api.APIClient)

	options := JSON{}

	if d.HasChange("krbcanonicalname") {
		_, err := client.ServiceMod(d.Id(), options)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGroupRead(ctx, d, m)
}

func resourceServiceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*api.APIClient)
	_, err := client.ServiceDel(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

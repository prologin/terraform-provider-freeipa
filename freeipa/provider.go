package freeipa

import (
	"context"
	api "terraform-provider-freeipa/freeipa/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type JSON = map[string]interface{}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FREEIPA_SERVER", nil),
			},
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FREEIPA_USER", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FREEIPA_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"freeipa_user":             resourceUser(),
			"freeipa_group":            resourceGroup(),
			"freeipa_service":          resourceService(),
			"freeipa_idp":              resourceIdentityProvider(),
			"freeipa_group_membership": resourceGroupMembership(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	server := d.Get("server").(string)
	user := d.Get("user").(string)
	password := d.Get("password").(string)

	var diags diag.Diagnostics
	if server == "" || user == "" || password == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to find server, user or password",
			Detail:   "The server, user and password must be provided",
		})
		return nil, diags
	}

	client, err := api.NewClient(server, user, password)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}

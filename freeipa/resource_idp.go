package freeipa

import (
	"context"
	"regexp"

	api "terraform-provider-freeipa/freeipa/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func schemaIdentityProvider() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cn": {
			Description:      "Identity Provider server name",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"clientid": {
			Description:      "Client ID",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"clientsecret": {
			Description:      "Client secret",
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"authendpoint": {
			Description:      "OAuth 2.0 authorization endpoint",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
		},
		"devauthendpoint": {
			Description:      "Device authorization endpoint",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
		},
		"tokenendpoint": {
			Description:      "Token endpoint",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
		},
		"userinfoendpoint": {
			Description:      "User information endpoint",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
		},
		"keysendpoint": {
			Description:      "JWKS endpoint",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
		},
		"issuerurl": {
			Description:      "OIDC URL",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
		},
		"scope": {
			Description:      "Scope (space separated)",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`^[\w\d.-]+([\s]?[\w\d.-]+)*$`), "Scope must be a space separated list of characters, numbers, and underscores")),
		},
		"sub": {
			Description:      "External IdP user identifier attribute",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
	}
}

func resourceIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage FreeIPA IDP (Identity Provider)",
		CreateContext: resourceIdentityProviderCreate,
		ReadContext:   resourceIdentityProviderRead,
		UpdateContext: resourceIdentityProviderUpdate,
		DeleteContext: resourceIdentityProviderDelete,
		Schema:        schemaIdentityProvider(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func flattenIdentityProvider(idp *api.IdentityProvider) JSON {
	flat := JSON{
		"cn":               idp.CN[0],
		"clientid":         idp.ClientID[0],
		"clientsecret":     idp.ClientSecret[0].Secret.Decode(),
		"authendpoint":     idp.AuthEndpoint[0],
		"devauthendpoint":  idp.DevAuthEndpoint[0],
		"tokenendpoint":    idp.TokenEndpoint[0],
		"userinfoendpoint": idp.UserInfoEndpoint[0],
		"keysendpoint":     idp.KeysEndpoint[0],
	}

	if len(idp.IssuerURL) > 0 {
		flat["issuerurl"] = idp.IssuerURL[0]
	} else {
		flat["issuerurl"] = ""
	}

	if len(idp.Scope) > 0 {
		flat["scope"] = idp.Scope[0]
	} else {
		flat["scope"] = ""
	}

	if len(idp.Sub) > 0 {
		flat["sub"] = idp.Sub[0]
	} else {
		flat["sub"] = ""
	}

	return nil
}

func resourceIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*api.APIClient)

	options := JSON{}
	if val, ok := d.GetOk("issuerurl"); ok {
		options["ipaidpissuerurl"] = val.(string)
	}
	if val, ok := d.GetOk("scope"); ok {
		options["ipaidpscope"] = val.(string)
	}
	if val, ok := d.GetOk("sub"); ok {
		options["ipaidpsub"] = val.(string)
	}

	idp, err := client.IdentityProviderAddGeneric(
		d.Get("cn").(string),
		d.Get("clientid").(string),
		d.Get("clientsecret").(string),
		d.Get("authendpoint").(string),
		d.Get("devauthendpoint").(string),
		d.Get("tokenendpoint").(string),
		d.Get("userinfoendpoint").(string),
		d.Get("keysendpoint").(string),
		options,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(idp.CN[0])

	return diags
}

func resourceIdentityProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*api.APIClient)
	idp, err := client.IdentityProviderShow(d.Id(), JSON{
		"all": true, // Otherwise we don't get the client secret
	})
	if err != nil {
		if err.(*api.APIError).Code == 4001 { // IDP not found
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	for key, value := range flattenIdentityProvider(idp) {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*api.APIClient)

	options := JSON{}
	if d.HasChange("cn") {
		options["rename"] = d.Get("cn").(string)
	}
	if d.HasChange("clientid") {
		options["ipaidpclientid"] = d.Get("clientid").(string)
	}
	if d.HasChange("clientsecret") {
		options["ipaidpclientsecret"] = d.Get("clientsecret").(string)
	}
	if d.HasChange("authendpoint") {
		options["ipaidpauthendpoint"] = d.Get("authendpoint").(string)
	}
	if d.HasChange("devauthendpoint") {
		options["ipaidpdevauthendpoint"] = d.Get("devauthendpoint").(string)
	}
	if d.HasChange("tokenendpoint") {
		options["ipaidptokenendpoint"] = d.Get("tokenendpoint").(string)
	}
	if d.HasChange("userinfoendpoint") {
		options["ipaidpuserinfoendpoint"] = d.Get("userinfoendpoint").(string)
	}
	if d.HasChange("keysendpoint") {
		options["ipaidpkeysendpoint"] = d.Get("keysendpoint").(string)
	}
	if d.HasChange("issuerurl") {
		options["ipaidpissuerurl"] = d.Get("issuerurl").(string)
	}
	if d.HasChange("scope") {
		options["ipaidpscope"] = d.Get("scope").(string)
	}
	if d.HasChange("sub") {
		options["ipaidpsub"] = d.Get("sub").(string)
	}

	_, err := client.IdentityProviderMod(d.Id(), options)
	if err != nil {
		return diag.FromErr(err)
	}

	// If the CN changed, update the ID
	if d.HasChange("cn") {
		d.SetId(d.Get("cn").(string))
	}

	return resourceIdentityProviderRead(ctx, d, m)
}

func resourceIdentityProviderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*api.APIClient)
	_, err := client.IdentityProviderDel(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

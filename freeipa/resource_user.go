package freeipa

import (
	"context"
	"regexp"

	api "terraform-provider-freeipa/freeipa/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func schemaUser() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"uid": {
			Description: "User UID (login)",
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.All(
				validation.StringIsNotWhiteSpace,
				StringContainsNoUpperLetter,
				StringIsNotOnlyDigits,
			)),
		},
		"givenname": {
			Description:      "First name",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"sn": {
			Description:      "Last name",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"password": { // TODO: fix password update after a resource import
			Description:      "User password",
			Type:             schema.TypeString,
			Required:         true,
			Sensitive:        true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"krbpasswordexpiration": {
			Description:      "Password expiration date (in RFC3339 format)\nIf not specified, the password will be immediately expired. This follows the default behavior of the API.",
			Type:             schema.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsRFC3339Time),
		},
		"mail": {
			Description: "Email addresses\nIf not specified, no email will be set. Note that this DOES NOT follows the API default behavior (that would have been to create UID@REALM email by default).",
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringMatch(regexp.MustCompile(`(?i)^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,4}$`), "must be a valid email address")),
			},
		},
		"homedirectory": {
			Description: "Home directory\nIf not specified, the default home directory will be used.",
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
	}
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage FreeIPA users",
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Schema:        schemaUser(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func flattenUser(user *api.User) JSON {
	flat := JSON{
		"uid": user.UID[0],
	}

	if len(user.KrbPasswordExpiration) > 0 {
		flat["krbpasswordexpiration"] = user.KrbPasswordExpiration[0].DateTime.String()
	} else {
		flat["krbpasswordexpiration"] = ""
	}

	if len(user.GivenName) > 0 {
		flat["givenname"] = user.GivenName[0]
	} else {
		flat["givenname"] = ""
	}

	if len(user.SN) > 0 {
		flat["sn"] = user.SN[0]
	} else {
		flat["sn"] = ""
	}

	if len(user.Mail) > 0 {
		flat["mail"] = user.Mail
	} else {
		flat["mail"] = make([]string, 0)
	}

	if len(user.HomeDirectory) > 0 {
		flat["homedirectory"] = user.HomeDirectory[0]
	} else {
		flat["homedirectory"] = ""
	}

	return flat
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*api.APIClient)

	options := JSON{
		"userpassword":          d.Get("password").(string),
		"krbpasswordexpiration": d.Get("krbpasswordexpiration").(string),
		"mail":                  d.Get("mail").([]interface{}),
	}
	homedir := d.Get("homedirectory").(string)
	if homedir != "" {
		options["homedirectory"] = homedir
	}

	user, err := client.UserAdd(d.Get("uid").(string),
		d.Get("givenname").(string),
		d.Get("sn").(string),
		options,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.UID[0])

	return resourceUserRead(ctx, d, m)
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*api.APIClient)
	user, err := client.UserShow(d.Id(), JSON{
		"all": true, // Retrieves all attributes, this is MANDATORY
	})
	if err != nil {
		if err.(*api.APIError).Code == 4001 { // User not found
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	for key, value := range flattenUser(user) {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*api.APIClient)

	options := JSON{}
	if d.HasChange("givenname") {
		options["givenname"] = d.Get("givenname").(string)
	}
	if d.HasChange("sn") {
		options["sn"] = d.Get("sn").(string)
	}
	if d.HasChange("password") {
		options["userpassword"] = d.Get("password").(string)
	}
	if d.HasChange("krbpasswordexpiration") {
		options["krbpasswordexpiration"] = d.Get("krbpasswordexpiration").(string)
	}
	if d.HasChange("mail") {
		options["mail"] = d.Get("mail").([]interface{})
	}
	if d.HasChange("homedirectory") {
		options["homedirectory"] = d.Get("homedirectory").(string)
	}

	if d.HasChangeExcept("uid") {
		_, err := client.UserMod(d.Id(), options)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*api.APIClient)
	_, err := client.UserDel(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

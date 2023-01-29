package freeipa

import (
	"context"

	api "terraform-provider-freeipa/freeipa/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func schemaGroup() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cn": {
			Description: "Group name",
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.All(
				validation.StringIsNotWhiteSpace,
				StringContainsNoUpperLetter,
				StringIsNotOnlyDigits,
			)),
		},
		"description": {
			Description: "First name",
			Type:        schema.TypeString,
			Optional:    true,
		},
	}
}

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage FreeIPA groups",
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Schema:        schemaGroup(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func flattenGroup(group *api.Group) JSON {
	flat := JSON{
		"cn": group.CN[0],
	}

	if len(group.Description) > 0 {
		flat["description"] = group.Description[0]
	}

	return flat
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*api.APIClient)
	group, err := client.GroupAdd(d.Get("cn").(string), JSON{
		"description": d.Get("description").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(group.CN[0])

	return diags
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*api.APIClient)
	group, err := client.GroupShow(d.Id(), nil)
	if err != nil {
		if err.(*api.APIError).Code == 4001 { // Group not found
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	for key, value := range flattenGroup(group) {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*api.APIClient)

	options := JSON{}
	if d.HasChange("description") {
		options["description"] = d.Get("description").(string)
	}

	if d.HasChangeExcept("cn") {
		_, err := client.GroupMod(d.Id(), options)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*api.APIClient)
	_, err := client.GroupDel(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

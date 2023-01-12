package freeipa

import (
	"context"
	api "terraform-provider-freeipa/freeipa/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Schema for group membership

func schemaGroupMembership() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group": {
			Description:      "Group name (CN)",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"member": {
			Description:      "Member identifier",
			Type:             schema.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"type": {
			Description:      `Member type (must be one of "user", "group" or "service")`,
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "user",
			ForceNew:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"user", "group", "service"}, false)),
		},
		"manager": {
			Description: "The member is a manager of the group (connot be used with type service)",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			ForceNew:    true,
		},
	}
}

func resourceGroupMembership() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage FreeIPA group membership",
		CreateContext: resourceGroupMembershipCreate,
		ReadContext:   resourceGroupMembershipRead,
		DeleteContext: resourceGroupMembershipDelete,
		Schema:        schemaGroupMembership(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*api.APIClient)
	var GroupAddMembership = client.GroupAddMember

	if d.Get("manager").(bool) {
		if d.Get("type").(string) == "service" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid group membership",
				Detail:   "A service cannot be a manager of a group",
			})
			return diags
		}

		GroupAddMembership = client.GroupAddMemberManager
	}

	_, err := GroupAddMembership(d.Get("group").(string), JSON{
		d.Get("type").(string): d.Get("member").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("group").(string) + ":" + d.Get("member").(string))

	return diags
}

// Cannot implement this since the API only allow to retrieve user's groups
func resourceGroupMembershipRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	if d.Get("type").(string) != "user" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Cannot check group membership",
			Detail:   "The group membership cannot be checked for non-user members, you may need to manually check it",
		})
		return diags
	}

	client := m.(*api.APIClient)
	groups, err := client.GetGroups(d.Get("member").(string), d.Get("type").(string))
	if err != nil {
		if err.(*api.APIError).Code == 4001 { // User, group or service not found
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	for _, group_cn := range groups.CNs {
		if group_cn == d.Get("group").(string) {
			return nil
		}
	}

	// If we reach this point, the group membership does not exist
	d.SetId("")

	return diags
}

func resourceGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*api.APIClient)
	var GroupRemoveMembership = client.GroupRemoveMember

	if d.Get("manager").(bool) {
		GroupRemoveMembership = client.GroupRemoveMemberManager
	}

	_, err := GroupRemoveMembership(d.Get("group").(string), JSON{
		d.Get("type").(string): d.Get("member").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

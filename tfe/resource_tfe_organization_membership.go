package tfe

import (
	"fmt"
	"log"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTFEOrganizationMembership() *schema.Resource {
	return &schema.Resource{
		Create: resourceTFEOrganizationMembershipCreate,
		Read:   resourceTFEOrganizationMembershipRead,
		Delete: resourceTFEOrganizationMembershipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTFEOrganizationMembershipCreate(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	// Get the email and organization.
	email := d.Get("email").(string)
	organization := d.Get("organization").(string)

	// Create a new options struct.
	options := tfe.OrganizationMembershipCreateOptions{
		Email: tfe.String(email),
	}

	log.Printf("[DEBUG] Create membership %s for organization: %s", email, organization)
	membership, err := tfeClient.OrganizationMemberships.Create(ctx, organization, options)
	if err != nil {
		return fmt.Errorf(
			"Error creating membership %s for organization %s: %w", email, organization, err)
	}

	d.SetId(membership.ID)

	return resourceTFEOrganizationMembershipRead(d, meta)
}

func resourceTFEOrganizationMembershipRead(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	options := tfe.OrganizationMembershipReadOptions{
		Include: []tfe.OrgMembershipIncludeOpt{tfe.OrgMembershipUser},
	}

	log.Printf("[DEBUG] Read configuration of membership: %s", d.Id())
	membership, err := tfeClient.OrganizationMemberships.ReadWithOptions(ctx, d.Id(), options)

	if err != nil {
		if err == tfe.ErrResourceNotFound {
			log.Printf("[DEBUG] Membership %s does no longer exist", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading configuration of membership %s: %w", d.Id(), err)
	}

	// Update the config.
	log.Printf("[INFO] User = %#v", membership.User)
	d.Set("email", membership.Email)
	d.Set("organization", membership.Organization.Name)
	d.Set("user_id", membership.User.ID)

	return nil
}

func resourceTFEOrganizationMembershipDelete(d *schema.ResourceData, meta interface{}) error {
	tfeClient := meta.(*tfe.Client)

	log.Printf("[DEBUG] Delete membership: %s", d.Id())
	err := tfeClient.OrganizationMemberships.Delete(ctx, d.Id())
	if err != nil {
		if err == tfe.ErrResourceNotFound {
			return nil
		}
		return fmt.Errorf("Error deleting membership %s: %w", d.Id(), err)
	}

	return nil
}

package jumpcloud

import (
	"context"
	"fmt"

	jcapiv1 "github.com/TheJumpCloud/jcapi-go/v1"
	jcapiv2 "github.com/TheJumpCloud/jcapi-go/v2"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"xorgid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"firstname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lastname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_mfa": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ldap_binding": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"id_sync": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"global_admin": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"passwordless_sudo": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"unix_uid": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"unix_guid": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			// Currently, only the options necessary for our use case are implemented
			// JumpCloud offers a lot more
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// We receive a v2config from the TF base code but need a v1config to continue. So, we take the only
// preloaded element (the x-api-key) and populate the v1config with it.
func convertV2toV1Config(v2config *jcapiv2.Configuration) *jcapiv1.Configuration {
	configv1 := jcapiv1.NewConfiguration()
	configv1.AddDefaultHeader("x-api-key", v2config.DefaultHeader["x-api-key"])
	return configv1
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	payload := jcapiv1.Systemuserputpost{
		Username:                    d.Get("username").(string),
		Email:                       d.Get("email").(string),
		Firstname:                   d.Get("firstname").(string),
		Lastname:                    d.Get("lastname").(string),
		EnableUserPortalMultifactor: d.Get("enable_mfa").(bool),
		LdapBindingUser:             d.Get("ldap_binding").(bool),
		EnableManagedUid:            d.Get("id_sync").(bool),
		Sudo:                        d.Get("global_admin").(bool),
		PasswordlessSudo:            d.Get("passwordless_sudo").(bool),
		UnixUid:                     d.Get("unix-uid").(int32),
		UnixGuid:                    d.Get("unix-guid").(int32),
	}

	req := map[string]interface{}{
		"body":   payload,
		"xOrgId": d.Get("xorgid").(string),
	}
	returnstruc, _, err := client.SystemusersApi.SystemusersPost(context.TODO(),
		"", "", req)
	if err != nil {
		return err
	}
	d.SetId(returnstruc.Id)
	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	res, _, err := client.SystemusersApi.SystemusersGet(context.TODO(),
		d.Id(), "", "", nil)

	// If the object does not exist in our infrastructure, we unset the ID
	// Unfortunately, the http request returns 200 even if the resource does not exist
	if err != nil {
		if err.Error() == "EOF" {
			d.SetId("")
			return nil
		}
		return err
	}

	d.SetId(res.Id)

	if err := d.Set("username", res.Username); err != nil {
		return err
	}
	if err := d.Set("email", res.Email); err != nil {
		return err
	}
	if err := d.Set("firstname", res.Firstname); err != nil {
		return err
	}
	if err := d.Set("lastname", res.Lastname); err != nil {
		return err
	}
	if err := d.Set("enable_mfa", res.EnableUserPortalMultifactor); err != nil {
		return err
	}
	if err := d.Set("ldap_binding", res.LdapBindingUser); err != nil {
		return err
	}
	if err := d.Set("id_sync", res.EnableManagedUid); err != nil {
		return err
	}
	if err := d.Set("global_admin", res.Sudo); err != nil {
		return err
	}
	if err := d.Set("passwordless_sudo", res.PasswordlessSudo); err != nil {
		return err
	}
	if err := d.Set("unix_uid", res.UnixUid); err != nil {
		return err
	}
	if err := d.Set("unix_gid", res.UnixGuid); err != nil {
		return err
	}

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	// The code from the create function is almost identical, but the structure is different :
	// jcapiv1.Systemuserput != jcapiv1.Systemuserputpost
	payload := jcapiv1.Systemuserput{
		Username:                    d.Get("username").(string),
		Email:                       d.Get("email").(string),
		Firstname:                   d.Get("firstname").(string),
		Lastname:                    d.Get("lastname").(string),
		EnableUserPortalMultifactor: d.Get("enable_mfa").(bool),
		LdapBindingUser:             d.Get("ldap_binding").(bool),
		EnableManagedUid:            d.Get("id_sync").(bool),
		Sudo:                        d.Get("global_admin").(bool),
		// PasswordlessSudo:            d.Get("passwordless_sudo").(bool),
		// This is not partof this object for some reason, changes to this will need to happen elsewhere
		UnixUid:  d.Get("unix-uid").(int32),
		UnixGuid: d.Get("unix-guid").(int32),
	}

	req := map[string]interface{}{
		"body":   payload,
		"xOrgId": d.Get("xorgid").(string),
	}
	_, _, err := client.SystemusersApi.SystemusersPut(context.TODO(),
		d.Id(), "", "", req)
	if err != nil {
		return err
	}
	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	configv1 := convertV2toV1Config(m.(*jcapiv2.Configuration))
	client := jcapiv1.NewAPIClient(configv1)

	res, _, err := client.SystemusersApi.SystemusersDelete(context.TODO(),
		d.Id(), "", headerAccept, nil)
	if err != nil {
		// TODO: sort out error essentials
		return fmt.Errorf("error deleting user group:%s; response = %+v", err, res)
	}
	d.SetId("")
	return nil
}

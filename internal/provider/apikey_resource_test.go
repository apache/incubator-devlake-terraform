// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "devlake_apikey" "tfresourcename" {
  allowed_path = ".*"
  expired_at   = "2025-02-28T09:12:00.153Z"
  name         = "should_not_exist"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devlake_apikey.tfresourcename", "allowed_path", ".*"),
					resource.TestCheckResourceAttr("devlake_apikey.tfresourcename", "expired_at", "2025-02-28T09:12:00.153Z"),
					resource.TestCheckResourceAttr("devlake_apikey.tfresourcename", "name", "should_not_exist"),
					// Verify apikey has Computed attributes filled. Can not check api_key here as it is dynamic
					resource.TestCheckResourceAttr("devlake_apikey.tfresourcename", "type", "devlake"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devlake_apikey.tfresourcename", "id"),
					resource.TestCheckResourceAttrSet("devlake_apikey.tfresourcename", "last_updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devlake_apikey.tfresourcename",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does exist in the devlake API, but
				// we rather want the terraform state here
				ImportStateVerifyIgnore: []string{"last_updated", "api_key"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "devlake_apikey" "tfresourcename" {
  allowed_path = ".*"
  expired_at   = "2026-02-28T09:12:00.153Z"
  name         = "should_not_exist"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify apikey was re-created, there is no way to update an apikey
					resource.TestCheckResourceAttr("devlake_apikey.tfresourcename", "allowed_path", ".*"),
					resource.TestCheckResourceAttr("devlake_apikey.tfresourcename", "expired_at", "2026-02-28T09:12:00.153Z"),
					resource.TestCheckResourceAttr("devlake_apikey.tfresourcename", "name", "should_not_exist"),
					// Verify apikey has Computed attributes filled. Can not check api_key here as it is dynamic
					resource.TestCheckResourceAttr("devlake_apikey.tfresourcename", "type", "devlake"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

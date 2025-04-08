// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBitbucketServerConnectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "devlake_bitbucketserver_connection" "tfresourcename" {
  endpoint  = "https://bitbucket-server.org"
  name      = "should_not_exist"
  password  = "whatever"
  username  = "serviceAccount"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.tfresourcename", "endpoint", "https://bitbucket-server.org"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.tfresourcename", "name", "should_not_exist"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.tfresourcename", "password", "whatever"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.tfresourcename", "proxy", ""),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.tfresourcename", "rate_limit_per_hour", "0"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.tfresourcename", "username", "serviceAccount"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.tfresourcename", "id"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.tfresourcename", "last_updated"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.tfresourcename", "created_at"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.tfresourcename", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "devlake_bitbucketserver_connection.tfresourcename",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does exist in the devlake API, but
				// we rather want the terraform state here
				ImportStateVerifyIgnore: []string{"password", "last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "devlake_bitbucketserver_connection" "tfresourcename" {
  endpoint  = "https://bitbucket-server2.org"
  name      = "should_not_exist"
  password  = "whatever"
  username  = "serviceAccount"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.tfresourcename", "endpoint", "https://bitbucket-server2.org"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.tfresourcename", "name", "should_not_exist"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.tfresourcename", "password", "whatever"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.tfresourcename", "proxy", ""),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.tfresourcename", "rate_limit_per_hour", "0"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.tfresourcename", "username", "serviceAccount"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.tfresourcename", "id"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.tfresourcename", "last_updated"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.tfresourcename", "created_at"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.tfresourcename", "updated_at"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	bitbucketServerConnectionConfig = providerConfig + `
resource "devlake_bitbucketserver_connection" "bbserver" {
  endpoint  = "https://bitbucket-server.org"
  name      = "should_not_exist"
  password  = "whatever"
  username  = "serviceAccount"
}
`
)

func TestAccBitbucketServerConnectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: bitbucketServerConnectionConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.bbserver", "endpoint", "https://bitbucket-server.org"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.bbserver", "name", "should_not_exist"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.bbserver", "password", "whatever"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.bbserver", "proxy", ""),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.bbserver", "rate_limit_per_hour", "0"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.bbserver", "username", "serviceAccount"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.bbserver", "id"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.bbserver", "last_updated"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.bbserver", "created_at"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.bbserver", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName: "devlake_bitbucketserver_connection.bbserver",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					if rs, ok := s.RootModule().Resources["devlake_bitbucketserver_connection.bbserver"]; ok {
						return rs.Primary.ID, nil
					} else {
						return "", fmt.Errorf("Resource devlake_bitbucketserver_connection.bbserver not found in state")
					}
				},
				ImportStateVerify: true,
				// The last_updated attribute does exist in the devlake API, but
				// we want the terraform state here
				ImportStateVerifyIgnore: []string{"password", "last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "devlake_bitbucketserver_connection" "bbserver" {
  endpoint  = "https://bitbucket-server2.org"
  name      = "should_not_exist"
  password  = "whatever"
  username  = "serviceAccount"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.bbserver", "endpoint", "https://bitbucket-server2.org"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.bbserver", "name", "should_not_exist"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.bbserver", "password", "whatever"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.bbserver", "proxy", ""),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.bbserver", "rate_limit_per_hour", "0"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection.bbserver", "username", "serviceAccount"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.bbserver", "id"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.bbserver", "last_updated"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.bbserver", "created_at"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection.bbserver", "updated_at"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

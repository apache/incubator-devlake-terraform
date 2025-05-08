// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	bitbucketServerConnectionScopeConfigConfig = bitbucketServerConnectionConfig + `
resource "devlake_bitbucketserver_connection_scopeconfig" "scopeconf" {
  connection_id	= devlake_bitbucketserver_connection.bbserver.id
  name          = "conf1"
}
`
)

func TestAccBitbucketServerConnectionScopeConfigResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: bitbucketServerConnectionScopeConfigConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "name", "conf1"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "entities.#", "3"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "entities.0", "CODEREVIEW"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "entities.1", "CROSS"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "entities.2", "CODE"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "pr_component", ""),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "pr_type", ""),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "ref_diff.tags_limit", "10"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "ref_diff.tags_pattern", `/v\d+\.\d+(\.\d+(-rc)*\d*)*$/`),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "connection_id"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "created_at"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "id"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "last_updated"),
				),
			},
			// ImportState testing
			{
				ResourceName: "devlake_bitbucketserver_connection_scopeconfig.scopeconf",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					var connectionId, scopeConfigId string
					if con, ok := s.RootModule().Resources["devlake_bitbucketserver_connection.bbserver"]; ok {
						connectionId = con.Primary.ID
					} else {
						return "", fmt.Errorf("Resource devlake_bitbucketserver_connection.bbserver not found in state")
					}
					if scope, ok := s.RootModule().Resources["devlake_bitbucketserver_connection_scopeconfig.scopeconf"]; ok {
						scopeConfigId = scope.Primary.ID
					} else {
						return "", fmt.Errorf("Resource devlake_bitbucketserver_connection_scopeconfig.scopeconf not found in state")
					}
					return fmt.Sprintf("%s,%s", connectionId, scopeConfigId), nil
				},
				ImportStateVerify: true,
				// The last_updated attribute does exist in the devlake API, but
				// we want the terraform state here
				ImportStateVerifyIgnore: []string{"last_updated", "connection_id"},
			},
			// Update and Read testing
			{
				Config: bitbucketServerConnectionConfig + `
resource "devlake_bitbucketserver_connection_scopeconfig" "scopeconf" {
  connection_id	= devlake_bitbucketserver_connection.bbserver.id
  name      = "conf2"
  pr_type	= "type: ([a-zA-Z0-9_-]+)"
  ref_diff  = {
    tags_limit	= 11
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "name", "conf2"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "entities.#", "3"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "entities.0", "CODEREVIEW"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "entities.1", "CROSS"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "entities.2", "CODE"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "pr_component", ""),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "pr_type", "type: ([a-zA-Z0-9_-]+)"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "ref_diff.tags_limit", "11"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "ref_diff.tags_pattern", `/v\d+\.\d+(\.\d+(-rc)*\d*)*$/`),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "connection_id"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "created_at"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "id"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scopeconfig.scopeconf", "last_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

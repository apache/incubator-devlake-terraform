// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	bitbucketServerConnectionScopeConfig = bitbucketServerConnectionScopeConfigConfig + `
resource "devlake_bitbucketserver_connection_scope" "scope" {
  id = "PROJECT/repos/REPO"
  clone_url = "https://bitbucket-server.org/scp/project/repos/repo.git"
  connection_id	= devlake_bitbucketserver_connection.bbserver.id
  description = "example repo"
  html_url = "https://bitbucket-server.org/projects/PROJECT/repos/REPO/browse"
  name = "PROJECT/REPO"
  scope_config_id = devlake_bitbucketserver_connection_scopeconfig.scopeconf.id
}
`
)

func TestAccBitbucketServerConnectionScopeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: bitbucketServerConnectionScopeConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scope.scope", "clone_url", "https://bitbucket-server.org/scp/project/repos/repo.git"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scope.scope", "description", "example repo"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scope.scope", "html_url", "https://bitbucket-server.org/projects/PROJECT/repos/REPO/browse"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scope.scope", "name", "PROJECT/REPO"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scope.scope", "id", "PROJECT/repos/REPO"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scope.scope", "connection_id"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scope.scope", "created_at"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scope.scope", "last_updated"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scope.scope", "scope_config_id"),
				),
			},
			// ImportState testing
			{
				ResourceName: "devlake_bitbucketserver_connection_scope.scope",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					var connectionId, scopeId string
					if con, ok := s.RootModule().Resources["devlake_bitbucketserver_connection.bbserver"]; ok {
						connectionId = con.Primary.ID
					} else {
						return "", fmt.Errorf("Resource devlake_bitbucketserver_connection.bbserver not found in state")
					}
					if scope, ok := s.RootModule().Resources["devlake_bitbucketserver_connection_scope.scope"]; ok {
						scopeId = scope.Primary.ID
					} else {
						return "", fmt.Errorf("Resource devlake_bitbucketserver_connection_scope.scope not found in state")
					}
					return fmt.Sprintf("%s,%s", connectionId, scopeId), nil
				},
				ImportStateVerify: true,
				// The last_updated attribute does exist in the devlake API, but
				// we want the terraform state here
				ImportStateVerifyIgnore: []string{"last_updated", "connection_id", "scope_config_id", "created_at"},
			},
			// Update and Read testing
			{
				Config: bitbucketServerConnectionScopeConfigConfig + `
resource "devlake_bitbucketserver_connection_scope" "scope" {
  id = "PROJECT/repos/REPO"
  clone_url = "https://bitbucket-server.org/scp/project/repos/repo.git"
  connection_id	= devlake_bitbucketserver_connection.bbserver.id
  description = "example repo"
  html_url = "https://bitbucket-server.org/projects/PROJECT/repos/REPO2/browse"
  name = "PROJECT/REPO2"
  scope_config_id = devlake_bitbucketserver_connection_scopeconfig.scopeconf.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scope.scope", "clone_url", "https://bitbucket-server.org/scp/project/repos/repo.git"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scope.scope", "description", "example repo"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scope.scope", "html_url", "https://bitbucket-server.org/projects/PROJECT/repos/REPO2/browse"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scope.scope", "name", "PROJECT/REPO2"),
					resource.TestCheckResourceAttr("devlake_bitbucketserver_connection_scope.scope", "id", "PROJECT/repos/REPO"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scope.scope", "connection_id"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scope.scope", "created_at"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scope.scope", "last_updated"),
					resource.TestCheckResourceAttrSet("devlake_bitbucketserver_connection_scope.scope", "scope_config_id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

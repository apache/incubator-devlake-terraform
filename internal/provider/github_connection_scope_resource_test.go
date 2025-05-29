// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	githubConnectionScopeConfig = githubConnectionScopeConfigConfig + `
resource "devlake_github_connection_scope" "scope" {
  full_name = "PROJECT/REPO"
  id = "42"
  connection_id	= devlake_github_connection.gh.id
  description = "example repo"
  scope_config_id = devlake_github_connection_scopeconfig.scopeconf.id
}
`
)

func TestAccGithubConnectionScopeResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: githubConnectionScopeConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devlake_github_connection_scope.scope", "description", "example repo"),
					resource.TestCheckResourceAttr("devlake_github_connection_scope.scope", "full_name", "PROJECT/REPO"),
					resource.TestCheckResourceAttr("devlake_github_connection_scope.scope", "id", "42"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devlake_github_connection_scope.scope", "connection_id"),
					resource.TestCheckResourceAttrSet("devlake_github_connection_scope.scope", "created_at"),
					resource.TestCheckResourceAttrSet("devlake_github_connection_scope.scope", "last_updated"),
					resource.TestCheckResourceAttrSet("devlake_github_connection_scope.scope", "scope_config_id"),
				),
			},
			// ImportState testing
			{
				ResourceName: "devlake_github_connection_scope.scope",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					var connectionId, scopeId string
					if con, ok := s.RootModule().Resources["devlake_github_connection.gh"]; ok {
						connectionId = con.Primary.ID
					} else {
						return "", fmt.Errorf("Resource devlake_github_connection.gh not found in state")
					}
					if scope, ok := s.RootModule().Resources["devlake_github_connection_scope.scope"]; ok {
						scopeId = scope.Primary.ID
					} else {
						return "", fmt.Errorf("Resource devlake_github_connection_scope.scope not found in state")
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
				Config: githubConnectionScopeConfigConfig + `
resource "devlake_github_connection_scope" "scope" {
  full_name = "PROJECT/REPO"
  id = "42"
  connection_id	= devlake_github_connection.gh.id
  description = "example repo desc"
  scope_config_id = devlake_github_connection_scopeconfig.scopeconf.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devlake_github_connection_scope.scope", "description", "example repo desc"),
					resource.TestCheckResourceAttr("devlake_github_connection_scope.scope", "full_name", "PROJECT/REPO"),
					resource.TestCheckResourceAttr("devlake_github_connection_scope.scope", "id", "42"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devlake_github_connection_scope.scope", "connection_id"),
					resource.TestCheckResourceAttrSet("devlake_github_connection_scope.scope", "created_at"),
					resource.TestCheckResourceAttrSet("devlake_github_connection_scope.scope", "last_updated"),
					resource.TestCheckResourceAttrSet("devlake_github_connection_scope.scope", "scope_config_id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

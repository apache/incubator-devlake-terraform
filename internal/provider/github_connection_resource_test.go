// Copyright (c) HashiCorp, Inc.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	githubConnectionConfig = providerConfig + `
resource "devlake_github_connection" "gh" {
  name      		= "should_not_exist"
  app_id    		= 123123
  installation_id	= 321321
  secret_key		= "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA8Y******\n******sm3C6hlD0XCuVGG1rPuh\n-----END RSA PRIVATE KEY-----"
}
`
)

func TestAccGithubConnectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: githubConnectionConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "app_id", "123123"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "auth_method", "AppKey"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "enable_graphql", "true"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "endpoint", "https://api.github.com/"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "installation_id", "321321"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "name", "should_not_exist"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "proxy", ""),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "rate_limit_per_hour", "0"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "secret_key", "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA8Y******\n******sm3C6hlD0XCuVGG1rPuh\n-----END RSA PRIVATE KEY-----"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "token", ""),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devlake_github_connection.gh", "id"),
					resource.TestCheckResourceAttrSet("devlake_github_connection.gh", "last_updated"),
					resource.TestCheckResourceAttrSet("devlake_github_connection.gh", "created_at"),
					resource.TestCheckResourceAttrSet("devlake_github_connection.gh", "updated_at"),
				),
			},
			// ImportState testing
			{
				ResourceName: "devlake_github_connection.gh",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					if rs, ok := s.RootModule().Resources["devlake_github_connection.gh"]; ok {
						return rs.Primary.ID, nil
					} else {
						return "", fmt.Errorf("Resource devlake_github_connection.gh not found in state")
					}
				},
				ImportStateVerify: true,
				// The last_updated attribute does exist in the devlake API, but
				// we want the terraform state here
				ImportStateVerifyIgnore: []string{"secret_key", "token", "last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "devlake_github_connection" "gh" {
  name      		= "should_not_exist"
  app_id    		= 42
  installation_id	= 321321
  secret_key		= "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA8Y******\n******sm3C6hlD0XCuVGG1rPuh\n-----END RSA PRIVATE KEY-----"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "app_id", "42"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "auth_method", "AppKey"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "enable_graphql", "true"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "endpoint", "https://api.github.com/"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "installation_id", "321321"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "name", "should_not_exist"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "proxy", ""),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "rate_limit_per_hour", "0"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "secret_key", "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA8Y******\n******sm3C6hlD0XCuVGG1rPuh\n-----END RSA PRIVATE KEY-----"),
					resource.TestCheckResourceAttr("devlake_github_connection.gh", "token", ""),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("devlake_github_connection.gh", "id"),
					resource.TestCheckResourceAttrSet("devlake_github_connection.gh", "last_updated"),
					resource.TestCheckResourceAttrSet("devlake_github_connection.gh", "created_at"),
					resource.TestCheckResourceAttrSet("devlake_github_connection.gh", "updated_at"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

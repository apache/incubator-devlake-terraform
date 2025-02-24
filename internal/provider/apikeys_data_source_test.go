// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApiKeysDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "devlake_apikeys" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.#", "1"),
					resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.0.allowed_path", ".*"),
					resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.0.api_key", ""),
					// resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.0.created_at", "2025-02-27T07:04:48.996Z"),
					resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.0.creator", ""),
					resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.0.creator_email", ""),
					// resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.0.expired_at", "2026-02-28T09:12:00.153Z"),
					resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.0.extra", ""),
					resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.0.id", "1"),
					resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.0.name", "terraform_integration_test"),
					resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.0.type", "devlake"),
					// resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.0.updated_at", "2025-02-27T07:04:48.996Z"),
					resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.0.updater", ""),
					resource.TestCheckResourceAttr("data.devlake_apikeys.test", "apikeys.0.updater_email", ""),
				),
			},
		},
	})
}

# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    devlake = {
      source = "registry.terraform.io/incubator-devlake-terraform/devlake"
    }
  }
}

provider "devlake" {
  host  = "http://localhost:4000/api"
  token = "whatever"
}

resource "devlake_github_connection" "gh" {
  name            = "should_not_exist"
  app_id          = 123123
  installation_id = 321321
  secret_key      = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA8Y******\n******sm3C6hlD0XCuVGG1rPuh\n-----END RSA PRIVATE KEY-----"
}

resource "devlake_github_connection_scopeconfig" "scopeconf" {
  connection_id = devlake_github_connection.gh.id
  name          = "conf1"
}

resource "devlake_github_connection_scope" "scope" {
  full_name       = "PROJECT/REPO"
  id              = "42"
  connection_id   = devlake_github_connection.gh.id
  description     = "example repo"
  scope_config_id = devlake_github_connection_scopeconfig.scopeconf.id
}

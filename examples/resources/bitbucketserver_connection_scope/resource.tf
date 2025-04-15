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

resource "devlake_bitbucketserver_connection" "bbserver" {
  endpoint = "https://bitbucket-server.org"
  name     = "should_not_exist"
  password = "whatever"
  username = "serviceAccount"
}

resource "devlake_bitbucketserver_connection_scopeconfig" "scopeconf" {
  connection_id = devlake_bitbucketserver_connection.bbserver.id
  name          = "conf"
}

resource "devlake_bitbucketserver_connection_scope" "scope" {
  id              = "PROJECT/repos/REPO"
  clone_url       = "https://bitbucket-server.org/scp/project/repos/repo.git"
  connection_id   = devlake_bitbucketserver_connection.bbserver.id
  description     = "example repo"
  html_url        = "https://bitbucket-server.org/projects/PROJECT/repos/REPO/browse"
  name            = "PROJECT/REPO"
  scope_config_id = devlake_bitbucketserver_connection_scopeconfig.scopeconf.id
}

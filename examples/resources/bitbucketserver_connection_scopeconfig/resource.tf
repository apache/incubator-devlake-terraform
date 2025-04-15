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
  name          = "conf2"
}

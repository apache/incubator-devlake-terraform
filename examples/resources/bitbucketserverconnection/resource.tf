# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    devlake = {
      source = "registry.terraform.io/incubator-devlake/devlake"
    }
  }
}

provider "devlake" {
  host  = "http://localhost:8080"
  token = "whatever"
}

resource "devlake_bitbucketserver_connection" "tfresourcename" {
  endpoint = "https://bitbucket-server.org"
  name     = "should_not_exist"
  password = "whatever"
  username = "serviceAccount"
}

output "bitbucketserver_connection_resource" {
  value     = devlake_bitbucketserver_connection.tfresourcename
  sensitive = true
}

# Copyright (c) HashiCorp, Inc.

terraform {
  required_providers {
    devlake = {
      source = "registry.terraform.io/incubator-devlake-terraform/devlake"
    }
  }
}

provider "devlake" {
  host  = "http://localhost:8080"
  token = "whatever"
}

resource "devlake_apikey" "tfresourcename" {
  allowed_path = ".*"
  expired_at   = "2025-02-28T09:12:00.153Z"
  name         = "devlakekeyname"
}

output "single_apikey_resource" {
  value     = devlake_apikey.tfresourcename
  sensitive = true
}

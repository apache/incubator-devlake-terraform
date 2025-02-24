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

data "devlake_apikeys" "all" {}

output "all_apikeys_data_source" {
  value = data.devlake_apikeys.all
}

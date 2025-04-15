#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.


# bitbucketserver connection can be imported by specifying the numeric identifier of the connection and the scopeconfig.
terraform import devlake_bitbucketserver_connection_scopeconfig.scopeconf "1,1"

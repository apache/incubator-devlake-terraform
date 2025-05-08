#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.


# bitbucketserver connection scope can be imported by specifying the connection id and the identifier in the form of '<PROJECT>/repos/<REPO>'.
terraform import devlake_bitbucketserver_connection_scope.scope "1,PROJECT/repos/REPO"

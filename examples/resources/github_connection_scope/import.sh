#!/usr/bin/env bash
# Copyright (c) HashiCorp, Inc.


# github connection scope can be imported by specifying the connection id and the github repo identifier.
terraform import devlake_github_connection_scope.scope "1,42"

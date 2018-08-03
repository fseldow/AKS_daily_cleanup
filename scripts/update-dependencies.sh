#!/bin/bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -e
cd $(dirname "$BASH_SOURCE")/..
export GLIDE_HOME=$(pwd)/.glide

# Install dependencies, first round
glide update --no-recursive
# A first round, update install dependencies
glide update
# A second round, this minimizes package in 'glide.lock'
glide update --strip-vendor
# Clean up, turn on use-lock-file, otherwise test dependencies are removed
glide-vc --use-lock-file --only-code --no-tests > /dev/null
 
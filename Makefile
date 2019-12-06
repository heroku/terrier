#   Copyright (c) 2019, salesforce.com, inc.
#   All rights reserved.
#   SPDX-License-Identifier: Apache License 2.0
#   For full license text, see the LICENSE file in the repo root

all: install all_binaries 

linux_binaries:
	gox -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch="linux/amd64"
	
all_binaries:
	gox -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch="darwin/amd64 linux/386 linux/amd64"

install:
	go build .
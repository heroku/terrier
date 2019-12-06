#!/bin/bash

# Copyright (c) 2019, salesforce.com, inc.
# All rights reserved.
# SPDX-License-Identifier: Apache License 2.0
# For full license text, see the LICENSE file in the repo root
 

usage () {
	echo "Usage: $0 <hash list> <output yml> "
}

if [ $# -ne 2 ];then
	usage
	exit 1
fi

if [ ! -f "$1" ];then
	echo "file $1 not found"
	exit 2
fi


OUTPUT_FILE=$2

echo "Converting $1 to Terrier YML: ${OUTPUT_FILE}"

rm -f $OUTPUT_FILE
cat << EOF > $OUTPUT_FILE
mode: image
#mode: container
image: $2
#path: path/to/container/merged
#verbose: true
#veryverbose: true
#UNCOMMENT AND COMMENT THE VALUES ABOVE ACCORDING TO WHAT YOU WANT TO DO
files:
EOF

input="$1"
while IFS= read -r line
do
  sha256=$(echo "$line" | awk '{ print $1; }')
  fullName=$(echo "$line" | awk '{ print $2; }')
  fileName="${fullName:1}"
  cat << EOF >> $OUTPUT_FILE
  - name: '$fileName'
    hashes:
       - hash: '$sha256'
EOF
done < "${input}"
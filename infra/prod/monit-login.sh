#!/bin/bash
set -eo pipefail

INFO=$(journalctl --unit=ssh --since="2 minutes ago" | grep "Accepted publickey for" ||: | tail -1)

if [[ $INFO ]]; then
    USER=$(echo "${INFO}" | awk '{print $9}')
    IP=$(echo "${INFO}" | awk '{print $11}')
    echo "Login of ${USER} user from ${IP} ip!"

    exit 1
else
    exit 0
fi

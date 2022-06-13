#!/bin/bash

if test -f "/etc/systemd/system/solocms.service"; then
    systemctl stop solocms
    systemctl disable solocms

    systemctl daemon-reload
    systemctl reset-failed
fi

if ! [ -d /var/lib/solocms/ ]; then
    mkdir /var/lib/solocms
fi
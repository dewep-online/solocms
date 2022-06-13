#!/bin/bash

#!/bin/bash

if [ -f "/etc/systemd/system/solocms.service" ]; then
    systemctl start solocms
    systemctl enable solocms
fi
#!/bin/bash
# set -e
/usr/share/nginx/html/env_load.sh &&
# sqlite_web -H 0.0.0.0 -x /usr/src/app/db/currencyData.sqlite -d true -u /db &
# sqlite_web -H 0.0.0.0 -p 8085 -x /usr/src/app/db/data-go.sqlite -d true -u /dbgo &
/server &
nginx -g "daemon off;"
# wait -n
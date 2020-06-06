#!/bin/bash
# set -e
/usr/share/nginx/html/env_load.sh &&
sqlite_web -H 0.0.0.0 -x /usr/src/app/data.sqlite -d true -u /db &
(cd /usr/src/app/backend; node index.js) &
nginx -g "daemon off;"
# wait -n
#!/bin/bash
# set -e

sqlite_web -H 0.0.0.0 -x test.sqlite -d true -u /db &
(cd /usr/src/app/backend; node index.js) &
nginx -g "daemon on;"
# wait -n
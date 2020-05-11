#!/bin/bash
# set -e

sqlite_web -H 0.0.0.0 -x test.sqlite -d true -u /db &
# npm run dev
(cd /usr/src/app/frontend; npm run serve) &
(cd /usr/src/app/backend; npm run server)
# wait -n
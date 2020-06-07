#!/bin/bash
# set -e

# go get github.com/derekparker/delve/cmd/dlv && 
# sh /usr/src/app/backend_go/start_script.sh &
# sqlite_web -H 0.0.0.0 -x data.sqlite -d true -u /db &
sqlite_web -H 0.0.0.0 -x ./db/currencyData.sqlite -d true -u /db &
sqlite_web -H 0.0.0.0 -p 8085 -x ./db/data-go.sqlite -d true -u /dbgo &
# npm run dev
(cd /usr/src/app/frontend; npm run serve)
# (cd /usr/src/app/frontend; npm run serve) &
# (cd /usr/src/app/backend; npm run server)
# wait -n
#!/bin/bash
# set -e

sqlite_web -H 0.0.0.0 -x test.sqlite -d true -u /db &
npm run dev
# wait -n
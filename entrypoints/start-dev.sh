#!/bin/bash
# set -e

# Read each line in .env file
# Each line represents key=value pairs
while read -r line || [[ -n "$line" ]];
do
  # Split env variables by character `=`
  if printf '%s\n' "$line" | grep -q -e '='; then
    varname=$(printf '%s\n' "$line" | sed -e 's/=.*//')
    varvalue=$(printf '%s\n' "$line" | sed -e 's/^[^=]*=//')
  fi

  # Read value of current variable if exists as Environment variable
  value=$(printf '%s\n' "${!varname}")
  # Otherwise use value from .env file
  [[ -z $value ]] && value=${varvalue}
  
  # Set new VUE_APP variables with existing variables
  newval=$(printf 'VUE_APP_%s=%s' "${varname}" "${value}")
  echo $newval
  export $newval
done < .env

sqlite_web -H 0.0.0.0 -x ./db/currencyData.sqlite -d true -u /db &
sqlite_web -H 0.0.0.0 -p 8085 -x ./db/data-go.sqlite -d true -u /dbgo &
(cd /usr/src/app/frontend; npm run serve)
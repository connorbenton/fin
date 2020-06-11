#!/bin/bash
# set -e

# go get github.com/derekparker/delve/cmd/dlv && 
# sh /usr/src/app/backend_go/start_script.sh &
# sqlite_web -H 0.0.0.0 -x data.sqlite -d true -u /db &
# /bin/bash -c "source /dev_env_set.sh" &&
# . /dev_env_set.sh &&

# Read each line in .env file
# Each line represents key=value pairs
while read -r line || [[ -n "$line" ]];
do
  # Filter out all unwanted variables for client
  # if [[ ! "$line" =~ ^BASE_URL* ]]; then
  #   continue
  # fi
  # Split env variables by character `=`
  if printf '%s\n' "$line" | grep -q -e '='; then
    varname=$(printf '%s\n' "$line" | sed -e 's/=.*//')
    varvalue=$(printf '%s\n' "$line" | sed -e 's/^[^=]*=//')
  fi

  # Read value of current variable if exists as Environment variable
  value=$(printf '%s\n' "${!varname}")
  # Otherwise use value from .env file
  [[ -z $value ]] && value=${varvalue}
  
  # Append configuration property to JS file
  # echo "  $varname: \"$value\"," >> ./env-config.js
  # Set new VUE_APP variables with existing variables
  # newval=$(printf 'VUE_APP_%s=\"%s\"' "${varname}" "${value}")
  newval=$(printf 'VUE_APP_%s=%s' "${varname}" "${value}")
  echo $newval
  export $newval
  # echo $newval >> /etc/environment
  # export $newval
  # export VUE_APP_$varname=\"$value\"
done < .env

sqlite_web -H 0.0.0.0 -x ./db/currencyData.sqlite -d true -u /db &
sqlite_web -H 0.0.0.0 -p 8085 -x ./db/data-go.sqlite -d true -u /dbgo &
# npm run dev
(cd /usr/src/app/frontend; npm run serve)
# (cd /usr/src/app/frontend; npm run serve) &
# (cd /usr/src/app/backend; npm run server)
# wait -n
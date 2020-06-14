#!/bin/bash

# Recreate config file
rm -rf ./env-config.js
touch ./env-config.js

# Add assignment 
echo "window._env_ = {" >> ./env-config.js

printenv >> ./temp.env

# Read each line in .env file
# Each line represents key=value pairs
while read -r line || [[ -n "$line" ]];
do
  # Filter out all unwanted variables for client
  if [[ ! "$line" =~ ^PLAID_PUBLIC_KEY*|^PLAID_ENVIRONMENT*|^USE_PLAID*|^USE_SALTEDGE*|^BASE_CURRENCY* ]]; then
    continue
  fi
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
  echo "  $varname: \"$value\"," >> ./env-config.js
done < ./temp.env
rm -rf ./temp.env

echo "}" >> ./env-config.js
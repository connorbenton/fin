# fin
### Automated expense tracking and analysis

To use this service beyond the limited [Mint CSV](https://help.mint.com/Accounts-and-Transactions/888960591/How-can-I-download-my-transactions.htm) import functionality, you will need development API keys to one or both of [SaltEdge Spectre](https://www.saltedge.com/products/spectre) and [Plaid](https://plaid.com/). 

Users can also format additional transaction data to import according to the [example CSV](https://github.com/connorbenton/fin/example.csv) (Mint CSVs will import without any modification necessary), where the optional column currency_code can be populated with the transaction currency (defaults to USD if left blank/column not present). Imports should always be done after accounts are linked and transactions fetched from APIs, because the server will try to identify duplicate transactions during import in order to associate imported transactions with already-existing accounts.

Data is displayed in the analysis tab according to the amount of each transaction normalized to the 'base currency' selected in the [.env file](#example-env-file), using daily exchange rates pulled from the [ECB SDMX API](https://sdw-wsrest.ecb.europa.eu/).

The [development build](#development-environment) also includes two instances of [sqlite-web](https://github.com/coleifer/sqlite-web) running in the frontend to view the two databases (one for currency rates, the other for the main connections/accounts/transactions data).

## Example docker-compose
A simple docker-compose port mapping to the container's 9028 port can be used if preferred for proxying (instead of the docker-networking shown in this example). The only inputs required are an [.env file](#example-env-file) path and a volume mount where the two databases will be stored.
```
version: '3.3'
services:
    fin:
        container_name: fin
        image: connorbenton/fin
        env_file: 
          - {path_to_your_env_file}
        volumes:
          - {path_to_your_db_directory}:/usr/src/app/db
        network_mode: {YOUR_PROXY_NETWORK_HERE} 
        expose:
            - '9028'
```


## Example nginx proxy setup
My personal setup uses [Authelia](https://github.com/authelia/authelia) to secure all proxied endpoints - in principle, any form of auth/protection will work with Fin (the backend and frontend are both internal to the container). Testing has only been done with subdomain (and not subfolder) routing.
```
server {
    server_name fin.{DOMAIN}.{TLD};
    listen 443 ssl http2;
    include /etc/nginx/snippets/strong-ssl.conf;
    include /etc/nginx/authelia.conf;
    client_max_body_size 50M;

    location / {
          set $upstream_fin fin;
          proxy_pass http://$upstream_fin:9028;
                include /etc/nginx/auth.conf;
                include /etc/nginx/proxy.conf;
    }
}
```

## Example env file
```
# Plaid credentials go here (don't use quotes)
PLAID_CLIENT_ID=XXX
PLAID_PUBLIC_KEY=XXX
PLAID_SECRET_DEVELOPMENT=XXX
PLAID_SECRET_SANDBOX=XXX

# SaltEdge credentials go here (don't use quotes)
SALTEDGE_APP_ID=XXX
SALTEDGE_APP_SECRET=XXX
SALTEDGE_CUSTOMER_ID=XXX

# Your base URL goes here, used for SaltEdge redirects on success
BASE_URL=https://SUBDOMAIN.DOMAIN.TLD

# This can be either 'development' or 'sandbox' (no quotes), affects Plaid only
PLAID_ENVIRONMENT=sandbox

# This is an all-caps string for currencies reported by the ECB - pick one from:
# (https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html)
BASE_CURRENCY=USD

# Values can be 'TRUE' or 'FALSE' to turn off either API if you don't need it
USE_SALTEDGE=TRUE
USE_PLAID=TRUE
```

## Logging
Mostly covering the Go backend:
```
$ docker logs -f fin 
```

# Development environment

```
$ git clone https://github.com/connorbenton/fin
```
Once you've edited the docker-compose.yml in ./docker/dev to expose/access the frontend on your own network:
```
$ cd docker/dev
$ docker-compose up -d --build
```
To actually start the backend (the Go backend development environment is built around [remote debugging with Visual Studio Code](https://github.com/golang/vscode-go/blob/master/docs/debugging.md)), you will need to find your fin-go docker IP address with:

```
$ docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' fin-go
```
Then VSCode can be set up to attach to the Docker container:
```
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Delve into Docker",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "remotePath": "/usr/src/app/backend_go",
            "port": 40000,
            "host": {YOUR-FIN-GO-DOCKER-IP-HERE},
            "cwd": "${workspaceRoot}/backend_go",
            "showLog": true
      }
    ]
}
```
It's normally a good idea to watch the backend during development to monitor compile errors and watch for any crashes:
```
$ docker logs -f fin-go
```
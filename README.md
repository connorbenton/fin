# fin
### Automated expense tracking and analysis

## Example docker-compose
```
version: '3.3'
services:
    fin:
        container_name: fin-stage
        image: fin-stage
        env_file: 
          - {path_to_your_env_file}
        volumes:
          - {path_to_your_db_directory}:/usr/src/app/db
        network_mode: {YOUR_PROXY_NETWORK_HERE} 
        expose:
            - '9028'
```
A simple docker-compose port mapping to the container's 9028 port can be used if preferred for proxying (over docker-networking shown in the example above)

### Compiles and hot-reloads for development
```
npm run serve
```

### Compiles and minifies for production
```
npm run build
```

### Run your tests
```
npm run test
```

### Lints and fixes files
```
npm run lint
```

### Customize configuration
See [Configuration Reference](https://cli.vuejs.org/config/).

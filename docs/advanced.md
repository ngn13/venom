## Environment Options
- `VENOM_DEBUG`: Enable debug output for the builds (agents)
- `VENOM_DISABLE_ANTIVM`: Disable Anti-VM feature for the builds
- `VENOM_DISABLE_ANTIDEBUG`: Disable Anti-Debug feature for the builds
- `VENOM_ALLINT`: Listen on all the interfaceses (0.0.0.0)
- `VENOM_URL`: URL for the agent access

## Nginx reverse proxy setup
If you want to use the server with a domain or with SSL, then you should use 
a reverse proxy, such as nginx. Here is an example configuration for nginx:
```conf
server {
  server_name   <domain>;

  location / {
    proxy_pass http://localhost:<local port>;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection 'upgrade';
    proxy_set_header Host $host;
    proxy_cache_bypass $http_upgrade;
  }
}
```

And for the container, make sure you use a local port, for example:
```bash
docker run -d -p 127.0.0.1:8080:8082 \
    -e VENOM_ALLINT=true             \
    -e VENOM_URL=http://<domain>     \
    -v $PWD/db:/venom/server/db      \
    -v $PWD/tor:/var/lib/tor/venom   \
    ghcr.io/ngn13/venom:latest
```

## Docker compose setup
Here is an example docker compose configuration:
```yaml
version: "3"
services:
    venom:
        image: ghcr.io/ngn13/venom:latest 
        volumes:
            - "./db:/venom/server/db"
            - "./tor:/var/lib/tor/venom"
        ports:
            - 80:8082
        environment:
            - VENOM_ALLINT=true
            - VENOM_URL=http://<ip>
```

# sessionid_server

[https://github.com/powerpuffpenguin/sessionid_server](https://github.com/powerpuffpenguin/sessionid_server)

# run
```
docker run \
    -p 9000:80/tcp \
    -d king011/sessionid_server:tag
```

configure 
```
docker run \
    -p 9000:80/tcp \
    -v your_configure:/opt/server/etc:ro \
    -d king011/sessionid_server:tag
```

# VOLUME

```
docker run \
    -p 9000:80/tcp \
    -v your_logs:/opt/server/data/logs \
    -v your_db:/opt/server/data/db \
    -d king011/sessionid_server:tag
```
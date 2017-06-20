# tls tunneling for fun

window 1
```
docker run -it --rm -v ~/go:/go golang /bin/bash
```

window 2
```
docker exec -it $container_handle /bin/bash
```

then from window 1:
```
cd /go/src/github.com/rosenhouse/tls-tunnel-experiments/
./build

bin/server -address 127.0.0.12:7007
```

and from window 2:
```
cd /go/src/github.com/rosenhouse/tls-tunnel-experiments/

echo "hello" | bin/client -address 127.0.0.12:7007 -clientAddr 127.0.0.11:6004 # or whatever
```

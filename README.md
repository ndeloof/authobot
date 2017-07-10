# Authobot 

![logo](logo.jpg)


Authobot is a simple authorization plugin for docker to prevent some API usages
that we know will expose hosts data and are not required for a legitimate use of docker
from a containerized jenkins build agent.

We prevent
- running container with bind mounts
- running privileged container
- _more to come_

We also provide a whitelist of authorized API URIs, based on docker-pipeline / declarative-pipeline requirements.
Any other API call will be rejected.

## Samples 

```
➜ #
➜ # let's run a container 
➜ #
➜ docker run --rm -t  ubuntu echo Hello
Hello

➜ #
➜ # let's now run a _privileged_ container (can access hosts' devices)
➜ #
➜ docker run --rm -t --privileged ubuntu echo Hello
docker: Error response from daemon: authorization denied by plugin authobot:latest: use of Privileged contianers is not allowed.
See 'docker run --help'.

➜ #
➜ # hum, let's bind mount host filesystem to hack all it's secrets
➜ #
➜ docker run --rm -t -v /:/host ubuntu echo Hello
docker: Error response from daemon: authorization denied by plugin authobot:latest: use of bind mounts is not allowed.
See 'docker run --help'.

➜ #
➜ # ok, let's just mount a volume to persist our data
➜ #
➜ docker run --rm -t -v some_volume:/host ubuntu echo Hello
Hello
```



## Contribute / hack

we use [dep](https://github.com/golang/dep) to manage dependencies.
run `dep ensure` to generate a local `vendor` folder so you can hack and build the plugin.

build with `docker build -t authobot .`

create rootfs directory, and export container filesystem
```
rm -rf rootfs
mkdir rootfs
ID=$(docker run -d authobot:latest)
docker export $ID | tar -x -C rootfs
docker kill $ID
docker rm $ID
```

Then, install and enable plugin on your local docker daemon
```
docker plugin create authobot $(pwd)
docker plugin enable authobot
``` 

change docker daemon configuration to include `--authorization-plugin=authobot` and restart daemon.

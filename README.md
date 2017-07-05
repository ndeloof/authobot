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

## Contribute / hack

we use [dep](https://github.com/golang/dep) to manage dependencies.
run `dep ensure` to generate a local `vendor` folder so you can hack and build the plugin.

build with `docker build -t authobot .`

create rootfs directory, and export container filesystem
```
rm -rf rootfs
mkdir rootfs
ID=$(docker run -d authobot)
docker export $ID | tar -x -C rootfs
docker kill $ID
docker rm $ID
```

Then, install and enable plugin on your local docker daemon
```
docker plugin create authobot $(pwd)
docker plugin enable authobot
``` 

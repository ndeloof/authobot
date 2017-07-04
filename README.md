# Authobot 

Authobot is a simple authorization plugin for docker to prevent some API usages
that we know will expose hosts data and are not required for a legitimate use of docker
from a containerized jenkins build agent.

We prevent
- running container with bind mounts
- running privileged container
- _more to come_

![logo](logo.jpg)




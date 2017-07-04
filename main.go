package main

import (
	"fmt"
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/authorization"
)

const (
	defaultDockerHost = "unix:///var/run/docker.sock"
)

var (
	flDockerHost = flag.String("host", defaultDockerHost, "Specifies the host where to contact the docker daemon")
)

func main() {
	fmt.Println("hello")

	flag.Parse()

	authobot, err := newPlugin(*flDockerHost)
	if err != nil {
		logrus.Fatal(err)
	}

	h := authorization.NewHandler(authobot)

	if err := h.ServeUnix("authobot", 0); err != nil {
		logrus.Fatal(err)
	}
}
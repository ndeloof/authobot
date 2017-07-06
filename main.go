package main

import (
	"fmt"
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/authorization"
)

func main() {
	fmt.Println("hello")

	flag.Parse()

	authobot, err := newPlugin()
	if err != nil {
		logrus.Fatal(err)
	}

	h := authorization.NewHandler(authobot)

	if err := h.ServeUnix("authobot", 0); err != nil {
		logrus.Fatal(err)
	}
}
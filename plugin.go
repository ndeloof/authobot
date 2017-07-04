package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"

	"github.com/docker/docker/api"
	"github.com/docker/engine-api/client"
	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/mount"
)

func newPlugin(dockerHost string) (*authobot, error) {
	var transport *http.Client
	client, err := client.NewClient(dockerHost, api.DefaultVersion, transport, nil)
	if err != nil {
		return nil, err
	}
	return &authobot{client: client}, nil
}

var (
	create = regexp.MustCompile(`/containers/create`)
)

type authobot struct {
	client *client.Client
}

type configWrapper struct {
	*container.Config
	HostConfig       *container.HostConfig
}

func (p *authobot) AuthZReq(req authorization.Request) authorization.Response {
	ruri, err := url.QueryUnescape(req.RequestURI)
	if err != nil {
		return authorization.Response{Err: err.Error()}
	}
	if req.RequestMethod == "POST" && create.MatchString(ruri) {
		if req.RequestBody != nil {
			body := &configWrapper{}
			if err := json.NewDecoder(bytes.NewReader(req.RequestBody)).Decode(body); err != nil {
				return authorization.Response{Err: err.Error()}
			}
			if len(body.Volumes) > 0 {
				return authorization.Response{Msg: "use of bind mounts is not allowed"}
			}
			if body.HostConfig.Privileged {
				return authorization.Response{Msg: "use of Privileged contianers is not allowed"}
			}
			for _, m := range body.HostConfig.Mounts {
				if m.Type == mount.TypeBind {
					return authorization.Response{Msg: "use of bind mounts is not allowed"}
				}

			}
		}
	}
	return authorization.Response{Allow: true}
}

func (p *authobot) AuthZRes(req authorization.Request) authorization.Response {
	return authorization.Response{Allow: true}
}
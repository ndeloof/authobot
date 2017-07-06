package main

import (
	"bytes"
	"encoding/json"
	"net/url"
	"regexp"

	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/mount"
	"github.com/pkg/errors"
	"fmt"
)

func newPlugin() (*authobot, error) {
	return &authobot{}, nil
}

var (
	create = regexp.MustCompile(`/containers/create`)

	whitelist = []regexp.Regexp{
		*regexp.MustCompile(`/_ping`),
		*regexp.MustCompile(`/version`),

		*regexp.MustCompile(`/containers/create`),
		*regexp.MustCompile(`/containers/.+/start`),
		*regexp.MustCompile(`/containers/.+/stop`),
		*regexp.MustCompile(`/containers/.+/kill`),
		*regexp.MustCompile(`/containers/.+/json`), // inspect
		*regexp.MustCompile(`/containers/.+/exec`),
		*regexp.MustCompile(`/exec/.+/start`),
		*regexp.MustCompile(`/exec/.+/json`),

		*regexp.MustCompile(`/build`),
		*regexp.MustCompile(`/images/create`), // pull
		*regexp.MustCompile(`/images/.+/json`), // inspect
		*regexp.MustCompile(`/images/.+/push`),
		*regexp.MustCompile(`/images/.+/tag`),
		*regexp.MustCompile(`/images/.+`), // remove
	}
)

type authobot struct {
}

type configWrapper struct {
	*container.Config
	HostConfig       *container.HostConfig
}

// --- implement authorization.Plugin

func (p *authobot) AuthZReq(req authorization.Request) authorization.Response {
	uri, err := url.QueryUnescape(req.RequestURI)
	if err != nil {
		return authorization.Response{Err: err.Error()}
	}

	fmt.Println("checking request to "+uri+" from user : "+req.User);

	err = p.Authorized(uri);
	if err != nil {
		return authorization.Response{Err: err.Error()}
	}


	if req.RequestMethod == "POST" && create.MatchString(uri) {
		if req.RequestBody != nil {
			body := &configWrapper{}
			if err := json.NewDecoder(bytes.NewReader(req.RequestBody)).Decode(body); err != nil {
				return authorization.Response{Err: err.Error()}
			}
			if len(body.Volumes) > 0 {
				return authorization.Response{Msg: "use of volumes is not allowed"}
			}
			if body.HostConfig.Privileged {
				return authorization.Response{Msg: "use of Privileged contianers is not allowed"}
			}
			for _, m := range body.HostConfig.Mounts {
				if m.Type == mount.TypeBind {
					return authorization.Response{Msg: "use of bind mounts is not allowed"}
				}

			}
			return authorization.Response{Allow: true}
		}
	}

	return authorization.Response{Msg: uri + " API is not allowed"}
}

func (p *authobot) AuthZRes(req authorization.Request) authorization.Response {
	return authorization.Response{Allow: true}
}

// ---

func (p *authobot) Authorized(uri string) error {
	for _, m := range whitelist {
		if m.MatchString(uri) {
			return nil;
		}
	}
	return errors.New(uri + " is not authorized")
}


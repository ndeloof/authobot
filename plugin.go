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
	"strings"
)

func newPlugin() (*authobot, error) {
	return &authobot{}, nil
}

var (
	create = regexp.MustCompile(`/containers/create`)

	whitelist = []regexp.Regexp{
		*regexp.MustCompile(`/_ping`),
		*regexp.MustCompile(`/v.*/version`),

		*regexp.MustCompile(`/v.*/containers/create`),
		*regexp.MustCompile(`/v.*/containers/.+/start`),
		*regexp.MustCompile(`/v.*/containers/.+/stop`),
		*regexp.MustCompile(`/v.*/containers/.+/kill`),
		*regexp.MustCompile(`/v.*/containers/.+/json`), // inspect
		*regexp.MustCompile(`/v.*/containers/.+/exec`),
		*regexp.MustCompile(`/v.*/exec/.+/start`),
		*regexp.MustCompile(`/v.*/exec/.+/json`),

		*regexp.MustCompile(`/v.*/build`),
		*regexp.MustCompile(`/v.*/images/create`), // pull
		*regexp.MustCompile(`/v.*/images/.+/json`), // inspect
		*regexp.MustCompile(`/v.*/images/.+/push`),
		*regexp.MustCompile(`/v.*/images/.+/tag`),
		*regexp.MustCompile(`/v.*/images/.+`), // remove
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

	// Remove query parameters
	i := strings.Index(uri, "?")
	if i > 0 {
		uri = uri[:i]
	}

	fmt.Println("checking request to "+uri+" from user : "+req.User)

	err = p.Authorized(uri)
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
	fmt.Println("AuthZRes");
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


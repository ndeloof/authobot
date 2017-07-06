package main

import "testing"

func TestDockerVersionWhitelisted(t *testing.T) {
	authobot, _ := newPlugin()
	err := authobot.Authorized("/v1.30/version")
	if err != nil {
		t.Error("/version was not authorized")
	}
}

func TestDockerInspectWhitelisted(t *testing.T) {
	authobot, _ := newPlugin()
	err := authobot.Authorized("/v1.30/containers/1234abcd/json")
	if err != nil {
		t.Error("/containers/1234abcd/json was not authorized")
	}

}

func TestDockerPSBlacklisted(t *testing.T) {
	authobot, _ := newPlugin()
	err := authobot.Authorized("/v1.30/containers/json")
	if err == nil {
		t.Error("/containers/json should have been rejected")
	}
}

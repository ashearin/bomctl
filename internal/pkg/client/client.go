// -----------------------------------------------------------------------------
// SPDX-FileCopyrightText: Copyright © 2024 bomctl a Series of LF Projects, LLC
// SPDX-FileName: internal/pkg/client/client.go
// SPDX-FileType: SOURCE
// SPDX-License-Identifier: Apache-2.0
// -----------------------------------------------------------------------------
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// -----------------------------------------------------------------------------

package client

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/bomctl/bomctl/internal/pkg/client/git"
	"github.com/bomctl/bomctl/internal/pkg/client/github"
	"github.com/bomctl/bomctl/internal/pkg/client/gitlab"
	"github.com/bomctl/bomctl/internal/pkg/client/http"
	"github.com/bomctl/bomctl/internal/pkg/client/oci"
	"github.com/bomctl/bomctl/internal/pkg/netutil"
	"github.com/bomctl/bomctl/internal/pkg/options"
	"github.com/bomctl/bomctl/internal/pkg/sliceutil"
)

var ErrUnsupportedURL = errors.New("failed to parse URL; see `--help` for valid URL patterns")

type (
	Client interface {
		Name() string
	}

	Fetcher interface {
		Client
		Fetch(fetchURL string, opts *options.FetchOptions) ([]byte, error)
		PrepareFetch(URL *url.URL, auth *netutil.BasicAuth, opts *options.Options) error
	}

	Pusher interface {
		Client
		AddFile(pushURL, id string, opts *options.PushOptions) error
		PreparePush(opts *options.PushOptions) error
		Push(opts *options.PushOptions) error
	}
)

// URL: [scheme:][//[userinfo@]host][/]path[?query][#fragment]
//
// https://username:password@github.com:12345/bomctl/bomctl.git#sbom.cdx.json?ref=main
//
// scheme: https
// userinfo: [username, password]
// host: github.com:12345 (parse port with url.Port())
// path: bomctl/bomctl.git
// fragment: sbom.cdx.json
// Query: ref=main (parse query with url.Query() and get value with query.Get("ref"))
//

func NewFetcher(fetchURL string) (Fetcher, error) {
	var fetch Fetcher
	c, err := New(fetchURL)
	fetch = c.(Fetcher)

	if err != nil {
		return nil, fmt.Errorf("Invalid URL: %s", fetchURL)
	}

	return fetch, nil
}

func New(sbomURL string) (Client, error) {
	parsedURL, err := url.Parse(sbomURL)
	if err != nil {
		return nil, fmt.Errorf("Invalid URL: %s", sbomURL)
	}

	return DetermineClient(parsedURL)
}

func DetermineClient(parsedURL *url.URL) (Client, error) {
	fmt.Printf("%#v\n", parsedURL)

	if client := checkScheme(parsedURL); client != nil {
		return client, nil
	}

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "https"
	}

	if strings.Contains(parsedURL.Path, ".git") || parsedURL.Fragment != "" {
		return &git.Client{}, nil
	}

	if strings.Contains(parsedURL.Host, "github") {
		return &github.Client{}, nil
	}

	if strings.Contains(parsedURL.Host, "gitlab") {
		return &gitlab.Client{}, nil
	}

	if sliceutil.Any(registrySlice(), func(s string) bool { return strings.Contains(parsedURL.Host, s) }) {
		return oci.Init(parsedURL)
	}

	if parsedURL.Scheme == "https" || parsedURL.Scheme == "http" {
		return &http.Client{}, nil
	}

	return nil, fmt.Errorf("%w", ErrUnsupportedURL)
}

func checkScheme(parsedURL *url.URL) Client {
	// check for explicit selection via scheme
	switch parsedURL.Scheme {
	case "git":
		return &git.Client{}
	case "github":
		return &github.Client{}
	case "gitlab":
		return &gitlab.Client{}
	case "oci", "oci-archive", "docker", "docker-archive":
		return &oci.Client{}
	default:
		return nil
	}
}

func registrySlice() []string {
	// maybe check for env variable like BOMCTL_REGISTRY or something and append to slice if it exists
	return []string{"registry", "harbor", "ghcr.io", "gcr.io", "docker.io"}
}

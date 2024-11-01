// -----------------------------------------------------------------------------
// SPDX-FileCopyrightText: Copyright © 2024 bomctl a Series of LF Projects, LLC
// SPDX-FileName: internal/pkg/client/github/client.go
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

package github

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/google/go-github/v53/github"

	"github.com/bomctl/bomctl/internal/pkg/netutil"
	"github.com/bomctl/bomctl/internal/pkg/options"
)

type Client struct {
	ghClient        *github.Client
	repo            *github.Repository
	owner, repoName string
}

func (*Client) Name() string {
	return "github"
}

func (*Client) RegExp() *regexp.Regexp {
	return regexp.MustCompile(
		fmt.Sprintf("^%s%s%s%s%s$",
			`((?:github\+)?(?P<scheme>https?|git|ssh):\/\/)?`,
			`((?P<username>[^:]+)(?::(?P<password>[^@]+))?(?:@))?`,
			`(?P<hostname>[^@\/?#:]+)(?::(?P<port>\d+))?`,
			`(?:[\/:](?P<path>[^@#]+))`,
			`(?:@(?P<gitRef>[^#]+))(?:#(?P<fragment>.*))?`,
		),
	)
}

func (client *Client) Parse(rawURL string) *netutil.URL {
	results := map[string]string{}
	pattern := client.RegExp()
	match := pattern.FindStringSubmatch(rawURL)

	for idx, name := range match {
		results[pattern.SubexpNames()[idx]] = name
	}

	if results["scheme"] == "" {
		results["scheme"] = "ssh"
	}

	// Ensure required map fields are present.
	for _, required := range []string{"scheme", "hostname", "path", "gitRef", "fragment"} {
		if value, ok := results[required]; !ok || value == "" {
			return nil
		}
	}

	path := results["path"]
	client.repoName = filepath.Base(path)
	client.owner = filepath.Dir(path)

	return &netutil.URL{
		Scheme:   results["scheme"],
		Username: results["username"],
		Password: results["password"],
		Hostname: results["hostname"],
		Port:     results["port"],
		Path:     results["path"],
		GitRef:   results["gitRef"],
		Query:    results["query"],
		Fragment: results["fragment"],
	}
}

func (client *Client) initRepo(opts *options.Options) error {
	repo, response, err := client.ghClient.Repositories.Get(opts.Context(), client.owner, client.repoName)
	if err != nil {
		return fmt.Errorf("initializing repo %s: %w", response.Body, err)
	}

	client.repo = repo

	return nil
}

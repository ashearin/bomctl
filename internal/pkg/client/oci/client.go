// -----------------------------------------------------------------------------
// SPDX-FileCopyrightText: Copyright © 2024 bomctl a Series of LF Projects, LLC
// SPDX-FileName: internal/pkg/client/oci/client.go
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

package oci

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"
	orasauth "oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"

	"github.com/bomctl/bomctl/internal/pkg/netutil"
	"github.com/bomctl/bomctl/internal/pkg/options"
)

var ErrMultipleSBOMs = errors.New("more than one SBOM document identified in OCI image")

type Client struct {
	TargetURL   *url.URL
	TargetRef   string
	ctx         context.Context
	store       *memory.Store
	repo        *remote.Repository
	descriptors []ocispec.Descriptor
}

func (*Client) Name() string {
	return "OCI"
}

func (*Client) RegExp() *regexp.Regexp {
	return regexp.MustCompile(
		fmt.Sprintf("^%s%s%s%s%s$",
			`((?P<scheme>oci|docker)(?:-archive)?:\/\/)?`,
			`((?P<username>[^:]+)(?::(?P<password>[^@]+))?(?:@))?`,
			`(?P<hostname>[^@\/?#:]+)(?::(?P<port>\d+))?`,
			`(?:\/(?P<path>[^:@]+))`,
			`((?::(?P<tag>[^@]+))|(?:@(?P<digest>sha256:[A-Fa-f0-9]{64})))?`,
		),
	)
}

func Init(targetURL *url.URL) (*Client, error) {
	if targetURL.Scheme == "" {
		targetURL.Scheme = "oci"
	}

	// Ensure required map fields are present.
	if targetURL.Host == "" || targetURL.Path == "" {
		return nil, errors.ErrUnsupported
	}

	client := Client{}
	client.TargetURL = targetURL

	queries := targetURL.Query()

	// Verify ref is present in query
	ref := queries.Get("ref")
	fmt.Println("ref: " + ref)
	if targetURL.RawQuery == "" {
		return nil, errors.ErrUnsupported
	}

	client.TargetRef = ref

	return &client, nil
}

func (client *Client) createRepository(auth *netutil.BasicAuth, opts *options.Options) (err error) {
	client.ctx = opts.Context()
	client.store = memory.New()

	repoPath := (&url.URL{
		Scheme: client.TargetURL.Scheme,
		Host:   client.TargetURL.Host,
		Path:   client.TargetURL.EscapedPath(),
	}).String()

	if client.repo != nil && client.repo.Reference.String() == repoPath {
		return nil
	}

	if client.repo, err = remote.NewRepository(repoPath); err != nil {
		return fmt.Errorf("creating OCI registry repository %s: %w", repoPath, err)
	}

	if auth != nil {
		client.repo.Client = &orasauth.Client{
			Client: retry.DefaultClient,
			Cache:  orasauth.DefaultCache,
			Credential: orasauth.StaticCredential(client.TargetURL.Host, orasauth.Credential{
				Username: auth.Username,
				Password: auth.Password,
			}),
		}
	}

	return nil
}

func descriptorJSON(obj *ocispec.Descriptor) string {
	output, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return ""
	}

	return string(output)
}

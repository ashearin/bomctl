// -----------------------------------------------------------------------------
// SPDX-FileCopyrightText: Copyright © 2024 bomctl a Series of LF Projects, LLC
// SPDX-FileName: internal/pkg/client/oci/internal_test.go
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
	"io"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"oras.land/oras-go/v2/content/memory"
	"oras.land/oras-go/v2/registry/remote"

	"github.com/bomctl/bomctl/internal/pkg/netutil"
	"github.com/bomctl/bomctl/internal/pkg/options"
)

var GetDocument = getDocument

func (client *Client) CreateRepository(url *netutil.URL, auth *netutil.BasicAuth, opts *options.Options) error {
	return client.createRepository(url, auth, opts)
}

func (client *Client) Descriptors() []ocispec.Descriptor {
	return client.descriptors
}

func (client *Client) GenerateManifest(annotations map[string]string) (ocispec.Descriptor, []byte, error) {
	return client.generateManifest(annotations)
}

func (client *Client) PushBlob(
	mediaType string, data io.Reader, annotations map[string]string,
) (ocispec.Descriptor, error) {
	return client.pushBlob(mediaType, data, annotations)
}

func (client *Client) Repo() *remote.Repository {
	return client.repo
}

func (client *Client) Store() *memory.Store {
	return client.store
}

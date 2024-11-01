// -----------------------------------------------------------------------------
// SPDX-FileCopyrightText: Copyright © 2024 bomctl a Series of LF Projects, LLC
// SPDX-FileName: internal/pkg/client/github/client_test.go
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

package github_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/bomctl/bomctl/internal/pkg/client/github"
	"github.com/bomctl/bomctl/internal/pkg/netutil"
)

type githubClientSuite struct {
	suite.Suite
}

func (ghcs *githubClientSuite) TestClient_Parse() {
	client := &github.Client{}

	for _, data := range []struct {
		expected *netutil.URL
		name     string
		url      string
		owner    string
		repoName string
	}{
		{
			name:     "github+http scheme",
			url:      "github+http://github.com/bomctl/bomctl@main#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "http",
				Hostname: "github.com",
				Path:     "bomctl/bomctl",
				GitRef:   "main",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "github+http scheme with tag",
			url:      "github+https://github.com/ossf/scorecard@v0.4.0#sbom.cdx.json",
			owner:    "ossf",
			repoName: "scorecard",
			expected: &netutil.URL{
				Scheme:   "https",
				Hostname: "github.com",
				Path:     "ossf/scorecard",
				GitRef:   "v0.4.0",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "github+http scheme with latest tag",
			url:      "github+http://github.com/bomctl/bomctl@latest#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "http",
				Hostname: "github.com",
				Path:     "bomctl/bomctl",
				GitRef:   "latest",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "github+https scheme with username, port",
			url:      "github+https://git@github.com:12345/bomctl/bomctl@main#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "https",
				Username: "git",
				Hostname: "github.com",
				Port:     "12345",
				Path:     "bomctl/bomctl",
				GitRef:   "main",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "github+https scheme with username, password, port",
			url:      "github+https://username:password@github.com:12345/bomctl/bomctl@main#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "https",
				Username: "username",
				Password: "password",
				Hostname: "github.com",
				Port:     "12345",
				Path:     "bomctl/bomctl",
				GitRef:   "main",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "github+https scheme with username",
			url:      "github+https://git@github.com/bomctl/bomctl@main#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "https",
				Username: "git",
				Hostname: "github.com",
				Path:     "bomctl/bomctl",
				GitRef:   "main",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "ssh scheme",
			url:      "ssh://github.com/bomctl/bomctl@main#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "ssh",
				Hostname: "github.com",
				Path:     "bomctl/bomctl",
				GitRef:   "main",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "ssh scheme with username, port",
			url:      "ssh://git@github.com:12345/bomctl/bomctl@main#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "ssh",
				Username: "git",
				Hostname: "github.com",
				Port:     "12345",
				Path:     "bomctl/bomctl",
				GitRef:   "main",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "ssh scheme with username, password, port",
			url:      "ssh://username:password@github.com:12345/bomctl/bomctl@main#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "ssh",
				Username: "username",
				Password: "password",
				Hostname: "github.com",
				Port:     "12345",
				Path:     "bomctl/bomctl",
				GitRef:   "main",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "ssh scheme with username",
			url:      "ssh://git@github.com/bomctl/bomctl@main#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "ssh",
				Username: "git",
				Hostname: "github.com",
				Path:     "bomctl/bomctl",
				GitRef:   "main",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "git scheme",
			url:      "git://github.com/bomctl/bomctl@main#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "git",
				Hostname: "github.com",
				Path:     "bomctl/bomctl",
				GitRef:   "main",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "git scheme with username, port",
			url:      "git://git@github.com:12345/bomctl/bomctl@main#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "git",
				Username: "git",
				Hostname: "github.com",
				Port:     "12345",
				Path:     "bomctl/bomctl",
				GitRef:   "main",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "git scheme with username, password, port",
			url:      "git://username:password@github.com:12345/bomctl/bomctl@main#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "git",
				Username: "username",
				Password: "password",
				Hostname: "github.com",
				Port:     "12345",
				Path:     "bomctl/bomctl",
				GitRef:   "main",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "git scheme with username",
			url:      "git://git@github.com/bomctl/bomctl@main#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "git",
				Username: "git",
				Hostname: "github.com",
				Path:     "bomctl/bomctl",
				GitRef:   "main",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "git scheme with username",
			url:      "git://git@github.com/bomctl/bomctl@main#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: &netutil.URL{
				Scheme:   "git",
				Username: "git",
				Hostname: "github.com",
				Path:     "bomctl/bomctl",
				GitRef:   "main",
				Fragment: "sbom.cdx.json",
			},
		},
		{
			name:     "missing git ref",
			url:      "github+https://github.com/bomctl/bomctl#sbom.cdx.json",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: nil,
		},
		{
			name:     "missing path to SBOM file",
			url:      "github+https://github.com/bomctl/bomctl@main",
			owner:    "bomctl",
			repoName: "bomctl",
			expected: nil,
		},
	} {
		ghcs.Run(data.name, func() {
			actual := client.Parse(data.url)
			ghcs.Require().Equal(data.expected, actual, data.url)
		})
	}
}

func TestGithubClientSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(githubClientSuite))
}

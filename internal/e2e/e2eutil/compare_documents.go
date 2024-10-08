// -----------------------------------------------------------------------------
// SPDX-FileCopyrightText: Copyright © 2024 bomctl a Series of LF Projects, LLC
// SPDX-FileName: internal/e2e/e2eutil/compare_documents.go
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

package e2eutil

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path"

	"github.com/protobom/protobom/pkg/reader"
	"github.com/protobom/protobom/pkg/sbom"
	"github.com/rogpeppe/go-internal/testscript"
)

const compareDocsRequiredArgNum = 3

// compareDocuments is a testscript command that compares the
// two given protobom documents and checks for equality.
func compareDocuments(script *testscript.TestScript, neg bool, args []string) {
	if len(args) != compareDocsRequiredArgNum {
		script.Fatalf("syntax: compare_docs directory_name file_name1 file_name2")
	}

	documents := []*sbom.Document{}
	sbomReader := reader.New()
	dir := args[0]

	for i := 1; i < len(args); i++ {
		fileString := path.Join(dir, args[i])

		file := getFile(script, fileString)

		data, err := io.ReadAll(file)
		if err != nil {
			script.Fatalf("failed to read from input file: %s", file.Name())
		}

		document, err := sbomReader.ParseStream(bytes.NewReader(data))
		if err != nil {
			script.Fatalf("failed to parse SBOM data from file: %s", file.Name())
		}

		documents = append(documents, document)
	}

	metaMatches := compareMetaData(script, documents[0].GetMetadata(), documents[1].GetMetadata())
	nodeListMatches := documents[0].GetNodeList().Equal(documents[1].GetNodeList())

	reportResult(script, metaMatches, nodeListMatches, neg)
}

func reportResult(script *testscript.TestScript, metaMatches, nodeListMatches, neg bool) { //nolint:revive
	docsMatch := metaMatches && nodeListMatches
	if !docsMatch && !neg {
		if !nodeListMatches {
			script.Logf("node lists do not match")
		}

		script.Fatalf("documents do not match")
	}

	if docsMatch && neg {
		script.Fatalf("documents Match")
	}
}

func getFile(script *testscript.TestScript, filePath string) *os.File {
	if !fileExists(filePath) {
		script.Fatalf("file not found %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		script.Fatalf("failed to open input file: %s", file.Name())
	}

	return file
}

func compareMetaData(script *testscript.TestScript, have, want *sbom.Metadata) bool {
	equal := true

	if have.GetId() != want.GetId() {
		script.Logf("MetaData ID does not match. have %s, want: %s", have.GetId(), want.GetId())

		equal = false
	}

	if have.GetVersion() != want.GetVersion() {
		script.Logf("MetaData Version does not match. have %s, want: %s", have.GetVersion(), want.GetVersion())

		equal = false
	}

	if have.GetName() != want.GetName() {
		script.Logf("MetaData Name does not match. have %s, want: %s", have.GetName(), want.GetName())

		equal = false
	}

	if have.GetName() != want.GetName() {
		script.Logf("MetaData Name does not match. have %s, want: %s", have.GetName(), want.GetName())

		equal = false
	}

	if len(have.GetAuthors()) != len(want.GetAuthors()) {
		script.Logf("MetaData Authors do not match. have %s, want: %s", have.GetAuthors(), want.GetAuthors())

		equal = false
	}

	if len(have.GetDocumentTypes()) != len(want.GetDocumentTypes()) {
		script.Logf("MetaData DocTypes do not match. have %s, want: %s",
			have.GetDocumentTypes(), want.GetDocumentTypes())

		equal = false
	}

	return equal
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)

	return !errors.Is(err, fs.ErrNotExist)
}

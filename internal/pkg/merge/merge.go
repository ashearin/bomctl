// -----------------------------------------------------------------------------
// SPDX-FileCopyrightText: Copyright © 2024 bomctl a Series of LF Projects, LLC
// SPDX-FileName: internal/pkg/merge/merge.go
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

package merge

import (
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/protobom/protobom/pkg/sbom"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bomctl/bomctl/internal/pkg/db"
	"github.com/bomctl/bomctl/internal/pkg/options"
)

var errMissingDocument = errors.New("one or more documents not found")

func Merge(documentIDs []string, opts *options.MergeOptions) (string, error) {
	backend, err := db.BackendFromContext(opts.Context())
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	backend.Logger.Info("Merging documents", "documentIDs", documentIDs)

	// Make document list a map so it can sort by the ids provided
	documentMap, tags, err := getSourceData(backend, documentIDs)
	if err != nil {
		return "", fmt.Errorf("failed to get source data: %w", err)
	}

	tags = append(tags, opts.Tags...)
	slices.Sort(tags)
	tags = slices.Compact(tags)

	merged, err := performTopLevelMerge(documentIDs, documentMap, opts)
	if err != nil {
		return "", err
	}

	if err := mergeRootNodes(merged); err != nil {
		return "", err
	}

	backend.Logger.Info("Adding merged document", "sbomID", merged.GetMetadata().GetId())

	if err := backend.Store(merged, nil); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if opts.Alias != "" {
		if err := backend.SetAlias(merged.GetMetadata().GetId(), opts.Alias, false); err != nil {
			opts.Logger.Warn("Alias could not be set.", "err", err)
		}
	}

	if err := backend.AddDocumentAnnotations(merged.GetMetadata().GetId(), db.TagAnnotation, tags...); err != nil {
		opts.Logger.Warn("Tag(s) could not be set.", "err", err)
	}

	return merged.GetMetadata().GetId(), err
}

func getSourceData(backend *db.Backend, documentIDs []string) (documentMap map[string]*sbom.Document, tags []string,
	err error,
) {
	documents, err := backend.GetDocumentsByIDOrAlias(documentIDs...)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	if len(documents) == 0 {
		return nil, nil, errMissingDocument
	}

	documentMap = make(map[string]*sbom.Document)
	tags = []string{}

	for idx := range documents {
		documentID := documents[idx].GetMetadata().GetId()
		documentIDs[idx] = documentID
		documentMap[documentID] = documents[idx]

		documentTags, err := backend.GetDocumentTags(documentID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get source document tags: %w", err)
		}

		tags = append(tags, documentTags...)
	}

	return documentMap, tags, nil
}

func performTopLevelMerge(sbomIDs []string, documentMap map[string]*sbom.Document,
	opts *options.MergeOptions,
) (*sbom.Document, error) {
	newDocumentID := uuid.New().URN()

	var err error

	merged := sbom.NewDocument()

	if opts.DocumentName != "" {
		merged.Metadata.Name = opts.DocumentName
	}

	merged.Metadata.Id = newDocumentID
	merged.Metadata.Date = timestamppb.Now()
	merged.Metadata.Version = "1"

	for _, sbomID := range sbomIDs {
		// Protobom performs add/update on all nodes in list and re-points edges to nodelist
		err = NewMerger(merged.GetNodeList()).MergeProtoMessage(documentMap[sbomID].GetNodeList())
		if err != nil {
			return nil, fmt.Errorf("failed to merge nodelist: %w", err)
		}

		err = NewMerger(merged.GetMetadata()).MergeProtoMessage(documentMap[sbomID].GetMetadata())
		if err != nil {
			return nil, fmt.Errorf("failed to merge metadata: %w", err)
		}
	}

	return merged, nil
}

func mergeRootNodes(merged *sbom.Document) error {
	var err error

	mergedRootNode := sbom.NewNode()
	mergedRootNode.Id = uuid.New().URN()

	// Augment new root node with other root node data
	for _, root := range merged.GetNodeList().GetRootElements() {
		rootNode := merged.GetNodeList().GetNodeByID(root)

		err = NewMerger(mergedRootNode).MergeProtoMessage(rootNode)
		if err != nil {
			return fmt.Errorf("failed to merge root node: %w", err)
		}
	}

	// Repoint all existing root edges to new root element
	for _, edge := range merged.GetNodeList().GetEdges() {
		if slices.Contains(merged.GetNodeList().GetRootElements(), edge.GetFrom()) {
			edge.From = mergedRootNode.GetId()
		}
	}

	// Reset root node
	oldRootElements := merged.GetNodeList().GetRootElements()
	merged.NodeList.RootElements = nil

	// Add root node first so when the old ones get removed from nodelist, it doesn't reset the edges
	merged.GetNodeList().AddRootNode(mergedRootNode)
	merged.GetNodeList().RemoveNodes(oldRootElements)

	return nil
}

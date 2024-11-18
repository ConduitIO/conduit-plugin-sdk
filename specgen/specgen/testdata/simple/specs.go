// Copyright © 2024 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package simple

import (
	sdk "github.com/conduitio/conduit-connector-sdk"
)

// GlobalConfig is a reusable config struct used in the source and destination
// config.
type GlobalConfig struct {
	// MyGlobalString is a required field in the global config with the name
	// "foo" and default value "bar".
	MyGlobalString string `json:"foo" default:"bar" validate:"required"`
}

type SourceConfig struct {
	sdk.UnimplementedSourceConfig

	// GlobalConfig parameters should be nested under "global". This comment
	// should be ignored.
	Global GlobalConfig `json:"global"`
	/*
		MyInt has a block comment.
		And it spans multiple lines.
	*/
	MyInt    int    `json:"myInt"`
	MyString string `json:"myString"` // MyString has a line comment.
}

var Connector = sdk.Connector{
	NewSpecification: nil,
	NewSource: func() sdk.Source {
		return &Source{}
	},
	NewDestination: nil,
}

type Source struct {
	sdk.UnimplementedSource
	config SourceConfig
}

func (s *Source) Config() sdk.SourceConfig {
	return &s.config
}

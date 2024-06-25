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

package sdk

import (
	"fmt"

	v1 "github.com/conduitio/conduit-connector-protocol/conduit/pschema/v1"
	"github.com/conduitio/conduit-connector-sdk/schema"
)

func initConnectorUtils() error {
	s, err := v1.NewClient()
	if err != nil {
		return fmt.Errorf("failed to initialize schema service client: %w", err)
	}
	schema.Service = s

	return nil
}

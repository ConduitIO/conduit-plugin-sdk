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
	"context"

	"github.com/conduitio/conduit-connector-sdk/internal"
)

// ConnectorIDFromContext fetches the connector ID from the context. If the
// context does not contain a connector ID it returns an empty string.
func ConnectorIDFromContext(ctx context.Context) string {
	connectorID := ctx.Value(internal.ConnectorIDCtxKey{})
	if connectorID != nil {
		return connectorID.(string) //nolint:forcetypeassert // only this package can set the value, it has to be a string
	}
	return ""
}

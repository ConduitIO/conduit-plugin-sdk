// Copyright © 2022 Meroxa, Inc.
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
	"errors"
	"testing"
	"time"

	"github.com/conduitio/conduit-commons/config"
	"github.com/matryer/is"
)

type testConfig struct {
	Foo    string `json:"foo"`
	Bar    int    `json:"bar"`
	Nested struct {
		Baz time.Duration `json:"baz"`
	} `json:"nested"`
	err error
}

func (c *testConfig) Validate(context.Context) error {
	return c.err
}

func TestParseConfig_ValidateCalled(t *testing.T) {
	is := is.New(t)

	wantErr := errors.New("validation error")
	cfg := config.Config{
		"foo": "bar",
	}

	params := config.Parameters{
		"foo": config.Parameter{Type: config.ParameterTypeString},
	}

	target := testConfig{
		err: wantErr,
	}
	err := Util.ParseConfig(context.Background(), cfg, &target, params)
	is.True(errors.Is(err, wantErr))
}

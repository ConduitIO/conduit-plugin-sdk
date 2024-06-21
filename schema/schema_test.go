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

package schema

import (
	"context"
	"errors"
	"testing"

	"github.com/conduitio/conduit-commons/schema"
	pschema "github.com/conduitio/conduit-connector-protocol/conduit/schema"
	"github.com/conduitio/conduit-connector-protocol/conduit/schema/client"
	"github.com/conduitio/conduit-connector-protocol/conduit/schema/mock"
	"github.com/google/go-cmp/cmp"
	"github.com/matryer/is"
	"go.uber.org/mock/gomock"
)

func TestSchemaService_Create_OK(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	schemaBytes := []byte(`
{
  "type": "record",
  "name": "ExampleRecord",
  "fields": [
    {
      "name": "id",
      "type": "int"
    }
  ]
}
`)
	want := schema.Instance{
		ID:      "12345",
		Subject: "schema-name",
		Version: 12,
		Type:    schema.TypeAvro,
		Bytes:   schemaBytes,
	}
	service := mock.NewService(gomock.NewController(t))
	service.EXPECT().
		Create(gomock.Any(), pschema.CreateRequest{
			Subject: "schema-name",
			Type:    pschema.TypeAvro,
			Bytes:   schemaBytes,
		}).
		Return(
			pschema.CreateResponse{
				Instance: want,
			},
			nil,
		)
	underTest, err := NewService(client.WithSchemaService(ctx, service))
	is.NoErr(err)

	got, err := underTest.Create(ctx, schema.TypeAvro, "schema-name", schemaBytes)
	is.NoErr(err)
	is.Equal("", cmp.Diff(want, got))
}

func TestSchemaService_Create_Err(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	schemaBytes := []byte{1, 2, 3}
	service := mock.NewService(gomock.NewController(t))
	serviceErr := errors.New("boom")
	service.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(pschema.CreateResponse{}, serviceErr)
	underTest, err := NewService(client.WithSchemaService(ctx, service))
	is.NoErr(err)

	_, err = underTest.Create(ctx, schema.TypeAvro, "schema-name", schemaBytes)
	is.True(errors.Is(err, serviceErr))
}

func TestSchemaService_Get_OK(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	schemaBytes := []byte(`
{
  "type": "record",
  "name": "ExampleRecord",
  "fields": [
    {
      "name": "id",
      "type": "int"
    }
  ]
}
`)
	want := schema.Instance{
		ID:      "12345",
		Subject: "schema-name",
		Version: 12,
		Type:    schema.TypeAvro,
		Bytes:   schemaBytes,
	}
	service := mock.NewService(gomock.NewController(t))
	service.EXPECT().
		Get(gomock.Any(), pschema.GetRequest{
			Subject: "schema-name",
			Version: 12,
		}).
		Return(
			pschema.GetResponse{
				Instance: want,
			},
			nil,
		)

	underTest, err := NewService(client.WithSchemaService(ctx, service))
	is.NoErr(err)

	got, err := underTest.Get(ctx, "schema-name", 12)
	is.NoErr(err)
	is.Equal("", cmp.Diff(want, got))
}

func TestSchemaService_Get_Err(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	service := mock.NewService(gomock.NewController(t))
	serviceErr := errors.New("boom")
	service.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		Return(pschema.GetResponse{}, serviceErr)

	underTest, err := NewService(client.WithSchemaService(ctx, service))
	is.NoErr(err)

	_, err = underTest.Get(ctx, "schema-name", 12)
	is.True(errors.Is(err, serviceErr))
}
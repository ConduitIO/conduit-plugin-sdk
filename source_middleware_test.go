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
	"errors"
	"testing"
	"time"

	"github.com/conduitio/conduit-commons/config"
	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/conduitio/conduit-commons/schema"
	"github.com/conduitio/conduit-connector-protocol/pconnector"
	"github.com/conduitio/conduit-connector-sdk/internal"
	sdkSchema "github.com/conduitio/conduit-connector-sdk/schema"
	"github.com/google/uuid"
	"github.com/matryer/is"
	"go.uber.org/mock/gomock"
)

func TestSourceWithSchema_Parameters(t *testing.T) {
	is := is.New(t)
	ctrl := gomock.NewController(t)
	src := NewMockSource(ctrl)

	s := SourceWithSchema{}.Wrap(src)

	want := config.Parameters{
		"foo": {
			Default:     "bar",
			Description: "baz",
		},
	}

	src.EXPECT().Parameters().Return(want)
	got := s.Parameters()

	is.Equal(got["foo"], want["foo"])
	is.Equal(len(got), 6) // expected middleware to inject 5 parameters
}

func TestSourceWithSchema_Configure(t *testing.T) {
	ctrl := gomock.NewController(t)
	src := NewMockSource(ctrl)
	ctx := context.Background()

	connectorID := uuid.NewString()
	ctx = internal.Enrich(ctx, pconnector.PluginConfig{ConnectorID: connectorID})
	boolPtr := func(b bool) *bool { return &b }

	testCases := []struct {
		name       string
		middleware SourceWithSchema
		have       config.Config

		wantErr            error
		wantSchemaType     schema.Type
		wantPayloadSubject string
		wantKeySubject     string
	}{{
		name:       "empty config",
		middleware: SourceWithSchema{},
		have:       config.Config{},

		wantSchemaType:     schema.TypeAvro,
		wantPayloadSubject: connectorID + "-payload",
		wantKeySubject:     connectorID + "-key",
	}, {
		name:       "invalid schema type",
		middleware: SourceWithSchema{},
		have: config.Config{
			configSourceSchemaType: "foo",
		},
		wantErr: schema.ErrUnsupportedType,
	}, {
		name: "disabled by default",
		middleware: SourceWithSchema{
			DefaultPayloadEncode: boolPtr(false),
			DefaultKeyEncode:     boolPtr(false),
		},
		have: config.Config{},

		wantSchemaType:     schema.TypeAvro,
		wantPayloadSubject: "",
		wantKeySubject:     "",
	}, {
		name:       "disabled by config",
		middleware: SourceWithSchema{},
		have: config.Config{
			configSourceSchemaPayloadEncode: "false",
			configSourceSchemaKeyEncode:     "false",
		},

		wantSchemaType:     schema.TypeAvro,
		wantPayloadSubject: "",
		wantKeySubject:     "",
	}, {
		name: "static default payload subject",
		middleware: SourceWithSchema{
			DefaultPayloadSubject: "foo",
			DefaultKeySubject:     "bar",
		},
		have: config.Config{},

		wantSchemaType:     schema.TypeAvro,
		wantPayloadSubject: "foo",
		wantKeySubject:     "bar",
	}, {
		name:       "payload subject by config",
		middleware: SourceWithSchema{},
		have: config.Config{
			configSourceSchemaPayloadSubject: "foo",
			configSourceSchemaKeySubject:     "bar",
		},

		wantSchemaType:     schema.TypeAvro,
		wantPayloadSubject: "foo",
		wantKeySubject:     "bar",
	}}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			s := tt.middleware.Wrap(src).(*sourceWithSchema)

			src.EXPECT().Configure(ctx, tt.have).Return(nil)

			err := s.Configure(ctx, tt.have)
			if tt.wantErr != nil {
				is.True(errors.Is(err, tt.wantErr))
				return
			}

			is.NoErr(err)

			is.Equal(s.schemaType, tt.wantSchemaType)
			is.Equal(s.payloadSubject, tt.wantPayloadSubject)
			is.Equal(s.keySubject, tt.wantKeySubject)
		})
	}
}

func TestSourceWithSchema_Read(t *testing.T) {
	is := is.New(t)
	ctrl := gomock.NewController(t)
	src := NewMockSource(ctrl)
	ctx := context.Background()

	connectorID := uuid.NewString()
	ctx = internal.Enrich(ctx, pconnector.PluginConfig{ConnectorID: connectorID})

	s := SourceWithSchema{}.Wrap(src)

	src.EXPECT().Configure(ctx, gomock.Any()).Return(nil)
	err := s.Configure(ctx, config.Config{})
	is.NoErr(err)

	testStructuredData := opencdc.StructuredData{
		"foo":   "bar",
		"int":   1,
		"float": 2.34,
		"time":  time.Now().UTC().Truncate(time.Microsecond), // avro precision is microseconds
	}

	testCases := []struct {
		name   string
		record opencdc.Record
	}{{
		name: "no key, no payload",
		record: opencdc.Record{
			Key: nil,
			Payload: opencdc.Change{
				Before: nil,
				After:  nil,
			},
		},
	}, {
		name: "raw key",
		record: opencdc.Record{
			Key: opencdc.RawData("this should not be encoded"),
			Payload: opencdc.Change{
				Before: nil,
				After:  nil,
			},
		},
	}, {
		name: "structured key",
		record: opencdc.Record{
			Key: testStructuredData.Clone(),
			Payload: opencdc.Change{
				Before: nil,
				After:  nil,
			},
		},
	}, {
		name: "raw payload before",
		record: opencdc.Record{
			Key: nil,
			Payload: opencdc.Change{
				Before: opencdc.RawData("this should not be encoded"),
				After:  nil,
			},
		},
	}, {
		name: "structured payload before",
		record: opencdc.Record{
			Key: nil,
			Payload: opencdc.Change{
				Before: testStructuredData.Clone(),
			},
		},
	}, {
		name: "raw payload after",
		record: opencdc.Record{
			Key: nil,
			Payload: opencdc.Change{
				Before: nil,
				After:  opencdc.RawData("this should not be encoded"),
			},
		},
	}, {
		name: "structured payload after",
		record: opencdc.Record{
			Key: nil,
			Payload: opencdc.Change{
				Before: nil,
				After:  testStructuredData.Clone(),
			},
		},
	}, {
		name: "all structured",
		record: opencdc.Record{
			Key: testStructuredData.Clone(),
			Payload: opencdc.Change{
				Before: testStructuredData.Clone(),
				After:  testStructuredData.Clone(),
			},
		},
	}, {
		name: "all raw",
		record: opencdc.Record{
			Key: opencdc.RawData("this should not be encoded"),
			Payload: opencdc.Change{
				Before: opencdc.RawData("this should not be encoded"),
				After:  opencdc.RawData("this should not be encoded"),
			},
		},
	}, {
		name: "key raw payload structured",
		record: opencdc.Record{
			Key: opencdc.RawData("this should not be encoded"),
			Payload: opencdc.Change{
				Before: nil,
				After:  testStructuredData.Clone(),
			},
		},
	}, {
		name: "key structured payload raw",
		record: opencdc.Record{
			Key: testStructuredData.Clone(),
			Payload: opencdc.Change{
				Before: opencdc.RawData("this should not be encoded"),
				After:  nil,
			},
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			src.EXPECT().Read(ctx).Return(tc.record, nil)

			var wantKey, wantPayloadBefore, wantPayloadAfter opencdc.Data
			if tc.record.Key != nil {
				wantKey = tc.record.Key.Clone()
			}
			if tc.record.Payload.Before != nil {
				wantPayloadBefore = tc.record.Payload.Before.Clone()
			}
			if tc.record.Payload.After != nil {
				wantPayloadAfter = tc.record.Payload.After.Clone()
			}

			got, err := s.Read(ctx)
			is.NoErr(err)

			gotKey := got.Key
			gotPayloadBefore := got.Payload.Before
			gotPayloadAfter := got.Payload.After

			if _, ok := wantKey.(opencdc.StructuredData); ok {
				subject, err := got.Metadata.GetKeySchemaSubject()
				is.NoErr(err)
				version, err := got.Metadata.GetKeySchemaVersion()
				is.NoErr(err)

				sch, err := sdkSchema.Get(ctx, subject, version)
				is.NoErr(err)

				var sd opencdc.StructuredData
				err = sch.Unmarshal(gotKey.Bytes(), &sd)
				is.NoErr(err)

				gotKey = sd
			} else {
				_, err := got.Metadata.GetKeySchemaSubject()
				is.True(errors.Is(err, opencdc.ErrMetadataFieldNotFound))
				_, err = got.Metadata.GetKeySchemaVersion()
				is.True(errors.Is(err, opencdc.ErrMetadataFieldNotFound))
			}

			_, isPayloadBeforeStructured := wantPayloadBefore.(opencdc.StructuredData)
			_, isPayloadAfterStructured := wantPayloadAfter.(opencdc.StructuredData)
			if isPayloadBeforeStructured || isPayloadAfterStructured {
				subject, err := got.Metadata.GetPayloadSchemaSubject()
				is.NoErr(err)
				version, err := got.Metadata.GetPayloadSchemaVersion()
				is.NoErr(err)

				sch, err := sdkSchema.Get(ctx, subject, version)
				is.NoErr(err)

				if isPayloadBeforeStructured {
					var sd opencdc.StructuredData
					err = sch.Unmarshal(gotPayloadBefore.Bytes(), &sd)
					is.NoErr(err)
					gotPayloadBefore = sd
				}
				if isPayloadAfterStructured {
					var sd opencdc.StructuredData
					err = sch.Unmarshal(gotPayloadAfter.Bytes(), &sd)
					is.NoErr(err)
					gotPayloadAfter = sd
				}
			} else {
				_, err := got.Metadata.GetPayloadSchemaSubject()
				is.True(errors.Is(err, opencdc.ErrMetadataFieldNotFound))
				_, err = got.Metadata.GetPayloadSchemaVersion()
				is.True(errors.Is(err, opencdc.ErrMetadataFieldNotFound))
			}

			is.Equal(gotKey, wantKey)
			is.Equal(gotPayloadBefore, wantPayloadBefore)
			is.Equal(gotPayloadAfter, wantPayloadAfter)
		})
	}
}

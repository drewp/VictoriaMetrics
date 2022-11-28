// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opentelemetrypb

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/gogo/protobuf/proto"
)

const spanIDSize = 8

var (
	errMarshalSpanID   = errors.New("marshal: invalid buffer length for SpanID")
	errUnmarshalSpanID = errors.New("unmarshal: invalid SpanID length")
)

// SpanID is a custom data type that is used for all span_id fields in OTLP
// Protobuf messages.
type SpanID [spanIDSize]byte

var _ proto.Sizer = (*SpanID)(nil)

// Size returns the size of the data to serialize.
func (sid SpanID) Size() int {
	if sid.IsEmpty() {
		return 0
	}
	return spanIDSize
}

// IsEmpty returns true if id contains at least one non-zero byte.
func (sid SpanID) IsEmpty() bool {
	return sid == [spanIDSize]byte{}
}

// MarshalTo converts trace ID into a binary representation. Called by Protobuf serialization.
func (sid SpanID) MarshalTo(data []byte) (n int, err error) {
	if sid.IsEmpty() {
		return 0, nil
	}

	if len(data) < spanIDSize {
		return 0, errMarshalSpanID
	}

	return copy(data, sid[:]), nil
}

// Unmarshal inflates this trace ID from binary representation. Called by Protobuf serialization.
func (sid *SpanID) Unmarshal(data []byte) error {
	if len(data) == 0 {
		*sid = [spanIDSize]byte{}
		return nil
	}

	if len(data) != spanIDSize {
		return errUnmarshalSpanID
	}

	copy(sid[:], data)
	return nil
}

// MarshalJSON converts SpanID into a hex string enclosed in quotes.
func (sid SpanID) MarshalJSON() ([]byte, error) {
	if sid.IsEmpty() {
		return []byte(`""`), nil
	}
	return marshalJSON(sid[:])
}

// UnmarshalJSON decodes SpanID from hex string, possibly enclosed in quotes.
// Called by Protobuf JSON deserialization.
func (sid *SpanID) UnmarshalJSON(data []byte) error {
	*sid = [spanIDSize]byte{}
	return unmarshalJSON(sid[:], data)
}

const traceIDSize = 16

var (
	errMarshalTraceID   = errors.New("marshal: invalid buffer length for TraceID")
	errUnmarshalTraceID = errors.New("unmarshal: invalid TraceID length")
)

// TraceID is a custom data type that is used for all trace_id fields in OTLP
// Protobuf messages.
type TraceID [traceIDSize]byte

var _ proto.Sizer = (*SpanID)(nil)

// Size returns the size of the data to serialize.
func (tid TraceID) Size() int {
	if tid.IsEmpty() {
		return 0
	}
	return traceIDSize
}

// IsEmpty returns true if id contains at leas one non-zero byte.
func (tid TraceID) IsEmpty() bool {
	return tid == [traceIDSize]byte{}
}

// MarshalTo converts trace ID into a binary representation. Called by Protobuf serialization.
func (tid TraceID) MarshalTo(data []byte) (n int, err error) {
	if tid.IsEmpty() {
		return 0, nil
	}

	if len(data) < traceIDSize {
		return 0, errMarshalTraceID
	}

	return copy(data, tid[:]), nil
}

// Unmarshal inflates this trace ID from binary representation. Called by Protobuf serialization.
func (tid *TraceID) Unmarshal(data []byte) error {
	if len(data) == 0 {
		*tid = [traceIDSize]byte{}
		return nil
	}

	if len(data) != traceIDSize {
		return errUnmarshalTraceID
	}

	copy(tid[:], data)
	return nil
}

// MarshalJSON converts trace id into a hex string enclosed in quotes.
func (tid TraceID) MarshalJSON() ([]byte, error) {
	if tid.IsEmpty() {
		return []byte(`""`), nil
	}
	return marshalJSON(tid[:])
}

// UnmarshalJSON inflates trace id from hex string, possibly enclosed in quotes.
// Called by Protobuf JSON deserialization.
func (tid *TraceID) UnmarshalJSON(data []byte) error {
	*tid = [traceIDSize]byte{}
	return unmarshalJSON(tid[:], data)
}

// marshalJSON converts trace id into a hex string enclosed in quotes.
// Called by Protobuf JSON deserialization.
func marshalJSON(id []byte) ([]byte, error) {
	// 2 chars per byte plus 2 quote chars at the start and end.
	hexLen := 2*len(id) + 2

	b := make([]byte, hexLen)
	hex.Encode(b[1:hexLen-1], id)
	b[0], b[hexLen-1] = '"', '"'

	return b, nil
}

// unmarshalJSON inflates trace id from hex string, possibly enclosed in quotes.
// Called by Protobuf JSON deserialization.
func unmarshalJSON(dst []byte, src []byte) error {
	if l := len(src); l >= 2 && src[0] == '"' && src[l-1] == '"' {
		src = src[1 : l-1]
	}
	nLen := len(src)
	if nLen == 0 {
		return nil
	}

	if len(dst) != hex.DecodedLen(nLen) {
		return errors.New("invalid length for ID")
	}

	_, err := hex.Decode(dst, src)
	if err != nil {
		return fmt.Errorf("cannot unmarshal ID from string '%s': %w", string(src), err)
	}
	return nil
}

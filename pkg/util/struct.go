package util

import (
	"bytes"
	"fmt"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
)

// MessageToStruct encodes a protobuf Message into a Struct. Hilariously, it
// uses JSON as the intermediary
func MessageToStruct(msg proto.Message) (*types.Struct, error) {
	if msg == nil {
		return nil, fmt.Errorf("nil message")
	}

	buf := &bytes.Buffer{}
	if err := (&jsonpb.Marshaler{OrigName: true}).Marshal(buf, msg); err != nil {
		return nil, err
	}

	pbs := &types.Struct{}
	if err := jsonpb.Unmarshal(buf, pbs); err != nil {
		return nil, err
	}

	return pbs, nil
}

// StructToMessage decodes a protobuf Message from a Struct.
func StructToMessage(pbst *types.Struct, out proto.Message) error {
	if pbst == nil {
		return fmt.Errorf("nil struct")
	}

	buf := &bytes.Buffer{}
	if err := (&jsonpb.Marshaler{OrigName: true}).Marshal(buf, pbst); err != nil {
		return err
	}

	return jsonpb.Unmarshal(buf, out)
}

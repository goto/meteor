package utils

import (
	"errors"
	"fmt"
	"reflect"

	v1beta2 "github.com/goto/meteor/models/gotocompany/assets/v1beta2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// GetAttributes returns custom properties of the given asset
func GetAttributes(asset *v1beta2.Asset) map[string]interface{} {
	msg, err := anypb.UnmarshalNew(asset.Data, proto.UnmarshalOptions{})
	if err != nil {
		return make(map[string]interface{})
	}

	attrMsg, ok := msg.(interface{ GetAttributes() *structpb.Struct })
	if !ok {
		return make(map[string]interface{})
	}

	return attrMsg.GetAttributes().AsMap()
}

// SetAttributes sets custom properties of the given asset
func SetAttributes(asset *v1beta2.Asset, customFields map[string]interface{}) (res *v1beta2.Asset, err error) {
	msg, err := anypb.UnmarshalNew(asset.Data, proto.UnmarshalOptions{})
	if err != nil {
		return nil, fmt.Errorf("unmarshal asset data: %w", err)
	}

	newAttrs, err := structpb.NewStruct(customFields)
	if err != nil {
		return nil, fmt.Errorf("error transforming map into structpb: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("could not find matching model: %v", r)
		}
	}()

	f := reflect.ValueOf(msg).Elem().FieldByName("Attributes")
	if !f.IsValid() || !f.CanSet() {
		return nil, errors.New("could not find matching model")
	}

	f.Set(reflect.ValueOf(newAttrs))

	data, err := anypb.New(msg)
	if err != nil {
		return nil, fmt.Errorf("error transforming msg into anypb: %w", err)
	}

	asset.Data = data

	return asset, nil
}

// TryParseMapToProto parses src into a *structpb.Struct.
// Returns a non-nil empty struct and nil error when src is nil.
// Callers must inspect the returned error before using the struct.
func TryParseMapToProto(src map[string]interface{}) (*structpb.Struct, error) {
	if src == nil {
		return &structpb.Struct{}, nil
	}
	return structpb.NewStruct(src)
}

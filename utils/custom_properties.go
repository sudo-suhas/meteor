package utils

import (
	"github.com/odpf/meteor/models"
	"github.com/odpf/meteor/models/odpf/assets"
	"github.com/odpf/meteor/models/odpf/assets/facets"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/structpb"
)

// GetCustomProperties returns custom properties of the given asset
func GetCustomProperties(metadata models.Metadata) map[string]interface{} {
	customProps := metadata.GetProperties()

	// if data's custom facet is nil, return new empty custom properties
	if customProps == nil {
		return make(map[string]interface{})
	}

	// return custom fields as map
	return parseToMap(customProps.Attributes)
}

// SetCustomProperties sets custom properties of the given asset
func SetCustomProperties(metadata models.Metadata, customFields map[string]interface{}) (models.Metadata, error) {
	properties, err := appendCustomFields(metadata, customFields)
	if err != nil {
		return metadata, errors.Wrap(err, "failed to append custom fields in metadata")
	}

	switch metadata := metadata.(type) {
	case *assets.Table:
		metadata.Properties = properties
	case *assets.Topic:
		metadata.Properties = properties
	case *assets.Dashboard:
		metadata.Properties = properties
	case *assets.Bucket:
		metadata.Properties = properties
	case *assets.Group:
		metadata.Properties = properties
	case *assets.Job:
		metadata.Properties = properties
	case *assets.User:
		metadata.Properties = properties
	}

	return metadata, nil
}

func appendCustomFields(metadata models.Metadata, customFields map[string]interface{}) (*facets.Properties, error) {
	properties := metadata.GetProperties()
	if properties == nil {
		properties = &facets.Properties{
			Attributes: &structpb.Struct{},
		}
	}

	protoStruct, err := parseMapToProto(customFields)
	if err != nil {
		return properties, errors.Wrap(err, "failed to parse map to proto")
	}
	properties.Attributes = protoStruct

	return properties, err
}

func parseToMap(src *structpb.Struct) map[string]interface{} {
	if src == nil {
		return nil
	}

	return src.AsMap()
}

func parseMapToProto(src map[string]interface{}) (*structpb.Struct, error) {
	if src == nil {
		return nil, nil
	}

	return structpb.NewStruct(src)
}

// TryParseMapToProto parses given map to proto struct
func TryParseMapToProto(src map[string]interface{}) *structpb.Struct {
	res, err := parseMapToProto(src)
	if err != nil {
		panic(err)
	}

	return res
}

package util

import (
	"encoding/json"

	"github.com/berquerant/jsonhttp/internal/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// IsInt returns true if v is an integer.
func IsInt(v float64) bool { return float64(int(v)) == v }

// MergeStringMap merges string maps.
// Overwrites dest by values from other for the same keys.
// Panic if dest is nil.
func MergeStringMap(dest, other map[string]string) {
	for k, v := range other {
		dest[k] = v
	}
}

// MergeMap merges interface maps.
// Overwrites dest by values from other for the same keys.
// Panic if dest is nil.
func MergeMap(dest, other map[string]interface{}) {
	for k, v := range other {
		if _, ok := dest[k]; !ok {
			dest[k] = v
			continue
		}
		vv, ok := v.(map[string]interface{})
		if !ok {
			dest[k] = v
			continue
		}
		dv, ok := dest[k].(map[string]interface{})
		if !ok {
			dest[k] = v
			continue
		}
		MergeMap(dv, vv)
		dest[k] = dv
	}
}

// JSON jsonify value for logging.
func JSON(v interface{}) string {
	if x, ok := v.(proto.Message); ok {
		b, err := protojson.Marshal(x)
		if err != nil {
			return errors.Wrapf(err, errors.Jsonify, "%v", x).Error()
		}
		return string(b)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return errors.Wrapf(err, errors.Jsonify, "%v", v).Error()
	}
	return string(b)
}

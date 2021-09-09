package pb

import (
	"encoding/json"
	"strings"

	"github.com/berquerant/jsonhttp/internal/errors"
)

// BodyBuilder extracts a value from http request body (application/json).
type BodyBuilder interface {
	Build(body []byte) (*Value, error)
}

func NewBodyBuilder(body *Value_Body, vc ValueConverter) BodyBuilder {
	return &bodyBuilder{
		body: body,
		vc:   vc,
	}
}

type bodyBuilder struct {
	body *Value_Body
	vc   ValueConverter
}

func (s *bodyBuilder) Build(body []byte) (*Value, error) {
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, errors.Wrap(err, errors.InvalidArgument, "body builder accept only json")
	}
	x, err := s.dig(0, v)
	if err != nil {
		return nil, errors.Wrap(err, errors.InvalidArgument, "body builder cannot build body")
	}
	return x, nil
}

func (s *bodyBuilder) dig(keyIndex int, v map[string]interface{}) (*Value, error) {
	keys := s.body.GetKeys()
	if keyIndex < 0 || keyIndex >= len(keys) {
		return nil, errors.New(errors.OutOfRange, "")
	}
	key := keys[keyIndex]
	value, ok := v[key]
	if !ok {
		return nil, errors.Newf(errors.NotFound, "%s not in body", strings.Join(s.body.GetKeys(), "."))
	}
	if keyIndex == len(keys)-1 {
		v, err := s.vc.Convert(value)
		if err != nil {
			return nil, errors.Wrap(err, errors.InvalidArgument, "")
		}
		return v, nil
	}
	if m, ok := value.(map[string]interface{}); ok {
		return s.dig(keyIndex+1, m)
	}
	return nil, errors.Newf(errors.InvalidArgument, "cannot dig %T %v", value, value)
}

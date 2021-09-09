package pb

import (
	"github.com/berquerant/jsonhttp/internal/errors"
	"github.com/berquerant/jsonhttp/internal/util"
)

// CastBuilder does typecast from Value into the specific Value.
type CastBuilder interface {
	Build(r TemplateSource) (*Value, error)
}

func NewCastBuilder(cast *Value_Cast, valueCoercer ValueCoercer, templateValueBuilder TemplateValueBuilder) CastBuilder {
	return &castBuilder{
		cast:                 cast,
		valueCoercer:         valueCoercer,
		templateValueBuilder: templateValueBuilder,
	}
}

type castBuilder struct {
	cast                 *Value_Cast
	valueCoercer         ValueCoercer
	templateValueBuilder TemplateValueBuilder
}

func (s *castBuilder) Build(r TemplateSource) (*Value, error) {
	b, err := s.templateValueBuilder.Build(s.cast.GetValue(), r)
	if err != nil {
		return nil, errors.Wrapf(err, errors.InvalidValue, "cannot build cast %s %s", s.cast.GetType(), util.JSON(s.cast.GetValue()))
	}
	v, err := s.coerce(b)
	if err != nil {
		return nil, errors.Wrapf(err, errors.InvalidValue, "cannot build cast %s %s", s.cast.GetType(), util.JSON(s.cast.GetValue()))
	}
	return v, nil
}

func (s *castBuilder) coerce(v *Value) (*Value, error) {
	switch s.cast.GetType() {
	case Value_Cast_BOOL:
		return s.valueCoercer.B(v)
	case Value_Cast_NUMBER:
		return s.valueCoercer.N(v)
	case Value_Cast_STRING:
		return s.valueCoercer.S(v)
	}
	return nil, errors.New(errors.UnknownError, "cannot coerce")
}

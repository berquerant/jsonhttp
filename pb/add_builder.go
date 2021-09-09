package pb

import (
	"strings"

	"github.com/berquerant/jsonhttp/internal/errors"
	"github.com/berquerant/jsonhttp/internal/util"
)

// AddBuilder adds up all its arguments.
type AddBuilder interface {
	Build(r TemplateSource) (*Value, error)
}

func NewAddBuilder(add *Value_Add, valueCaster ValueCaster, templateValueBuilder TemplateValueBuilder) AddBuilder {
	return &addBuilder{
		add:                  add,
		valueCaster:          valueCaster,
		templateValueBuilder: templateValueBuilder,
	}
}

type addBuilder struct {
	add                  *Value_Add
	valueCaster          ValueCaster
	templateValueBuilder TemplateValueBuilder
}

func (s *addBuilder) Build(r TemplateSource) (*Value, error) {
	if len(s.add.GetValues()) == 0 {
		return NewNull(), nil
	}
	values := make([]*Value, len(s.add.GetValues()))
	for i, v := range s.add.GetValues() {
		x, err := s.templateValueBuilder.Build(v, r)
		if err != nil {
			return nil, errors.Wrapf(err, errors.InvalidValue, "cannot build add %d", i)
		}
		values[i] = x
	}
	switch s.add.GetType() {
	case Value_Add_STRING:
		ss := make([]string, len(values))
		for i, x := range values {
			v, err := s.valueCaster.String(x)
			if err != nil {
				return nil, errors.Wrapf(err, errors.TypeCast, "cannot build add %d %s", i, util.JSON(x))
			}
			ss[i] = v
		}
		return NewS(strings.Join(ss, "")), nil
	case Value_Add_NUMBER:
		var sum float64
		for i, x := range values {
			v, err := s.valueCaster.Float(x)
			if err != nil {
				return nil, errors.Wrapf(err, errors.TypeCast, "cannot build add %d %s", i, util.JSON(x))
			}
			sum += v
		}
		return NewN(sum), nil
	}
	return nil, errors.New(errors.UnknownError, "add builder")
}

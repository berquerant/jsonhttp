package pb

import (
	"fmt"
	reflect "reflect"
	"strconv"

	"github.com/berquerant/jsonhttp/internal/errors"
	"github.com/berquerant/jsonhttp/internal/util"
)

// ValueCoercer does typecast on Value.
type ValueCoercer interface {
	B(value *Value) (*Value, error)
	N(value *Value) (*Value, error)
	S(value *Value) (*Value, error)
}

func NewValueCoercer(valueCaster ValueCaster) ValueCoercer {
	return &valueCoercer{
		valueCaster: valueCaster,
	}
}

type valueCoercer struct {
	valueCaster ValueCaster
}

func (s *valueCoercer) B(value *Value) (*Value, error) {
	x, err := s.valueCaster.Bool(value)
	if err != nil {
		return nil, errors.Wrapf(err, errors.TypeCoerce, "cannot type coerce B %s", util.JSON(value))
	}
	return NewB(x), nil
}

func (s *valueCoercer) N(value *Value) (*Value, error) {
	x, err := s.valueCaster.Float(value)
	if err != nil {
		return nil, errors.Wrapf(err, errors.TypeCoerce, "cannot type coerce N %s", util.JSON(value))
	}
	return NewN(x), nil
}

func (s *valueCoercer) S(value *Value) (*Value, error) {
	x, err := s.valueCaster.String(value)
	if err != nil {
		return nil, errors.Wrapf(err, errors.TypeCoerce, "cannot type coerce S %s", util.JSON(value))
	}
	return NewS(x), nil
}

// ValueCaster does typecast.
type ValueCaster interface {
	Bool(value *Value) (bool, error)
	Int(value *Value) (int, error)
	Float(value *Value) (float64, error)
	String(value *Value) (string, error)
}

func NewValueCaster() ValueCaster { return &valueCaster{} }

type valueCaster struct {
}

func (*valueCaster) String(value *Value) (string, error) {
	switch value.GetValue().(type) {
	case *Value_Null:
		return "null", nil
	case *Value_B:
		return fmt.Sprint(value.GetB()), nil
	case *Value_N:
		return fmt.Sprint(value.GetN()), nil
	case *Value_S:
		return value.GetS(), nil
	}
	return "", errors.Newf(errors.TypeCast, "cannot type cast string %s", util.JSON(value))
}

func (*valueCaster) Float(value *Value) (float64, error) {
	switch value.GetValue().(type) {
	case *Value_Null:
		return 0, nil
	case *Value_B:
		if value.GetB() {
			return 1, nil
		}
		return 0, nil
	case *Value_N:
		return value.GetN(), nil
	case *Value_S:
		x, err := strconv.ParseFloat(value.GetS(), 64)
		if err != nil {
			return 0, errors.Wrapf(err, errors.TypeCast, "cannot type cast int %s", util.JSON(value))
		}
		return x, nil
	}
	return 0, errors.Newf(errors.TypeCast, "cannot type cast float %s", util.JSON(value))
}

func (s *valueCaster) Int(value *Value) (int, error) {
	x, err := s.Float(value)
	if err != nil {
		return 0, errors.Wrapf(err, errors.TypeCast, "cannot type cast int %s", util.JSON(value))
	}
	return int(x), nil
}

func (*valueCaster) Bool(value *Value) (bool, error) {
	switch value.GetValue().(type) {
	case *Value_Null:
		return false, nil
	case *Value_B:
		return value.GetB(), nil
	case *Value_N:
		return value.GetN() != 0, nil
	case *Value_S:
		return value.GetS() != "", nil
	}
	return false, errors.Newf(errors.TypeCast, "cannot type cast bool %s", util.JSON(value))
}

// ValueInverter translates pb value for json into interface.
type ValueInverter interface {
	Invert(value *Value) (interface{}, error)
}

func NewValueInverter() ValueInverter { return &valueInverter{} }

type valueInverter struct {
}

func (s *valueInverter) Invert(value *Value) (interface{}, error) {
	switch value.GetValue().(type) {
	case *Value_Null:
		return nil, nil
	case *Value_B:
		return value.GetB(), nil
	case *Value_N:
		if util.IsInt(value.GetN()) {
			return int(value.GetN()), nil
		}
		return value.GetN(), nil
	case *Value_S:
		return value.GetS(), nil
	case *Value_L:
		p := []interface{}{}
		for i, x := range value.GetL().GetValues() {
			v, err := s.Invert(x)
			if err != nil {
				return nil, errors.Wrapf(err, errors.InvalidValue, "%v at %d index in list", x, i)
			}
			p = append(p, v)
		}
		return p, nil
	case *Value_M:
		p := map[string]interface{}{}
		for k, x := range value.GetM().GetValues() {
			v, err := s.Invert(x)
			if err != nil {
				return nil, errors.Wrapf(err, errors.InvalidValue, "%v at key %s in map", x, k)
			}
			p[k] = v
		}
		return p, nil
	}
	return nil, errors.Newf(errors.InvalidArgument, "cannot invert %v", util.JSON(value))
}

// ValueConverter translates value for json into pb.
type ValueConverter interface {
	Convert(value interface{}) (*Value, error)
}

func NewValueConverter() ValueConverter { return &valueConverter{} }

type valueConverter struct {
}

func (s *valueConverter) Convert(value interface{}) (*Value, error) {
	switch v := value.(type) {
	case nil:
		return NewNull(), nil
	case bool:
		return NewB(v), nil
	case float64:
		return NewN(v), nil
	case float32:
		return NewN(float64(v)), nil
	case int:
		return NewN(float64(v)), nil
	case string:
		return NewS(v), nil
	}

	var (
		t = reflect.TypeOf(value)
		v = reflect.ValueOf(value)
	)
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		p := []*Value{}
		for i := 0; i < v.Len(); i++ {
			x, err := s.Convert(v.Index(i).Interface())
			if err != nil {
				return nil, errors.Wrapf(err, errors.InvalidValue, "%v at index %d in list", v.Index(i).Interface(), i)
			}
			p = append(p, x)
		}
		return NewL(p), nil
	case reflect.Map:
		p := map[string]*Value{}
		it := v.MapRange()
		for it.Next() {
			k, ok := it.Key().Interface().(string)
			if !ok {
				continue
			}
			x, err := s.Convert(it.Value().Interface())
			if err != nil {
				return nil, errors.Wrapf(err, errors.InvalidValue, "%v key %s in map", it.Value().Interface(), k)
			}
			p[k] = x
		}
		return NewM(p), nil
	}
	return nil, errors.Newf(errors.InvalidArgument, "cannot convert %T %v", v, v)
}

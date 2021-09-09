package pb

import (
	"math/rand"
	"time"

	"github.com/berquerant/jsonhttp/internal/errors"
	"github.com/google/uuid"
)

// UtilBuilder yields some values based on something other than http request.
type UtilBuilder interface {
	Build() (*Value, error)
}

func NewUtilBuilder(util *Value_Util) UtilBuilder {
	return &utilBuilder{
		util: util,
	}
}

type utilBuilder struct {
	util *Value_Util
}

func (s *utilBuilder) Build() (*Value, error) {
	switch s.util.GetValue().(type) {
	case *Value_Util_Now_:
		return s.buildNow()
	case *Value_Util_Random_:
		return s.buildRandom()
	}
	return nil, errors.New(errors.InvalidSettings, "util builder")
}

func (s *utilBuilder) buildNow() (*Value, error) {
	n := s.util.GetNow()
	switch n.GetType() {
	case Value_Util_Now_TIMESTAMP:
		return NewN(float64(time.Now().Unix())), nil
	}
	return nil, errors.New(errors.UnknownError, "util builder")
}

func (s *utilBuilder) buildRandom() (*Value, error) {
	r := s.util.GetRandom()
	switch r.GetValue().(type) {
	case *Value_Util_Random_Dice_:
		d := r.GetDice()
		min := d.GetMin()
		max := d.GetMax()
		if max <= min {
			return nil, errors.Newf(errors.InvalidSettings, "max must be greater than min %d %d", max, min)
		}
		return NewN(float64(min + rand.Int31()%(max-min))), nil
	case *Value_Util_Random_Type_:
		switch r.GetType() {
		case Value_Util_Random_UUID:
			return NewS(uuid.NewString()), nil
		case Value_Util_Random_STDU:
			return NewN(rand.Float64()), nil
		}
	}
	return nil, errors.New(errors.UnknownError, "util builder")
}

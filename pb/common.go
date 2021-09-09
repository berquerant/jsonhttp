package pb

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func newValue(x isValue_Value) *Value {
	return &Value{
		Value: x,
	}
}

func NewNull() *Value {
	return newValue(&Value_Null{
		Null: structpb.NullValue(0),
	})
}

func NewB(b bool) *Value {
	return newValue(&Value_B{
		B: b,
	})
}

func NewN(f float64) *Value {
	return newValue(&Value_N{
		N: f,
	})
}

func NewS(s string) *Value {
	return newValue(&Value_S{
		S: s,
	})
}

func NewL(l []*Value) *Value {
	return newValue(&Value_L{
		L: &Value_List{
			Values: l,
		},
	})
}

func NewM(m map[string]*Value) *Value {
	return newValue(&Value_M{
		M: &Value_Map{
			Values: m,
		},
	})
}

func (s *Value) Clone() *Value { return proto.Clone(s).(*Value) }

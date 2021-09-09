package pb_test

import (
	"fmt"
	"testing"

	"github.com/berquerant/jsonhttp/pb"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

type mockValueCaster struct {
	value interface{}
	err   error
}

func (s *mockValueCaster) Bool(_ *pb.Value) (bool, error) {
	if s.err != nil {
		return false, s.err
	}
	return s.value.(bool), nil
}

func (s *mockValueCaster) Int(_ *pb.Value) (int, error) {
	if s.err != nil {
		return 0, s.err
	}
	return s.value.(int), nil
}

func (s *mockValueCaster) Float(_ *pb.Value) (float64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return s.value.(float64), nil
}

func (s *mockValueCaster) String(_ *pb.Value) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.value.(string), nil
}

func TestValueCoercer(t *testing.T) {
	valueCastErr := fmt.Errorf("value cast error")
	t.Run("B", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			value interface{}
			err   error
		}{
			{
				title: "false",
				value: false,
			},
			{
				title: "true",
				value: true,
			},
			{
				title: "err",
				err:   valueCastErr,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				got, err := pb.NewValueCoercer(&mockValueCaster{
					value: tc.value,
					err:   tc.err,
				}).B(nil)
				if tc.err != nil {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.value, got.GetB())
			})
		}
	})

	t.Run("N", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			value interface{}
			err   error
		}{
			{
				title: "number",
				value: 1.2,
			},
			{
				title: "err",
				err:   valueCastErr,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				got, err := pb.NewValueCoercer(&mockValueCaster{
					value: tc.value,
					err:   tc.err,
				}).N(nil)
				if tc.err != nil {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.value, got.GetN())
			})
		}
	})

	t.Run("S", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			value interface{}
			err   error
		}{
			{
				title: "string",
				value: "light",
			},
			{
				title: "err",
				err:   valueCastErr,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				got, err := pb.NewValueCoercer(&mockValueCaster{
					value: tc.value,
					err:   tc.err,
				}).S(nil)
				if tc.err != nil {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.value, got.GetS())
			})
		}
	})
}

func TestValueCaster(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			value *pb.Value
			want  interface{}
			isErr bool
		}{
			{
				title: "null",
				value: pb.NewNull(),
				want:  "null",
			},
			{
				title: "true",
				value: pb.NewB(true),
				want:  "true",
			},
			{
				title: "false",
				value: pb.NewB(false),
				want:  "false",
			},
			{
				title: "zero",
				value: pb.NewN(0),
				want:  "0",
			},
			{
				title: "not zero",
				value: pb.NewN(1),
				want:  "1",
			},
			{
				title: "empty string",
				value: pb.NewS(""),
				want:  "",
			},
			{
				title: "string",
				value: pb.NewS("unite"),
				want:  "unite",
			},
			{
				title: "list",
				value: pb.NewL([]*pb.Value{}),
				isErr: true,
			},
			{
				title: "map",
				value: pb.NewM(map[string]*pb.Value{}),
				isErr: true,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				got, err := pb.NewValueCaster().String(tc.value)
				if tc.isErr {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("Float", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			value *pb.Value
			want  interface{}
			isErr bool
		}{
			{
				title: "null",
				value: pb.NewNull(),
				want:  0.0,
			},
			{
				title: "true",
				value: pb.NewB(true),
				want:  1.0,
			},
			{
				title: "false",
				value: pb.NewB(false),
				want:  0.0,
			},
			{
				title: "zero",
				value: pb.NewN(0),
				want:  0.0,
			},
			{
				title: "not zero",
				value: pb.NewN(1),
				want:  1.0,
			},
			{
				title: "empty string",
				value: pb.NewS(""),
				isErr: true,
			},
			{
				title: "not number string",
				value: pb.NewS("unite"),
				isErr: true,
			},
			{
				title: "int string",
				value: pb.NewS("120"),
				want:  120.0,
			},
			{
				title: "float string",
				value: pb.NewS("120.2"),
				want:  120.2,
			},
			{
				title: "list",
				value: pb.NewL([]*pb.Value{}),
				isErr: true,
			},
			{
				title: "map",
				value: pb.NewM(map[string]*pb.Value{}),
				isErr: true,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				got, err := pb.NewValueCaster().Float(tc.value)
				if tc.isErr {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("Int", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			value *pb.Value
			want  interface{}
			isErr bool
		}{
			{
				title: "null",
				value: pb.NewNull(),
				want:  0,
			},
			{
				title: "true",
				value: pb.NewB(true),
				want:  1,
			},
			{
				title: "false",
				value: pb.NewB(false),
				want:  0,
			},
			{
				title: "zero",
				value: pb.NewN(0),
				want:  0,
			},
			{
				title: "not zero",
				value: pb.NewN(1),
				want:  1,
			},
			{
				title: "empty string",
				value: pb.NewS(""),
				isErr: true,
			},
			{
				title: "not number string",
				value: pb.NewS("unite"),
				isErr: true,
			},
			{
				title: "int string",
				value: pb.NewS("120"),
				want:  120,
			},
			{
				title: "float string",
				value: pb.NewS("120.2"),
				want:  120,
			},
			{
				title: "list",
				value: pb.NewL([]*pb.Value{}),
				isErr: true,
			},
			{
				title: "map",
				value: pb.NewM(map[string]*pb.Value{}),
				isErr: true,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				got, err := pb.NewValueCaster().Int(tc.value)
				if tc.isErr {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("Bool", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			value *pb.Value
			want  interface{}
			isErr bool
		}{
			{
				title: "null",
				value: pb.NewNull(),
				want:  false,
			},
			{
				title: "true",
				value: pb.NewB(true),
				want:  true,
			},
			{
				title: "false",
				value: pb.NewB(false),
				want:  false,
			},
			{
				title: "zero",
				value: pb.NewN(0),
				want:  false,
			},
			{
				title: "not zero",
				value: pb.NewN(1),
				want:  true,
			},
			{
				title: "empty string",
				value: pb.NewS(""),
				want:  false,
			},
			{
				title: "not empty string",
				value: pb.NewS("unite"),
				want:  true,
			},
			{
				title: "list",
				value: pb.NewL([]*pb.Value{}),
				isErr: true,
			},
			{
				title: "map",
				value: pb.NewM(map[string]*pb.Value{}),
				isErr: true,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				got, err := pb.NewValueCaster().Bool(tc.value)
				if tc.isErr {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})
}

func TestValueInverter(t *testing.T) {
	t.Run("Invert", func(t *testing.T) {
		t.Run("null", func(t *testing.T) {
			got, err := pb.NewValueInverter().Invert(pb.NewNull())
			assert.Nil(t, err)
			assert.Nil(t, got)
		})
		t.Run("bool", func(t *testing.T) {
			for _, tc := range []*struct {
				title string
				value *pb.Value
				want  bool
			}{
				{
					title: "true",
					value: pb.NewB(true),
					want:  true,
				},
				{
					title: "false",
					value: pb.NewB(false),
					want:  false,
				},
			} {
				t.Run(tc.title, func(t *testing.T) {
					got, err := pb.NewValueInverter().Invert(tc.value)
					assert.Nil(t, err)
					assert.Equal(t, tc.want, got)
				})
			}
		})
		t.Run("number", func(t *testing.T) {
			t.Run("float", func(t *testing.T) {
				got, err := pb.NewValueInverter().Invert(pb.NewN(1.5))
				assert.Nil(t, err)
				assert.Equal(t, 1.5, got)
			})
			t.Run("int", func(t *testing.T) {
				got, err := pb.NewValueInverter().Invert(pb.NewN(1.0))
				assert.Nil(t, err)
				assert.Equal(t, 1, got)
			})
		})
		t.Run("string", func(t *testing.T) {
			got, err := pb.NewValueInverter().Invert(pb.NewS("inv"))
			assert.Nil(t, err)
			assert.Equal(t, "inv", got)
		})
		t.Run("slice", func(t *testing.T) {
			t.Run("zero", func(t *testing.T) {
				got, err := pb.NewValueInverter().Invert(pb.NewL([]*pb.Value{}))
				assert.Nil(t, err)
				assert.Equal(t, "", cmp.Diff([]interface{}{}, got))
			})
			t.Run("string", func(t *testing.T) {
				got, err := pb.NewValueInverter().Invert(pb.NewL([]*pb.Value{
					pb.NewS("element"),
				}))
				assert.Nil(t, err)
				assert.Equal(t, "", cmp.Diff([]interface{}{"element"}, got))
			})
			t.Run("elements", func(t *testing.T) {
				got, err := pb.NewValueInverter().Invert(pb.NewL([]*pb.Value{
					pb.NewS("element"),
					pb.NewB(true),
				}))
				assert.Nil(t, err)
				assert.Equal(t, "", cmp.Diff([]interface{}{"element", true}, got))
			})
			t.Run("nest", func(t *testing.T) {
				got, err := pb.NewValueInverter().Invert(pb.NewL([]*pb.Value{
					pb.NewL([]*pb.Value{
						pb.NewS("element"),
						pb.NewB(true),
					}),
				}))
				assert.Nil(t, err)
				assert.Equal(t, "", cmp.Diff([]interface{}{
					[]interface{}{"element", true},
				}, got))
			})
		})
		t.Run("map", func(t *testing.T) {
			t.Run("zero", func(t *testing.T) {
				got, err := pb.NewValueInverter().Invert(pb.NewM(map[string]*pb.Value{}))
				assert.Nil(t, err)
				assert.Equal(t, "", cmp.Diff(map[string]interface{}{}, got))
			})
			t.Run("string", func(t *testing.T) {
				got, err := pb.NewValueInverter().Invert(pb.NewM(map[string]*pb.Value{
					"str": pb.NewS("val"),
				}))
				assert.Nil(t, err)
				assert.Equal(t, "", cmp.Diff(map[string]interface{}{
					"str": "val",
				}, got))
			})
			t.Run("assort", func(t *testing.T) {
				got, err := pb.NewValueInverter().Invert(pb.NewM(map[string]*pb.Value{
					"str": pb.NewS("val"),
					"m": pb.NewM(map[string]*pb.Value{
						"i": pb.NewN(1),
					}),
				}))
				assert.Nil(t, err)
				assert.Equal(t, "", cmp.Diff(map[string]interface{}{
					"str": "val",
					"m": map[string]interface{}{
						"i": 1,
					},
				}, got))
			})
		})
	})
}

func TestValueConverter(t *testing.T) {
	const tolerance = 0.0001
	t.Run("Convert", func(t *testing.T) {
		t.Run("null", func(t *testing.T) {
			got, err := pb.NewValueConverter().Convert(nil)
			assert.Nil(t, err)
			assert.False(t, got.GetB())
			assert.InDelta(t, 0, got.GetN(), tolerance)
			assert.Equal(t, "", got.GetS())
			assert.Nil(t, got.GetL())
			assert.Nil(t, got.GetM())
			assert.NotNil(t, got.GetNull())
		})

		t.Run("bool", func(t *testing.T) {
			for _, tc := range []*struct {
				title string
				value interface{}
				want  bool
			}{
				{
					title: "true",
					value: true,
					want:  true,
				},
				{
					title: "false",
					value: false,
					want:  false,
				},
			} {
				t.Run(tc.title, func(t *testing.T) {
					got, err := pb.NewValueConverter().Convert(tc.value)
					assert.Nil(t, err)
					assert.Equal(t, tc.want, got.GetB())
				})
			}
		})

		t.Run("int or float", func(t *testing.T) {
			for _, tc := range []*struct {
				title string
				value interface{}
				want  float64
			}{
				{
					title: "zero",
					value: 0,
					want:  0,
				},
				{
					title: "int",
					value: 10,
					want:  10.0,
				},
				{
					title: "float",
					value: 10.1,
					want:  10.1,
				},
				{
					title: "negative float",
					value: -27.1,
					want:  -27.1,
				},
			} {
				t.Run(tc.title, func(t *testing.T) {
					got, err := pb.NewValueConverter().Convert(tc.value)
					assert.Nil(t, err)
					assert.InDelta(t, tc.want, got.GetN(), tolerance)
				})
			}
		})

		t.Run("string", func(t *testing.T) {
			for _, tc := range []*struct {
				title string
				value interface{}
				want  string
			}{
				{
					title: "zero",
					value: "",
					want:  "",
				},
				{
					title: "string",
					value: "command",
					want:  "command",
				},
			} {
				got, err := pb.NewValueConverter().Convert(tc.value)
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got.GetS())
			}
		})

		t.Run("slice", func(t *testing.T) {
			t.Run("zero", func(t *testing.T) {
				got, err := pb.NewValueConverter().Convert([]int{})
				assert.Nil(t, err)
				assert.Equal(t, 0, len(got.GetL().GetValues()))
			})
			t.Run("int", func(t *testing.T) {
				got, err := pb.NewValueConverter().Convert([]int{1})
				assert.Nil(t, err)
				v := got.GetL().GetValues()
				assert.Equal(t, 1, len(v))
				assert.InDelta(t, 1, v[0].GetN(), tolerance)
			})
			t.Run("strings", func(t *testing.T) {
				got, err := pb.NewValueConverter().Convert([]string{"wonder", "world"})
				assert.Nil(t, err)
				v := got.GetL().GetValues()
				assert.Equal(t, 2, len(v))
				assert.Equal(t, "wonder", v[0].GetS())
				assert.Equal(t, "world", v[1].GetS())
			})
			t.Run("list of bools", func(t *testing.T) {
				got, err := pb.NewValueConverter().Convert([][]bool{
					{true},
					{false, true},
				})
				assert.Nil(t, err)
				v := got.GetL().GetValues()
				assert.Equal(t, 2, len(v))
				{
					v := v[0].GetL().GetValues()
					assert.Equal(t, 1, len(v))
					assert.True(t, v[0].GetB())
				}
				{
					v := v[1].GetL().GetValues()
					assert.Equal(t, 2, len(v))
					assert.False(t, v[0].GetB())
					assert.True(t, v[1].GetB())
				}
			})
		})

		t.Run("map", func(t *testing.T) {
			t.Run("zero", func(t *testing.T) {
				got, err := pb.NewValueConverter().Convert(map[string]interface{}{})
				assert.Nil(t, err)
				assert.Equal(t, 0, len(got.GetM().GetValues()))
			})
			t.Run("string", func(t *testing.T) {
				got, err := pb.NewValueConverter().Convert(map[string]interface{}{
					"key": "value",
				})
				assert.Nil(t, err)
				v := got.GetM().GetValues()
				assert.Equal(t, 1, len(v))
				assert.Contains(t, v, "key")
				assert.Equal(t, "value", v["key"].GetS())
			})
			t.Run("assort", func(t *testing.T) {
				got, err := pb.NewValueConverter().Convert(map[string]interface{}{
					"l": []string{
						"a",
						"b",
					},
					"m": map[string]interface{}{
						"s": "innerstring",
					},
				})
				assert.Nil(t, err)
				v := got.GetM().GetValues()
				assert.Equal(t, 2, len(v))
				assert.Contains(t, v, "l")
				assert.Contains(t, v, "m")
				{
					v := v["l"].GetL().GetValues()
					assert.Equal(t, 2, len(v))
					assert.Equal(t, "a", v[0].GetS())
					assert.Equal(t, "b", v[1].GetS())
				}
				{
					v := v["m"].GetM().GetValues()
					assert.Equal(t, 1, len(v))
					assert.Contains(t, v, "s")
					assert.Equal(t, "innerstring", v["s"].GetS())
				}
			})
		})
	})
}

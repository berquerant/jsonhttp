package pb_test

import (
	"fmt"
	"testing"

	"github.com/berquerant/jsonhttp/pb"
	"github.com/stretchr/testify/assert"
)

type mockTemplateValueBuilder struct {
	err error
}

func (s *mockTemplateValueBuilder) Build(value *pb.Value, _ pb.TemplateSource) (*pb.Value, error) {
	if s.err != nil {
		return nil, s.err
	}
	return value, nil
}

type mockRepeatValueCaster struct {
	i      int
	values []interface{}
	err    error
}

func (s *mockRepeatValueCaster) Bool(_ *pb.Value) (bool, error) {
	defer func() {
		s.i++
	}()
	if s.err != nil {
		return false, s.err
	}
	return s.values[s.i].(bool), nil
}

func (s *mockRepeatValueCaster) Int(_ *pb.Value) (int, error) {
	defer func() {
		s.i++
	}()
	if s.err != nil {
		return 0, s.err
	}
	return s.values[s.i].(int), nil
}

func (s *mockRepeatValueCaster) Float(_ *pb.Value) (float64, error) {
	defer func() {
		s.i++
	}()
	if s.err != nil {
		return 0, s.err
	}
	return s.values[s.i].(float64), nil
}

func (s *mockRepeatValueCaster) String(_ *pb.Value) (string, error) {
	defer func() {
		s.i++
	}()
	if s.err != nil {
		return "", s.err
	}
	return s.values[s.i].(string), nil
}

func TestAddBuilder(t *testing.T) {
	templateValueBuildErr := fmt.Errorf("template value build error")
	castErr := fmt.Errorf("cast error")
	newValues := func(n int) []*pb.Value {
		v := make([]*pb.Value, n)
		for i := 0; i < n; i++ {
			v[i] = pb.NewNull()
		}
		return v
	}

	t.Run("Build", func(t *testing.T) {
		t.Run("Null", func(t *testing.T) {
			got, err := pb.NewAddBuilder(&pb.Value_Add{}, nil, nil).Build(nil)
			assert.Nil(t, err)
			_, ok := got.GetValue().(*pb.Value_Null)
			assert.True(t, ok)
		})

		t.Run("Template error", func(t *testing.T) {
			t.Run("String", func(t *testing.T) {
				_, err := pb.NewAddBuilder(&pb.Value_Add{
					Type:   pb.Value_Add_STRING,
					Values: newValues(1),
				}, nil, &mockTemplateValueBuilder{
					err: templateValueBuildErr,
				}).Build(nil)
				assert.NotNil(t, err)
			})
			t.Run("Number", func(t *testing.T) {
				_, err := pb.NewAddBuilder(&pb.Value_Add{
					Type:   pb.Value_Add_NUMBER,
					Values: newValues(1),
				}, nil, &mockTemplateValueBuilder{
					err: templateValueBuildErr,
				}).Build(nil)
				assert.NotNil(t, err)
			})
		})

		t.Run("Cast error", func(t *testing.T) {
			_, err := pb.NewAddBuilder(&pb.Value_Add{
				Type:   pb.Value_Add_STRING,
				Values: newValues(1),
			}, &mockRepeatValueCaster{
				err: castErr,
			}, &mockTemplateValueBuilder{}).Build(nil)
			assert.NotNil(t, err)
		})

		t.Run("String", func(t *testing.T) {
			for _, tc := range []*struct {
				title  string
				add    *pb.Value_Add
				values []interface{}
				want   string
			}{
				{
					title: "a string",
					add: &pb.Value_Add{
						Type:   pb.Value_Add_STRING,
						Values: newValues(1),
					},
					values: []interface{}{
						"pyro",
					},
					want: "pyro",
				},
				{
					title: "two strings",
					add: &pb.Value_Add{
						Type:   pb.Value_Add_STRING,
						Values: newValues(2),
					},
					values: []interface{}{
						"pyro", "mancer",
					},
					want: "pyromancer",
				},
			} {
				t.Run(tc.title, func(t *testing.T) {
					got, err := pb.NewAddBuilder(tc.add, &mockRepeatValueCaster{
						values: tc.values,
					}, &mockTemplateValueBuilder{}).Build(nil)
					assert.Nil(t, err)
					assert.Equal(t, got.GetS(), tc.want)
				})
			}
		})

		t.Run("Number", func(t *testing.T) {
			for _, tc := range []*struct {
				title  string
				add    *pb.Value_Add
				values []interface{}
				want   float64
			}{
				{
					title: "a number",
					add: &pb.Value_Add{
						Type:   pb.Value_Add_NUMBER,
						Values: newValues(1),
					},
					values: []interface{}{
						1.0,
					},
					want: 1.0,
				},
				{
					title: "two numbers",
					add: &pb.Value_Add{
						Type:   pb.Value_Add_NUMBER,
						Values: newValues(2),
					},
					values: []interface{}{
						1.2,
						3.5,
					},
					want: 4.7,
				},
			} {
				t.Run(tc.title, func(t *testing.T) {
					got, err := pb.NewAddBuilder(tc.add, &mockRepeatValueCaster{
						values: tc.values,
					}, &mockTemplateValueBuilder{}).Build(nil)
					assert.Nil(t, err)
					assert.Equal(t, got.GetN(), tc.want)
				})
			}
		})
	})
}

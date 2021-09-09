package pb_test

import (
	"testing"

	"github.com/berquerant/jsonhttp/pb"
	"github.com/stretchr/testify/assert"
)

type mockValueConverter struct {
	v interface{}
}

func (s *mockValueConverter) Convert(value interface{}) (*pb.Value, error) {
	s.v = value
	return nil, nil
}

func TestBodyBuilder(t *testing.T) {
	t.Run("Build", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			vb    *pb.Value_Body
			b     []byte
			want  interface{}
			isErr bool
		}{
			{
				title: "invalid json",
				vb: &pb.Value_Body{
					Keys: []string{"t"},
				},
				b:     []byte(`vvv`),
				isErr: true,
			},
			{
				title: "no keys",
				vb:    &pb.Value_Body{},
				b:     []byte(`{"t":true}`),
				isErr: true,
			},
			{
				title: "got a value",
				vb: &pb.Value_Body{
					Keys: []string{"t"},
				},
				b:    []byte(`{"t":true}`),
				want: true,
			},
			{
				title: "no hits",
				vb: &pb.Value_Body{
					Keys: []string{"t", "x"},
				},
				b:     []byte(`{"t":true}`),
				isErr: true,
			},
			{
				title: "deep",
				vb: &pb.Value_Body{
					Keys: []string{"t", "x"},
				},
				b:    []byte(`{"t":{"x":"deep"},"x":true}`),
				want: "deep",
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				vc := &mockValueConverter{}
				_, err := pb.NewBodyBuilder(tc.vb, vc).Build(tc.b)
				if tc.isErr {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.want, vc.v)
			})
		}
	})
}

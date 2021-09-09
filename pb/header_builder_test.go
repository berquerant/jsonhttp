package pb_test

import (
	"net/http"
	"testing"

	"github.com/berquerant/jsonhttp/pb"
	"github.com/stretchr/testify/assert"
)

func TestHeaderBuilder(t *testing.T) {
	t.Run("Build", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			vh    *pb.Value_Header
			h     *http.Header
			want  *pb.Value
		}{
			{
				title: "no headers",
				vh: &pb.Value_Header{
					Key: "Host",
				},
			},
			{
				title: "cannot hit",
				vh: &pb.Value_Header{
					Key: "Host",
				},
				h: &http.Header{
					"Accept": []string{
						"application/json",
					},
				},
			},
			{
				title: "hit",
				vh: &pb.Value_Header{
					Key: "Accept",
				},
				h: &http.Header{
					"Accept": []string{
						"application/json",
					},
				},
				want: &pb.Value{
					Value: &pb.Value_S{
						S: "application/json",
					},
				},
			},
			{
				title: "hit head",
				vh: &pb.Value_Header{
					Key: "Accept",
				},
				h: &http.Header{
					"Accept": []string{
						"application/json",
						"text/plain",
					},
				},
				want: &pb.Value{
					Value: &pb.Value_S{
						S: "application/json",
					},
				},
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				b := pb.NewHeaderBuilder(tc.vh)
				got, err := b.Build(tc.h)
				if tc.want == nil {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.want.GetS(), got.GetS())
			})
		}
	})
}

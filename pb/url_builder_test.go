package pb_test

import (
	"net/url"
	"testing"

	"github.com/berquerant/jsonhttp/pb"
	"github.com/stretchr/testify/assert"
)

func TestURLBuilder(t *testing.T) {
	t.Run("BuildQuery", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			query *pb.Value_Url_Query
			url   string
			want  string
		}{
			{
				title: "no query",
				query: &pb.Value_Url_Query{
					Key: "q",
				},
				url: "https://example.com",
			},
			{
				title: "no hits",
				query: &pb.Value_Url_Query{
					Key: "q",
				},
				url: "https://example.com?x=1",
			},
			{
				title: "hit",
				query: &pb.Value_Url_Query{
					Key: "x",
				},
				url:  "https://example.com?x=1",
				want: "1",
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				u, err := url.Parse(tc.url)
				assert.Nil(t, err)
				got, err := pb.NewURLBuilder(&pb.Value_Url{
					Value: &pb.Value_Url_Query_{
						Query: tc.query,
					},
				}).Build(u)
				if tc.want == "" {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got.GetS())
			})
		}
	})
	t.Run("BuildPath", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			path  *pb.Value_Url_Path
			url   string
			want  string
		}{
			{
				title: "out of range",
				path: &pb.Value_Url_Path{
					Index: 10,
				},
				url: "https://example.com/a",
			},
			{
				title: "no path",
				path: &pb.Value_Url_Path{
					Index: 0,
				},
				url: "https://example.com",
			},
			{
				title: "get first",
				path: &pb.Value_Url_Path{
					Index: 0,
				},
				url:  "https://example.com/a/b/",
				want: "a",
			},
			{
				title: "get first",
				path: &pb.Value_Url_Path{
					Index: 1,
				},
				url:  "https://example.com/a/b/",
				want: "b",
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				u, err := url.Parse(tc.url)
				assert.Nil(t, err)
				got, err := pb.NewURLBuilder(&pb.Value_Url{
					Value: &pb.Value_Url_Path_{
						Path: tc.path,
					},
				}).Build(u)
				if tc.want == "" {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got.GetS())
			})
		}
	})
	const (
		url1 = "https://esto@localhost:6661/a/b?q1=1&q2=2#frag"
	)
	t.Run("BuildPart", func(t *testing.T) {
		for _, tc := range []*struct {
			title string
			part  pb.Value_Url_Part
			url   string
			want  string
		}{
			{
				title: "scheme",
				part:  pb.Value_Url_SCHEME,
				url:   url1,
				want:  "https",
			},
			{
				title: "authority",
				part:  pb.Value_Url_AUTHORITY,
				url:   url1,
				want:  "esto@localhost:6661",
			},
			{
				title: "host",
				part:  pb.Value_Url_HOST,
				url:   url1,
				want:  "localhost",
			},
			{
				title: "port",
				part:  pb.Value_Url_PORT,
				url:   url1,
				want:  "6661",
			},
			{
				title: "path",
				part:  pb.Value_Url_PATH,
				url:   url1,
				want:  "/a/b",
			},
			{
				title: "path",
				part:  pb.Value_Url_QUERY,
				url:   url1,
				want:  "q1=1&q2=2",
			},
			{
				title: "path",
				part:  pb.Value_Url_FRAGMENT,
				url:   url1,
				want:  "frag",
			},
			{
				title: "path",
				part:  pb.Value_Url_ALL,
				url:   url1,
				want:  url1,
			},
		} {
			t.Run(tc.title, func(t *testing.T) {
				u, err := url.Parse(tc.url)
				assert.Nil(t, err)
				got, err := pb.NewURLBuilder(&pb.Value_Url{
					Value: &pb.Value_Url_Part_{
						Part: tc.part,
					},
				}).Build(u)
				assert.Nil(t, err)
				assert.Equal(t, tc.want, got.GetS())
			})
		}
	})
}

package util_test

import (
	"testing"

	"github.com/berquerant/jsonhttp/internal/util"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestIsInt(t *testing.T) {
	t.Run("yes", func(t *testing.T) {
		assert.True(t, util.IsInt(1.0))
	})
	t.Run("no", func(t *testing.T) {
		assert.False(t, util.IsInt(1.5))
	})
}

func TestMergeMap(t *testing.T) {
	for _, tc := range []*struct {
		title             string
		dest, other, want map[string]interface{}
	}{
		{
			title: "zero",
			want:  map[string]interface{}{},
		},
		{
			title: "dest zero",
			other: map[string]interface{}{
				"k1": 1,
			},
			want: map[string]interface{}{
				"k1": 1,
			},
		},
		{
			title: "other zero",
			dest: map[string]interface{}{
				"k1": 1,
			},
			want: map[string]interface{}{
				"k1": 1,
			},
		},
		{
			title: "independent",
			dest: map[string]interface{}{
				"k1": 1,
			},
			other: map[string]interface{}{
				"k2": 2,
			},
			want: map[string]interface{}{
				"k1": 1,
				"k2": 2,
			},
		},
		{
			title: "dependent",
			dest: map[string]interface{}{
				"k1": 1,
				"k3": 3,
			},
			other: map[string]interface{}{
				"k2": 2,
				"k3": "3",
			},
			want: map[string]interface{}{
				"k1": 1,
				"k2": 2,
				"k3": "3",
			},
		},
		{
			title: "dependent but merge map",
			dest: map[string]interface{}{
				"k1": 1,
				"k3": map[string]interface{}{
					"i1": 100,
					"i2": 200,
				},
			},
			other: map[string]interface{}{
				"k2": 2,
				"k3": map[string]interface{}{
					"i2": 1000,
					"i3": 300,
				},
			},
			want: map[string]interface{}{
				"k1": 1,
				"k2": 2,
				"k3": map[string]interface{}{
					"i1": 100,
					"i2": 1000,
					"i3": 300,
				},
			},
		},
		{
			title: "dependent internal key overwrite",
			dest: map[string]interface{}{
				"k3": map[string]interface{}{
					"i1": 100,
					"i2": 200,
				},
			},
			other: map[string]interface{}{
				"k3": 100,
			},
			want: map[string]interface{}{
				"k3": 100,
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got := make(map[string]interface{}, len(tc.dest))
			for k, v := range tc.dest {
				got[k] = v
			}
			util.MergeMap(got, tc.other)
			assert.Equal(t, "", cmp.Diff(got, tc.want))
		})
	}
}

func TestMergeStringMap(t *testing.T) {
	for _, tc := range []*struct {
		title             string
		dest, other, want map[string]string
	}{
		{
			title: "zero",
			want:  map[string]string{},
		},
		{
			title: "dest zero",
			other: map[string]string{
				"k1": "v1",
			},
			want: map[string]string{
				"k1": "v1",
			},
		},
		{
			title: "other zero",
			dest: map[string]string{
				"k1": "v1",
			},
			want: map[string]string{
				"k1": "v1",
			},
		},
		{
			title: "independent",
			dest: map[string]string{
				"k2": "v2",
			},
			other: map[string]string{
				"k1": "v1",
			},
			want: map[string]string{
				"k1": "v1",
				"k2": "v2",
			},
		},
		{
			title: "dependent",
			dest: map[string]string{
				"k2": "v2",
				"k3": "k3",
			},
			other: map[string]string{
				"k1": "v1",
				"k3": "x",
			},
			want: map[string]string{
				"k1": "v1",
				"k2": "v2",
				"k3": "x",
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got := make(map[string]string, len(tc.dest))
			for k, v := range tc.dest {
				got[k] = v
			}
			util.MergeStringMap(got, tc.other)
			assert.Equal(t, "", cmp.Diff(got, tc.want))
		})
	}
}

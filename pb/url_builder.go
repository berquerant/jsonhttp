package pb

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/berquerant/jsonhttp/internal/errors"
	"github.com/berquerant/jsonhttp/internal/util"
)

// URLBuilder extracts a value from URL.
type URLBuilder interface {
	Build(u *url.URL) (*Value, error)
}

func NewURLBuilder(urlValue *Value_Url) URLBuilder {
	return &urlBuilder{
		urlValue: urlValue,
	}
}

type urlBuilder struct {
	urlValue *Value_Url
}

func (s *urlBuilder) Build(u *url.URL) (*Value, error) {
	switch s.urlValue.GetValue().(type) {
	case *Value_Url_Path_:
		return s.buildPath(u)
	case *Value_Url_Query_:
		return s.buildQuery(u)
	case *Value_Url_Part_:
		return s.buildPart(u)
	}
	return nil, errors.Newf(errors.InvalidSettings, "url builder %s", util.JSON(s.urlValue))
}

func (s *urlBuilder) buildQuery(u *url.URL) (*Value, error) {
	var (
		x = s.urlValue.GetQuery()
		q = u.Query()
	)
	if v := q.Get(x.GetKey()); v != "" {
		return NewS(v), nil
	}
	return nil, errors.Newf(errors.NotFound, "%s not in query %s", x.GetKey(), q)
}

func (s *urlBuilder) buildPath(u *url.URL) (*Value, error) {
	if u.Path == "" {
		return nil, errors.New(errors.InvalidArgument, "no path")
	}
	var (
		p = strings.Split(strings.TrimLeft(u.Path, "/"), "/")
		x = s.urlValue.GetPath()
	)
	if x.GetIndex() < 0 || int(x.GetIndex()) >= len(p) {
		return nil, errors.Newf(errors.OutOfRange, "at %d in path %s", x.GetIndex(), u.Path)
	}
	return NewS(p[x.GetIndex()]), nil
}

func (s *urlBuilder) buildPart(u *url.URL) (*Value, error) {
	switch s.urlValue.GetPart() {
	case Value_Url_SCHEME:
		return NewS(u.Scheme), nil
	case Value_Url_AUTHORITY:
		if u.User != nil {
			return NewS(fmt.Sprintf("%s@%s", u.User, u.Host)), nil
		}
		return NewS(u.Host), nil
	case Value_Url_HOST:
		return NewS(u.Hostname()), nil
	case Value_Url_PORT:
		return NewS(u.Port()), nil
	case Value_Url_PATH:
		return NewS(u.Path), nil
	case Value_Url_QUERY:
		return NewS(u.RawQuery), nil
	case Value_Url_FRAGMENT:
		return NewS(u.Fragment), nil
	case Value_Url_ALL:
		return NewS(u.String()), nil
	}
	return nil, errors.Newf(errors.UnknownError, "cannot build url from %s", u)
}

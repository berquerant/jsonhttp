package pb

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/berquerant/jsonhttp/internal/errors"
	"github.com/berquerant/jsonhttp/internal/util"
)

type TemplateSource interface {
	Body() []byte
	URL() *url.URL
	Header() *http.Header
}

type templateSource struct {
	body   []byte
	url    *url.URL
	header *http.Header
}

func NewTemplateSource(url *url.URL, header *http.Header, body []byte) TemplateSource {
	return &templateSource{
		body:   body,
		url:    url,
		header: header,
	}
}

func (s *templateSource) Body() []byte         { return s.body }
func (s *templateSource) URL() *url.URL        { return s.url }
func (s *templateSource) Header() *http.Header { return s.header }

// TemplatesBuilder extracts and builds elements from http request.
type TemplatesBuilder interface {
	Add(t *Template, r TemplateSource) error
	Headers() map[string]string
	Body() map[string]interface{}
}

func NewTemplatesBuilder(templateValueBuilder TemplateValueBuilder, valueInverter ValueInverter) TemplatesBuilder {
	return &templatesBuilder{
		body:                 map[string]interface{}{},
		headers:              map[string]string{},
		templateValueBuilder: templateValueBuilder,
		valueInverter:        valueInverter,
	}
}

type templatesBuilder struct {
	body                 map[string]interface{}
	headers              map[string]string
	templateValueBuilder TemplateValueBuilder
	valueInverter        ValueInverter
}

func (s *templatesBuilder) Body() map[string]interface{} { return s.body }
func (s *templatesBuilder) Headers() map[string]string   { return s.headers }

func (s *templatesBuilder) Add(t *Template, r TemplateSource) error {
	b := newTemplateBuilder(s.templateValueBuilder, s.valueInverter)
	if err := b.Build(t, r); err != nil {
		return errors.Wrap(err, errors.InvalidArgument, "cannot build templates")
	}
	util.MergeStringMap(s.headers, b.headers)
	util.MergeMap(s.body, b.body)
	return nil
}

func newTemplateBuilder(templateValueBuilder TemplateValueBuilder, valueInverter ValueInverter) *templateBuilder {
	return &templateBuilder{
		body:                 map[string]interface{}{},
		headers:              map[string]string{},
		templateValueBuilder: templateValueBuilder,
		valueInverter:        valueInverter,
	}
}

type templateBuilder struct {
	body                 map[string]interface{}
	headers              map[string]string
	templateValueBuilder TemplateValueBuilder
	valueInverter        ValueInverter
}

func (s *templateBuilder) Build(t *Template, r TemplateSource) error {
	v, err := s.templateValueBuilder.Build(t.GetValue(), r)
	if err != nil {
		return errors.Wrap(err, errors.InvalidArgument, "cannot build template")
	}
	if v.GetM() == nil {
		return nil
	}

	switch t.GetType() {
	case Template_HEADER:
		s.headers = s.mapShallow(v)
	case Template_BODY:
		v, err := s.mapDeep(v)
		if err != nil {
			return errors.Wrap(err, errors.InvalidArgument, "cannot build template")
		}
		s.body = v
	}
	return nil
}

func (s *templateBuilder) mapDeep(value *Value) (map[string]interface{}, error) {
	if value.GetM() == nil {
		return nil, errors.Newf(errors.InvalidValue, "not map %s", util.JSON(value))
	}
	v, err := s.valueInverter.Invert(value)
	if err != nil {
		return nil, errors.Wrapf(err, errors.InvalidValue, "cannot invert %s", util.JSON(value))
	}
	if m, ok := v.(map[string]interface{}); ok {
		return m, nil
	}
	return nil, errors.Newf(errors.InvalidValue, "not map %s", util.JSON(value))
}

func (*templateBuilder) mapShallow(value *Value) map[string]string {
	if value.GetM() == nil {
		return nil
	}
	d := map[string]string{}
	for k, v := range value.GetM().GetValues() {
		switch v.GetValue().(type) {
		case *Value_Null:
			d[k] = "null"
		case *Value_B:
			d[k] = fmt.Sprint(v.GetB())
		case *Value_N:
			d[k] = fmt.Sprint(v.GetN())
		case *Value_S:
			d[k] = v.GetS()
		}
	}
	return d
}

// TemplateValueBuilder applies value templates.
type TemplateValueBuilder interface {
	Build(value *Value, r TemplateSource) (*Value, error)
}

func NewTemplateValueBuilder(
	bodyBF func(*Value_Body) BodyBuilder,
	utilBF func(*Value_Util) UtilBuilder,
	headerBF func(*Value_Header) HeaderBuilder,
	urlBF func(*Value_Url) URLBuilder,
	addBF func(*Value_Add, TemplateValueBuilder) AddBuilder,
	castBF func(*Value_Cast, TemplateValueBuilder) CastBuilder,
) TemplateValueBuilder {
	return &templateValueBuilder{
		bodyBF:   bodyBF,
		utilBF:   utilBF,
		headerBF: headerBF,
		urlBF:    urlBF,
		addBF:    addBF,
		castBF:   castBF,
	}
}

type templateValueBuilder struct {
	bodyBF   func(*Value_Body) BodyBuilder
	utilBF   func(*Value_Util) UtilBuilder
	headerBF func(*Value_Header) HeaderBuilder
	urlBF    func(*Value_Url) URLBuilder
	addBF    func(*Value_Add, TemplateValueBuilder) AddBuilder
	castBF   func(*Value_Cast, TemplateValueBuilder) CastBuilder
}

func (s *templateValueBuilder) Build(value *Value, r TemplateSource) (*Value, error) {
	v, err := s.build(value, r)
	if err != nil {
		return nil, errors.Wrapf(err, errors.InvalidValue, "cannot build template value %s", util.JSON(value))
	}
	return v, nil
}

func (s *templateValueBuilder) build(value *Value, r TemplateSource) (*Value, error) {
	if r == nil {
		return nil, errors.New(errors.InvalidArgument, "request is nil")
	}
	switch value.GetValue().(type) {
	case *Value_Null, *Value_B, *Value_N, *Value_S:
		return value.Clone(), nil
	case *Value_L:
		p := []*Value{}
		for i, x := range value.GetL().GetValues() {
			v, err := s.build(x, r)
			if err != nil {
				return nil, errors.Wrapf(err, errors.InvalidValue, "at %d index in list", i)
			}
			p = append(p, v)
		}
		return NewL(p), nil
	case *Value_M:
		p := map[string]*Value{}
		for k, x := range value.GetM().GetValues() {
			v, err := s.build(x, r)
			if err != nil {
				return nil, errors.Wrapf(err, errors.InvalidValue, "key %s in map", k)
			}
			p[k] = v
		}
		return NewM(p), nil
	case *Value_Header_:
		return s.headerBF(value.GetHeader()).Build(r.Header())
	case *Value_Body_:
		return s.bodyBF(value.GetBody()).Build(r.Body())
	case *Value_Url_:
		return s.urlBF(value.GetUrl()).Build(r.URL())
	case *Value_Util_:
		return s.utilBF(value.GetUtil()).Build()
	case *Value_Add_:
		return s.addBF(value.GetAdd(), s).Build(r)
	case *Value_Cast_:
		return s.castBF(value.GetCast(), s).Build(r)
	}
	return nil, errors.New(errors.UnknownError, "template value builder")
}

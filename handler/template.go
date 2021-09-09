package handler

import (
	"encoding/json"

	"github.com/berquerant/jsonhttp/pb"
)

func NewTemplateValueBuilder() pb.TemplateValueBuilder {
	valueConverter := pb.NewValueConverter()
	valueCaster := pb.NewValueCaster()
	valueCoercer := pb.NewValueCoercer(valueCaster)
	bodyBF := func(b *pb.Value_Body) pb.BodyBuilder {
		return pb.NewBodyBuilder(b, valueConverter)
	}
	addBF := func(a *pb.Value_Add, s pb.TemplateValueBuilder) pb.AddBuilder {
		return pb.NewAddBuilder(a, valueCaster, s)
	}
	castBF := func(c *pb.Value_Cast, s pb.TemplateValueBuilder) pb.CastBuilder {
		return pb.NewCastBuilder(c, valueCoercer, s)
	}
	return pb.NewTemplateValueBuilder(
		bodyBF,
		pb.NewUtilBuilder,
		pb.NewHeaderBuilder,
		pb.NewURLBuilder,
		addBF,
		castBF,
	)
}

func NewTemplatesBuilder() pb.TemplatesBuilder {
	return pb.NewTemplatesBuilder(NewTemplateValueBuilder(), pb.NewValueInverter())
}

func WriteResultFromSource(w ResultWriter, src pb.TemplateSource) error {
	var m map[string]interface{}
	if err := json.Unmarshal(src.Body(), &m); err != nil {
		return err
	}
	for k, v := range m {
		w.Body().Set(k, v)
	}
	for k, v := range *src.Header() {
		if len(v) > 0 {
			w.Headers().Set(k, v[0])
		}
	}
	return nil
}

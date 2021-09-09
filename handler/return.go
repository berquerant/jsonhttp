package handler

import (
	"net/http"
	"time"

	"github.com/berquerant/jsonhttp/internal/errors"
	"github.com/berquerant/jsonhttp/internal/util"
	"github.com/berquerant/jsonhttp/pb"
)

// ReturnHandler applies the request to the template.
func ReturnHandler(ret *pb.Action_Return) Handler {
	const tag = "[return]"
	return func(w ResultWriter, r *http.Request) error {
		var (
			c       = FromContext(r.Context())
			src     = pb.NewTemplateSource(r.URL, &r.Header, c.Body())
			doDelay = func(src pb.TemplateSource) error {
				if ret.GetDelay() == nil {
					return nil
				}
				v, err := NewTemplateValueBuilder().Build(ret.GetDelay(), src)
				if err != nil {
					c.Log().Error("%s at calculate sleep %v", tag, err)
					return err
				}
				d, err := pb.NewValueCaster().Int(v)
				if err != nil {
					return errors.Wrapf(err, errors.Handler, "%s sleep duration is not int %v", tag, util.JSON(v))
				}
				c.Log().Info("%s sleep %d ms", tag, d)
				time.Sleep(time.Duration(d) * time.Millisecond)
				return nil
			}
		)

		w.Status().Set(int(ret.GetStatus()))
		writeTemplate := func() error {
			b := NewTemplatesBuilder()
			for i, t := range ret.GetTemplates() {
				if err := b.Add(t, src); err != nil {
					c.Log().Error("%s template %d %s %v", tag, i, util.JSON(t), err)
					return err
				}
			}
			w.Status().Set(int(ret.GetStatus()))
			for k, v := range b.Headers() {
				w.Headers().Set(k, v)
			}
			for k, v := range b.Body() {
				w.Body().Set(k, v)
			}
			return doDelay(src)
		}
		switch ret.GetTemplateType() {
		case pb.Action_APPEND:
			if err := WriteResultFromSource(w, pb.NewTemplateSource(nil, &r.Header, c.Body())); err != nil {
				return errors.Wrapf(err, errors.Handler, "%s write result %s", tag, c.Body())
			}
			return writeTemplate()
		case pb.Action_SELECT:
			if len(ret.GetTemplates()) == 0 {
				if err := WriteResultFromSource(w, pb.NewTemplateSource(nil, &r.Header, c.Body())); err != nil {
					return errors.Wrapf(err, errors.Handler, "%s write result %s", tag, c.Body())
				}
				return doDelay(src)
			}
			return writeTemplate()
		}
		return errors.Newf(errors.Handler, "unknown template type %s", ret.GetTemplateType())
	}
}

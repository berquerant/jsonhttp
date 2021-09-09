package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/berquerant/jsonhttp/internal/errors"
	"github.com/berquerant/jsonhttp/internal/util"
	"github.com/berquerant/jsonhttp/pb"
)

// GatewayHandler wraps a request to other url.
func GatewayHandler(gw *pb.Action_Gateway) Handler {
	const tag = "[gateway]"
	return func(w ResultWriter, r *http.Request) error {
		var (
			c                    = FromContext(r.Context())
			templateValueBuilder = NewTemplateValueBuilder()
			err                  error
		)
		// build url
		var u string
		if u, err = func() (string, error) {
			path, err := templateValueBuilder.Build(gw.GetPath(), pb.NewTemplateSource(r.URL, &r.Header, c.Body()))
			if err != nil {
				return "", err
			}
			return path.GetS(), nil
		}(); err != nil {
			return err
		}
		// build headers and body
		var (
			headers map[string]string
			body    map[string]interface{}
		)
		if headers, body, err = func() (headers map[string]string, body map[string]interface{}, err error) {
			headers = map[string]string{}
			body = map[string]interface{}{}
			fromRaw := func() error {
				if err := json.Unmarshal(c.Body(), &body); err != nil {
					return errors.Wrapf(err, errors.Handler, "%s unmarshal body %s", tag, c.Body())
				}
				for k, v := range r.Header {
					if len(v) > 0 {
						headers[k] = v[0]
					}
				}
				return nil
			}
			writeTemplate := func() error {
				b := NewTemplatesBuilder()
				for i, t := range gw.GetTemplates() {
					if err := b.Add(t, pb.NewTemplateSource(r.URL, &r.Header, c.Body())); err != nil {
						c.Log().Error("%s template %d %s %v", tag, i, util.JSON(t), err)
						return err
					}
				}
				for k, v := range b.Headers() {
					headers[k] = v
				}
				for k, v := range b.Body() {
					body[k] = v
				}
				return nil
			}

			switch gw.GetTemplateType() {
			case pb.Action_APPEND:
				if err = fromRaw(); err != nil {
					return
				}
				if err = writeTemplate(); err != nil {
					return
				}
				return
			case pb.Action_SELECT:
				if len(gw.GetTemplates()) == 0 {
					err = fromRaw()
					return
				}
				err = writeTemplate()
				return
			}
			err = errors.Newf(errors.UnknownError, "unknown template type %s", gw.GetTemplateType())
			return
		}(); err != nil {
			return err
		}

		var requestBody bytes.Buffer
		{
			b, err := json.Marshal(body)
			if err != nil {
				return errors.Wrapf(err, errors.Handler, "%s marshal body %v", tag, body)
			}
			if _, err := requestBody.Write(b); err != nil {
				return errors.Wrapf(err, errors.Handler, "%s marshal body %v", tag, body)
			}
		}
		// set request timeout
		ctx := r.Context()
		if gw.GetTimeout() != nil {
			x, err := templateValueBuilder.Build(gw.GetTimeout(), pb.NewTemplateSource(r.URL, &r.Header, c.Body()))
			if err != nil {
				return errors.Wrapf(err, errors.Handler, "%s build timeout %s", tag, util.JSON(gw.GetTimeout()))
			}
			d, err := pb.NewValueCaster().Int(x)
			if err != nil {
				return errors.Wrapf(err, errors.Handler, "%s type cast timeout %s", tag, util.JSON(x))
			}
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(r.Context(), time.Duration(d)*time.Millisecond)
			defer cancel()
		}
		// build http request
		var req *http.Request
		if req, err = http.NewRequestWithContext(ctx, gw.GetMethodType().String(), u, &requestBody); err != nil {
			return errors.Wrapf(err, errors.Handler, "%s build request", tag)
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		// do http request
		c.Log().Info("%s request to %s", tag, u)
		c.Log().Debug(`%s request with headers %s body "%s"`, tag, util.JSON(headers), requestBody)
		var res *http.Response
		if res, err = http.DefaultClient.Do(req); err != nil {
			return errors.Wrapf(err, errors.Handler, "%s do request", tag)
		}
		defer res.Body.Close()
		// build response
		w.Status().Set(res.StatusCode)
		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.Wrapf(err, errors.Handler, "%s read response body", tag)
		}
		c.Log().Debug(`%s got response status %d headers %s body "%s"`, tag, res.StatusCode, util.JSON(res.Header), responseBody)
		writeTemplate := func() error {
			ru, err := url.Parse(u)
			if err != nil {
				errors.Wrapf(err, errors.Handler, "%s parse request url %s", tag, u)
			}
			src := pb.NewTemplateSource(ru, &res.Header, responseBody)
			builder := NewTemplatesBuilder()
			for i, t := range gw.GetResponseTemplates() {
				if err := builder.Add(t, src); err != nil {
					c.Log().Error("%s response template %d %s %v", tag, i, util.JSON(t), err)
					return err
				}
			}
			for k, v := range builder.Headers() {
				w.Headers().Set(k, v)
			}
			for k, v := range builder.Body() {
				w.Body().Set(k, v)
			}
			return nil
		}
		switch gw.GetResponseTemplateType() {
		case pb.Action_APPEND:
			if err := WriteResultFromSource(w, pb.NewTemplateSource(nil, &res.Header, responseBody)); err != nil {
				return errors.Wrapf(err, errors.Handler, "%s write result %s", tag, responseBody)
			}
			return writeTemplate()
		case pb.Action_SELECT:
			if len(gw.GetResponseTemplates()) == 0 {
				if err := WriteResultFromSource(w, pb.NewTemplateSource(nil, &res.Header, responseBody)); err != nil {
					return errors.Wrapf(err, errors.Handler, "%s write result %s", tag, responseBody)
				}
				return nil
			}
			return writeTemplate()
		}
		return errors.Newf(errors.UnknownError, "unknown response template type %s", gw.GetResponseTemplateType())
	}
}

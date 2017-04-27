package opentracing

import (
	"context"

	"net/http"

	log "github.com/Sirupsen/logrus"
	gcloudtracer "github.com/hellofresh/gcloud-opentracing"
	"config"
	"opentracing/appdash"
	opentracing "github.com/opentracing/opentracing-go"
)

const (
	// CtxSpanID is used to store the SpanID in a request's context
	CtxSpanID = 0
)

// Build a tracer based on the configuration provided
func Build(config config.Tracing) (opentracing.Tracer, error) {
	if config.IsGoogleCloudEnabled() {
		log.Debug("Using google cloud platform (stackdriver trace) as tracing system")

		tracer, err := gcloudtracer.NewTracer(
			context.Background(),
			gcloudtracer.WithLogger(log.StandardLogger()),
			gcloudtracer.WithProject(config.GoogleCloudTracing.ProjectID),
			gcloudtracer.WithJWTCredentials(gcloudtracer.JWTCredentials{
				Email:        config.GoogleCloudTracing.Email,
				PrivateKey:   []byte(config.GoogleCloudTracing.PrivateKey),
				PrivateKeyID: config.GoogleCloudTracing.PrivateKeyID,
			}),
		)
		if err != nil {
			return nil, err
		}

		return tracer, nil
	} else if config.IsAppdashEnabled() {
		server := appdash.NewServer(config.AppdashTracing.DSN, config.AppdashTracing.URL)

		appdashFields := log.WithFields(log.Fields{
			"appdash_dsn":    config.AppdashTracing.DSN,
			"appdash_web_ui": config.AppdashTracing.URL,
		})

		if config.AppdashTracing.URL != "" {
			appdashFields.Debug("Using local appdash server as tracing system")
			err := server.Listen()
			if err != nil {
				return nil, err
			}
		} else {
			appdashFields.Debug("Using remote appdash server as tracing system")
		}

		return server.GetTracer(), nil
	} else {
		log.Debug("Not using a tracer as tracing system")
		return &opentracing.NoopTracer{}, nil
	}
}

// FromContext creates a span from a context that contains a parent span
func FromContext(ctx context.Context, name string) opentracing.Span {
	span, _ := opentracing.StartSpanFromContext(ctx, name)
	return span
}

// ToContext sets a span to a context
func ToContext(r *http.Request, span opentracing.Span) *http.Request {
	return r.WithContext(opentracing.ContextWithSpan(r.Context(), span))
}

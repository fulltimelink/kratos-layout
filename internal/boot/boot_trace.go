package boot

import (
	"github.com/fulltimelink/kratos-layout/internal/conf"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"os"
)

type BootTrace struct {
	conf *conf.Bootstrap
}

func NewBootTrace(conf *conf.Bootstrap) *BootTrace {
	return &BootTrace{
		conf,
	}
}

func (t *BootTrace) Run() {
	var tp *tracesdk.TracerProvider
	currResource, err := t.resource()
	if err != nil {
		panic(err)
	}
	if t.conf.Trace.Enable {
		exporter, err := t.exporter(t.conf.Trace.Exporter)
		if err != nil {
			panic(err)
		}
		tp = tracesdk.NewTracerProvider(
			tracesdk.WithBatcher(exporter),
			tracesdk.WithResource(currResource),
		)
	} else {
		tp = tracesdk.NewTracerProvider(
			tracesdk.WithResource(currResource),
		)
	}
	otel.SetTracerProvider(tp)
}

func (t *BootTrace) exporter(types string) (tracesdk.SpanExporter, error) {
	var exporter tracesdk.SpanExporter
	var err error

	switch types {
	case "stdout":
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint())
	case "file":
		var osFile *os.File
		osFile, err = os.Create(t.conf.Trace.TraceFilePath)
		exporter, err = stdouttrace.New(
			stdouttrace.WithWriter(osFile))
	case "jaeger":
		exporter, err = jaeger.New(
			jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(t.conf.Trace.Endpoint),
			))
	}
	return exporter, err
}

func (t *BootTrace) resource() (*resource.Resource, error) {
	currResource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(t.conf.Server.Name),          //实例名称
			attribute.String("environment", t.conf.Server.Environment), // 相关环境
			attribute.String("ID", t.conf.Server.Version),              //版本
			attribute.String("token", t.conf.Trace.Token),              //token
		),
	)
	return currResource, err
}

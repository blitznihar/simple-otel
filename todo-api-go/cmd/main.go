package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"todo-api-go/internal/config"
	"todo-api-go/internal/db"
	httpapi "todo-api-go/internal/http"

	"github.com/joho/godotenv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func initTracer(cfg config.Config) (func(context.Context) error, error) {
	// If no OTEL endpoint set, just run without exporting.
	if cfg.OtelEndpoint == "" {
		log.Println("OTEL_EXPORTER_OTLP_ENDPOINT not set; running without trace export")
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(cfg.OtelServiceName),
			)),
		)
		otel.SetTracerProvider(tp)
		return tp.Shutdown, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpointURL(cfg.OtelEndpoint),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.OtelServiceName),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	shutdown, err := initTracer(cfg)
	if err != nil {
		log.Fatalf("otel init error: %v", err)
	}
	defer func() { _ = shutdown(context.Background()) }()

	_, coll, err := db.ConnectMongo(cfg)
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}
	log.Println("Mongo connected OK")

	h := httpapi.Handlers{Coll: coll}
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpapi.Router(h),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("todo-api-go listening on :%s", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}

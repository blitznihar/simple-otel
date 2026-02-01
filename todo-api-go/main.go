package main

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type Todo struct {
	ID         int    `json:"id" bson:"id"`
	Title      string `json:"title" bson:"title"`
	Note       string `json:"note" bson:"note"`
	DueDate    string `json:"due_date,omitempty" bson:"due_date,omitempty"`
	CreatedAt  string `json:"created_at,omitempty" bson:"created_at,omitempty"`
	IsComplete bool   `json:"is_complete" bson:"is_complete"`
}

func initTracer() func(context.Context) error {
	ctx := context.Background()

	exp, _ := otlptracehttp.New(ctx)

	res, _ := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("todo-api-go"),
		),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	return tp.Shutdown
}

func main() {
	fmt.Println("Hello, Todo API in Go!")

	shutdown := initTracer()
	defer shutdown(context.Background())
}

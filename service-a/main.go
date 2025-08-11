package main

import (
    "bytes"
    "context"
    "encoding/json"
    "io"
    "log"
    "net/http"

    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type CepRequest struct {
    Cep string `json:"cep"`
}

func main() {
    tp := initTracer()
    defer func() { _ = tp.Shutdown(context.Background()) }()

    mux := http.NewServeMux()
    mux.Handle("/cep", otelhttp.NewHandler(http.HandlerFunc(handleCep), "HandleCep"))

    log.Println("Serviço A rodando na porta 8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func handleCep(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req CepRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }

    if len(req.Cep) != 8 {
        http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
        return
    }

    body, _ := json.Marshal(req)
    resp, err := otelhttp.Post(r.Context(), "http://service-b:8081/weather", "application/json", bytes.NewBuffer(body))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}

func initTracer() *sdktrace.TracerProvider {
    ctx := context.Background()
    exporter, err := otlptracehttp.New(ctx)
    if err != nil {
        log.Fatal(err)
    }
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String("service-a"),
        )),
    )
    otel.SetTracerProvider(tp)
    return tp
}

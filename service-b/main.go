package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

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

type WeatherResponse struct {
    City  string  `json:"city"`
    TempC float64 `json:"temp_C"`
    TempF float64 `json:"temp_F"`
    TempK float64 `json:"temp_K"`
}

func main() {
    tp := initTracer()
    defer func() { _ = tp.Shutdown(context.Background()) }()

    mux := http.NewServeMux()
    mux.Handle("/weather", otelhttp.NewHandler(http.HandlerFunc(handleWeather), "HandleWeather"))

    log.Println("Serviço B rodando na porta 8081")
    log.Fatal(http.ListenAndServe(":8081", mux))
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
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

    city, err := getCityFromCEP(req.Cep)
    if err != nil {
        http.Error(w, "can not find zipcode", http.StatusNotFound)
        return
    }

    tempC, err := getTempFromWeatherAPI(city)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    resp := WeatherResponse{
        City:  city,
        TempC: tempC,
        TempF: tempC*1.8 + 32,
        TempK: tempC + 273,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func getCityFromCEP(cep string) (string, error) {
    url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    var data map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return "", err
    }
    if _, ok := data["erro"]; ok {
        return "", fmt.Errorf("not found")
    }
    city := fmt.Sprintf("%v", data["localidade"])
    return city, nil
}

func getTempFromWeatherAPI(city string) (float64, error) {
    apiKey := os.Getenv("WEATHER_API_KEY")
    url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, city)
    resp, err := http.Get(url)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()
    var data map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return 0, err
    }
    if current, ok := data["current"].(map[string]interface{}); ok {
        if tempC, ok := current["temp_c"].(float64); ok {
            return tempC, nil
        }
    }
    return 0, fmt.Errorf("temp not found")
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
            semconv.ServiceNameKey.String("service-b"),
        )),
    )
    otel.SetTracerProvider(tp)
    return tp
}

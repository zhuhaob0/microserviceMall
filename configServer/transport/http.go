package transport

import (
	"context"
	"encoding/json"
	"errors"
	"final-design/configServer/endpoint"
	"net/http"

	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints endpoint.ConfigEndpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Path("/metrics").Handler(promhttp.Handler())

	r.Methods("GET").Path("/master/{.*}").Handler(kithttp.NewServer(
		endpoints.DownloadFileEndpoint,
		decodeDownloadRequest,
		encodeDownloadResponse,
		options...,
	))

	// create health check handler
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeConfigResponse,
		options...,
	))

	return r
}

// encodeStringResponse encode response to return，
// response是makeXXXEndpoint函数里的XXXResponse{}
func encodeConfigResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	// ret := response.(endpoint.DownloadResponse)
	// fmt.Println("ret.Content=", ret.Content)
	return json.NewEncoder(w).Encode(response)
}

func encodeDownloadResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	ret := response.(endpoint.DownloadResponse)
	return json.NewEncoder(w).Encode(ret)
}

func decodeDownloadRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var filename string = "." + r.URL.Path

	return endpoint.DownloadRequest{
		Filename: filename,
	}, nil
}

// decodeHealthCheckRequest decode request
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpoint.HealthRequest{}, nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

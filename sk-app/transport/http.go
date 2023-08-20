package transport

import (
	"context"
	"encoding/json"
	"errors"
	endpts "final-design/sk-app/endpoint"
	"final-design/sk-app/model"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	gozipkin "github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ErrBadRequest = errors.New("invalid request parameters")
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints endpts.SkAppEndpoints,
	zipkinTracer *gozipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	zipkinServer := zipkin.HTTPServerTrace(zipkinTracer, zipkin.Name("http-transport"))

	options := []kithttp.ServerOption{
		// kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}

	r.Methods("POST").Path("/sec/info").Handler(kithttp.NewServer(
		endpoints.GetSecInfoEndpoint,
		decodeSecInfoRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/sec/list").Handler(kithttp.NewServer(
		endpoints.GetSecInfoListEndpoint,
		decodeSecInfoListRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/sec/kill").Handler(kithttp.NewServer(
		endpoints.SecKillEndpoint,
		decodeSecKillRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/sec/test").Handler(kithttp.NewServer(
		endpoints.TestEndpoint,
		decodeTestRequest,
		encodeResponse,
		options...,
	))

	r.Path("/metrics").Handler(promhttp.Handler())

	// create health check handler
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeResponse,
		options...,
	))

	return r
}

// decodeUserRequest decode request params to struct
func decodeSecInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var secInfoRequest endpts.SecInfoRequest
	if err := json.NewDecoder(r.Body).Decode(&secInfoRequest); err != nil {
		return nil, err
	}
	return secInfoRequest, nil
}

func decodeSecInfoListRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	fmt.Println("decode sec list request")
	return nil, nil
}

func decodeSecKillRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var secRequest model.SecRequest
	if err := json.NewDecoder(r.Body).Decode(&secRequest); err != nil {
		return nil, err
	}
	secRequest.AccessToken = r.Header.Get("Authorization")
	return secRequest, nil
}
func decodeTestRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}

// decodeHealthCheckRequest decode request
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}

// encode errors from business-logic
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

// encodeArithmeticResponse encode response to return
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

package transport

import (
	"context"
	"encoding/json"
	"errors"
	endpts "final-design/user-service/endpoint"
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
	ErrBadRequest = errors.New("invalid request parameter")
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints endpts.UserEndpoints, zipkinTracer *gozipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	zipkinSever := zipkin.HTTPServerTrace(zipkinTracer, zipkin.Name("http-transport"))

	options := []kithttp.ServerOption{
		// kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		// kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinSever,
	}
	// ===============================普通用户注册和检查=======================================================
	r.Methods("POST").Path("/user/create").Handler(kithttp.NewServer(
		endpoints.CreateUserEndpoint,
		decodeCreateUserRequest,
		encodeUserResponse,
		options...,
	))
	r.Methods("POST").Path("/user/check").Handler(kithttp.NewServer(
		endpoints.UserEndpoint,
		decodeUserRequest,
		encodeUserResponse,
		options...,
	))

	// ==============================Admin用户注册和检查=======================================================
	r.Methods("POST").Path("/user/admin/create").Handler(kithttp.NewServer(
		endpoints.CreateAdminUserEndpoint,
		decodeCreateAdminUserRequest,
		encodeUserResponse,
		options...,
	))
	r.Methods("POST").Path("/user/admin/check").Handler(kithttp.NewServer(
		endpoints.AdminUserEndpoint,
		decodeAdminUserRequest,
		encodeUserResponse,
		options...,
	))

	// =================================健康检查============================================================
	r.Path("/metrics").Handler(promhttp.Handler())

	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeUserResponse,
		options...,
	))

	return r
}

// 解码检查User的request
func decodeUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var userRequest endpts.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		return nil, err
	}
	return userRequest, nil
}

// 解码检查AdminUser的request
func decodeAdminUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var userRequest endpts.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		return nil, err
	}
	return userRequest, nil
}

// 解码创建User的request
func decodeCreateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var createUserRequest endpts.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		return nil, err
	}
	return createUserRequest, nil
}

// 解码创建AdminUser的request
func decodeCreateAdminUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var createUserRequest endpts.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&createUserRequest); err != nil {
		return nil, err
	}
	return createUserRequest, nil
}

// encode errors from bussiness-logic
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
func encodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// decodeHealthCheckRequest decode request
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}

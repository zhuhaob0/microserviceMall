package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	endpts "final-design/sk-admin/endpoint"
	"final-design/sk-admin/model"

	"github.com/go-kit/kit/tracing/zipkin"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	gozipkin "github.com/openzipkin/zipkin-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ErrBadRequest = errors.New("invalid request parameter")
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints endpts.SkAdminEndpoints, zipkinTracer *gozipkin.Tracer, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	zipkinServer := zipkin.HTTPServerTrace(zipkinTracer, zipkin.Name("http-transport"))

	options := []kithttp.ServerOption{
		// kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		// kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
		kithttp.ServerErrorEncoder(encodeError),
		zipkinServer,
	}

	// =========================================商品管理====================================================
	r.Methods("GET").Path("/product/list").Handler(kithttp.NewServer(
		endpoints.GetProductEndpoint,
		decodeGetListRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/product/create").Handler(kithttp.NewServer(
		endpoints.CreateProductEndpoint,
		decodeCreateProductRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/product/update").Handler(kithttp.NewServer(
		endpoints.UpdateProductEndpoint,
		decodeUpdateProductRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/product/delete").Handler(kithttp.NewServer(
		endpoints.DeleteProductEndpoint,
		decodeDeleteProductRequest,
		encodeResponse,
		options...,
	))
	// ==========================================活动管理===================================================
	r.Methods("GET").Path("/activity/list").Handler(kithttp.NewServer(
		endpoints.GetActivityEndpoint,
		decodeGetListRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/activity/create").Handler(kithttp.NewServer(
		endpoints.CreateActivityEndpoint,
		decodeCreateActivityRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/activity/update").Handler(kithttp.NewServer(
		endpoints.UpdateActivityEndpoint,
		decodeUpdateActivityRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/activity/delete").Handler(kithttp.NewServer(
		endpoints.DeleteActivityEndpoint,
		decodeDeleteActivityRequest,
		encodeResponse,
		options...,
	))
	// ==========================================订单管理====================================================
	r.Methods("GET").Path("/order/list").Handler(kithttp.NewServer(
		endpoints.GetOrderListEndpoint,
		decodeGetOrderRequest,
		encodeResponse,
		options...,
	))

	r.Methods("POST").Path("/order/buyer").Handler(kithttp.NewServer(
		endpoints.GetBuyerOrderEndpoint,
		decodeGetBuyerOrderRequest,
		encodeResponse,
		options...,
	))
	// ==========================================健康检查====================================================
	r.Path("/metrics").Handler(promhttp.Handler())

	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeResponse,
		options...,
	))

	loggedRouter := handlers.LoggingHandler(os.Stderr, r)

	return loggedRouter
}

// 对获取所有商品或者活动请求解码
func decodeGetListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpts.GetListRequest{}, nil
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
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// decodeHealthCheckRequest decode request
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}

// =====================================decodeProduct===============================================================

func decodeCreateProductRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return nil, err
	}
	// fmt.Println("decodeCreateProduct, product=", product)
	return product, nil
}
func decodeUpdateProductRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return nil, err
	}
	return product, nil
}
func decodeDeleteProductRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		return nil, err
	}
	return product, nil
}

// =====================================decodeActivity==================================================================

func decodeCreateActivityRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var activity model.Activity
	if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
		return nil, err
	}
	return activity, nil
}
func decodeUpdateActivityRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var activity model.Activity
	if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
		return nil, err
	}
	return activity, nil
}
func decodeDeleteActivityRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var activity model.Activity
	if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
		return nil, err
	}
	return activity, nil
}

// =====================================decodeOrder==================================================================
func decodeGetOrderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.GetListRequest{}, nil
}

func decodeGetBuyerOrderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var orderReq endpts.BuyerOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		return nil, err
	}
	return orderReq, nil
}

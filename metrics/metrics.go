package metrics

import (
	"context"
	"fmt"
	prometheusapi "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"os"
	"time"
)

var prometheusAddress = os.Getenv("PROMETHEUS_ADDRESS")

func newPrometheusClient(serverAddress string) (prometheus.API, error) {
	client, err := prometheusapi.NewClient(prometheusapi.Config{
		Address: serverAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating metrics client: %v", err)
	}
	return prometheus.NewAPI(client), nil
}

func GetAvgSvcTraffic(appName, serviceName string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		(sum(
			rate(istio_request_bytes_sum{app="`+appName+`", svc="`+serviceName+`", source_workload!="unknown", destination_workload!="unknown"}[5m])
		) by (source_workload, destination_workload)
		+
		sum(
			rate(istio_response_bytes_sum{app="`+appName+`", svc="`+serviceName+`", source_workload!="unknown", destination_workload!="unknown"}[5m])
		) by (source_workload, destination_workload)
		)
		or 
		sum(
			rate(istio_tcp_sent_bytes_total{app="`+appName+`", svc="`+serviceName+`", source_workload!="unknown", destination_workload!="unknown"}[5m]) 
			+ 
			rate(istio_tcp_received_bytes_total{app="`+appName+`", svc="`+serviceName+`", source_workload!="unknown", destination_workload!="unknown"}[5m])
		) by (source_workload, destination_workload)
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAvgNodeLatencies(nodeName string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		(rate(node_latency_sum{origin_node="`+nodeName+`"}[5m]) / rate(node_latency_count{origin_node="`+nodeName+`"}[5m])) * 1000
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

package main

import (
	"fmt"
	"github.com/amarchese96/telemetry-service/metrics"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAvgSvcTraffic(c *gin.Context) {
	appName := c.Query("app")
	serviceName := c.Query("svc")

	trafficValues := map[string]float64{}

	results, warnings, err := metrics.GetAvgSvcTraffic(appName, serviceName)
	fmt.Println(warnings)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}

	for _, result := range results {
		if string(result.Metric["source_workload"]) == serviceName {
			trafficValues[string(result.Metric["destination_workload"])] = float64(result.Value)
		} else if string(result.Metric["destination_workload"]) == serviceName {
			trafficValues[string(result.Metric["source_workload"])] = float64(result.Value)
		}
	}
	c.IndentedJSON(http.StatusOK, trafficValues)
}

func getAvgNodeLatencies(c *gin.Context) {
	nodeName := c.Query("node")

	latencyValues := map[string]float64{}

	results, _, err := metrics.GetAvgNodeLatencies(nodeName)
	//fmt.Println(warnings)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}

	for _, result := range results {
		latencyValues[string(result.Metric["destination_node"])] = float64(result.Value)
	}
	c.IndentedJSON(http.StatusOK, latencyValues)
}

func main() {
	router := gin.Default()
	router.GET("/metrics/svc/avg-traffic", getAvgSvcTraffic)
	router.GET("/metrics/node/avg-latencies", getAvgNodeLatencies)

	router.Run("0.0.0.0:8080")
}

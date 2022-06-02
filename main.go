package main

import (
	"fmt"
	"github.com/amarchese96/telemetry-service/metrics"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getAvgTraffic(c *gin.Context) {
	appName := c.Query("app")
	nodeName := c.Query("node")

	trafficValues := map[string]float64{}

	results, warnings, err := metrics.GetAvgTraffic(appName, nodeName)
	fmt.Println(warnings)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}

	for _, result := range results {
		if string(result.Metric["source_workload"]) == nodeName {
			trafficValues[string(result.Metric["destination_workload"])] = float64(result.Value)
		} else if string(result.Metric["destination_workload"]) == nodeName {
			trafficValues[string(result.Metric["source_workload"])] = float64(result.Value)
		}
	}
	c.IndentedJSON(http.StatusOK, trafficValues)
}

func main() {
	router := gin.Default()
	router.GET("/metrics/avg-traffic", getAvgTraffic)

	router.Run("0.0.0.0:8081")
}

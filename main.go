package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/influxdata/influxdb1-client"
	client "github.com/influxdata/influxdb1-client/v2"
)

// Application configuration constants.
const (
	DATABASE = "system_metrics"
)

func main() {
	httpClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	handleError("Couldn't connect to the InfluxDB server", err)
	defer httpClient.Close()

	points, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  DATABASE,
		Precision: "us",
	})
	handleError("Couldn't create batch points", err)

	go func() {
		rand.Seed(time.Now().UTC().UnixNano())

		for i := 0; i < 30; i++ {
			cpuTemperature := rand.Float64() * 100.0
			gpuTemperature := rand.Float64() * 100.0

			tags := map[string]string{
				"cluster_index": strconv.Itoa(rand.Intn(3)),
				"host_index":    strconv.Itoa(rand.Intn(3)),
			}

			fields := map[string]interface{}{
				"cpu_temperature": cpuTemperature,
				"gpu_temperature": gpuTemperature,
			}

			point, err := client.NewPoint("system_stat", tags, fields, time.Now())
			handleError("Couldn't create a point", err)
			points.AddPoint(point)
		}

		err = httpClient.Write(points)
		handleError("Couldn't write data to the server", err)
	}()

	time.Sleep(5 * time.Second)

	query := client.Query{
		Command:  "SELECT * FROM system_stat",
		Database: DATABASE,
	}

	resp, err := httpClient.Query(query)
	handleError("Couldn't query the database", err)
	handleError("Error from the database", resp.Error())
	jsonBytes, err := json.MarshalIndent(resp, "", "    ")
	handleError("Couldn't marshal the response to JSON", err)

	fmt.Println(string(jsonBytes))
}

func handleError(message string, err error) {
	if err != nil {
		log.Fatalf("%s: %s\n", message, err)
	}
}

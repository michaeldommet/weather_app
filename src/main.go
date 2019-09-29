package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

//Insert function
func insertQuery(w weather, i inside) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://IPADD:8086",
		Username: "****",
		Password: "****",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	//q := client.NewQuery("SELECT \"temperature C\"  FROM \"sensor_data\".\"autogen\".\"rpi-bme280\"", "", "")
	//if response, err := c.Query(q); err == nil && response.Error() == nil {
	//	fmt.Println(response.Results)
	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "sensor_data",
		Precision: "s",
	})
	// Create a point and add to batch
	temptags := map[string]string{"location": "Home"}
	tempfields := map[string]interface{}{
		"outside-Temp":     w.Temp,
		"outside-Pressure": w.Pressure,
		"outside-Humidity": w.Humidity,
		"outside-TempMax":  w.TempMax,
		"outside-TempMin":  w.TempMin,
		"inside-Temp":      i.InTemp,
		"inside-Pressure":  i.InPressure,
		"inside-Humidity":  i.InHumidity,
	}
	pt, err := client.NewPoint("rpi-bme280", temptags, tempfields, time.Now())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	bp.AddPoint(pt)
	// Write the batch
	c.Write(bp)
}

func main() {

	url := fmt.Sprint("http://api.openweathermap.org/data/2.5/weather?q=Berlin,de&units=metric&APPID=***********")

	//Build the API requst:
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Could not make the api call", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return
	}
	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	var record WeatherStatus

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	insidereads := sensorData()
	outtemp := record.Main
	insertQuery(outtemp, insidereads)

}

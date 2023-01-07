package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) Live(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "UP"})
}

func (server *Server) Ready(ctx *gin.Context) {

	geoDB := "UP"
	weatherAPI := "UP"

	// GeoDB
	url := server.config.GeoDBAddress + "/cities"
	apiKey := server.config.GeoDBAPIKey
	host := server.config.GeoDBAPIHost
	err := PingService(url, apiKey, host)

	if err != nil {
		log.Println(err)
		geoDB = "DOWN"
	}

	// Weather API
	url = server.config.AerisWeatherAddress + "/observations/" + fmt.Sprintf("%f", 51.507222222) + "," + fmt.Sprintf("%f", -0.1275)
	apiKey = server.config.AerisWeatherAPIKey
	host = server.config.AerisWeatherAPIHost
	err = PingService(url, apiKey, host)

	if err != nil {
		log.Println(err)
		weatherAPI = "DOWN"
	}

	// Return health status
	ctx.JSON(http.StatusOK, gin.H{"status": gin.H{
		"geoDB":      geoDB,
		"weatherAPI": weatherAPI,
	}})
}

func PingService(url string, apiKey string, host string) error {

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", apiKey)
	req.Header.Add("X-RapidAPI-Host", host)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Cannot unmarshal Response")
		return errors.New("GeoDB service unavailable!")
	}

	defer res.Body.Close()

	log.Println(res.StatusCode)
	if res.StatusCode != 200 {
		return errors.New("GeoDB service unavailable!")
	}

	return nil
}

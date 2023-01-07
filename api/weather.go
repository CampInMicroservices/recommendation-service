package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type getLocationWeather struct {
	Lat  float64 `form:"lat" binding:"required"`
	Long float64 `form:"long" binding:"required"`
}

type weatherApiResponse struct {
	Success  bool        `json:"success"`
	Error    interface{} `json:"error"`
	Response struct {
		ID         string `json:"id"`
		DataSource string `json:"dataSource"`
		Loc        struct {
			Long float64 `json:"long"`
			Lat  float64 `json:"lat"`
		} `json:"loc"`
		Place struct {
			Name    string `json:"name"`
			City    string `json:"city"`
			State   string `json:"state"`
			Country string `json:"country"`
		} `json:"place"`
		Profile struct {
			Tz       string `json:"tz"`
			Tzname   string `json:"tzname"`
			Tzoffset int    `json:"tzoffset"`
			IsDST    bool   `json:"isDST"`
			ElevM    int    `json:"elevM"`
			ElevFT   int    `json:"elevFT"`
		} `json:"profile"`
		ObTimestamp int       `json:"obTimestamp"`
		ObDateTime  time.Time `json:"obDateTime"`
		Ob          struct {
			Type                string      `json:"type"`
			Timestamp           int         `json:"timestamp"`
			DateTimeISO         time.Time   `json:"dateTimeISO"`
			RecTimestamp        int         `json:"recTimestamp"`
			RecDateTimeISO      time.Time   `json:"recDateTimeISO"`
			TempC               int         `json:"tempC"`
			TempF               int         `json:"tempF"`
			DewpointC           int         `json:"dewpointC"`
			DewpointF           int         `json:"dewpointF"`
			Humidity            int         `json:"humidity"`
			PressureMB          int         `json:"pressureMB"`
			PressureIN          float64     `json:"pressureIN"`
			SpressureMB         int         `json:"spressureMB"`
			SpressureIN         float64     `json:"spressureIN"`
			AltimeterMB         int         `json:"altimeterMB"`
			AltimeterIN         float64     `json:"altimeterIN"`
			WindKTS             float64     `json:"windKTS"`
			WindKPH             float64     `json:"windKPH"`
			WindMPH             float64     `json:"windMPH"`
			WindSpeedKTS        float64     `json:"windSpeedKTS"`
			WindSpeedKPH        float64     `json:"windSpeedKPH"`
			WindSpeedMPH        float64     `json:"windSpeedMPH"`
			WindDirDEG          float64     `json:"windDirDEG"`
			WindDir             string      `json:"windDir"`
			WindGustKTS         interface{} `json:"windGustKTS"`
			WindGustKPH         interface{} `json:"windGustKPH"`
			WindGustMPH         interface{} `json:"windGustMPH"`
			FlightRule          string      `json:"flightRule"`
			VisibilityKM        float64     `json:"visibilityKM"`
			VisibilityMI        float64     `json:"visibilityMI"`
			Weather             string      `json:"weather"`
			WeatherShort        string      `json:"weatherShort"`
			WeatherCoded        string      `json:"weatherCoded"`
			WeatherPrimary      string      `json:"weatherPrimary"`
			WeatherPrimaryCoded string      `json:"weatherPrimaryCoded"`
			CloudsCoded         string      `json:"cloudsCoded"`
			Icon                string      `json:"icon"`
			HeatindexC          float64     `json:"heatindexC"`
			HeatindexF          int         `json:"heatindexF"`
			WindchillC          float64     `json:"windchillC"`
			WindchillF          int         `json:"windchillF"`
			FeelslikeC          float64     `json:"feelslikeC"`
			FeelslikeF          int         `json:"feelslikeF"`
			IsDay               bool        `json:"isDay"`
			Sunrise             int         `json:"sunrise"`
			SunriseISO          time.Time   `json:"sunriseISO"`
			Sunset              int         `json:"sunset"`
			SunsetISO           time.Time   `json:"sunsetISO"`
			SnowDepthCM         interface{} `json:"snowDepthCM"`
			SnowDepthIN         interface{} `json:"snowDepthIN"`
			PrecipMM            int         `json:"precipMM"`
			PrecipIN            int         `json:"precipIN"`
			SolradWM2           int         `json:"solradWM2"`
			SolradMethod        string      `json:"solradMethod"`
			CeilingFT           int         `json:"ceilingFT"`
			CeilingM            float64     `json:"ceilingM"`
			Light               int         `json:"light"`
			Uvi                 interface{} `json:"uvi"`
			Qc                  string      `json:"QC"`
			QCcode              int         `json:"QCcode"`
			TrustFactor         int         `json:"trustFactor"`
			Sky                 int         `json:"sky"`
		} `json:"ob"`
		Raw        string `json:"raw"`
		RelativeTo struct {
			Lat        float64 `json:"lat"`
			Long       float64 `json:"long"`
			Bearing    int     `json:"bearing"`
			BearingENG string  `json:"bearingENG"`
			DistanceKM float64 `json:"distanceKM"`
			DistanceMI float64 `json:"distanceMI"`
		} `json:"relativeTo"`
	} `json:"response"`
}

func (server *Server) GetLocationWeather(ctx *gin.Context) {

	// Check if request has lat and long field in URI.
	var req getLocationWeather
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	url := server.config.AerisWeatherAddress + "/observations/" + fmt.Sprintf("%f", req.Lat) + "," + fmt.Sprintf("%f", req.Long)
	apiReq, _ := http.NewRequest("GET", url, nil)

	apiReq.Header.Add("X-RapidAPI-Key", server.config.AerisWeatherAPIKey)
	apiReq.Header.Add("X-RapidAPI-Host", server.config.AerisWeatherAPIHost)

	res, _ := http.DefaultClient.Do(apiReq)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panic("Could not close Body")
		}
	}(res.Body)
	body, _ := ioutil.ReadAll(res.Body)

	var r weatherApiResponse
	err := json.Unmarshal([]byte(string(body)), &r)
	if err != nil {
		log.Panic("Cannot unmarshal Response")
	}

	ctx.JSON(http.StatusOK, r)
}

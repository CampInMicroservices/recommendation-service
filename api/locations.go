package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"

	gql "github.com/graphql-go/graphql"
)

type graphRequest struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

type data struct {
	ID          int     `json:"id"`
	WikiDataID  string  `json:"wikiDataId"`
	Type        string  `json:"type"`
	City        string  `json:"city"`
	Name        string  `json:"name"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionCode  string  `json:"regionCode"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Population  int     `json:"population"`
}

type LocationsResponse struct {
	Data  []data `json:"data"`
	Error string
}

func (server *Server) GetPopularLocations(ctx *gin.Context) {

	var dataType = gql.NewObject(
		gql.ObjectConfig{
			Name: "Data",
			Fields: gql.Fields{
				"id": &gql.Field{
					Type: gql.Int,
				},
				"wikiDataId": &gql.Field{
					Type: gql.String,
				},
				"Type": &gql.Field{
					Type: gql.String,
				},
				"city": &gql.Field{
					Type: gql.String,
				},
				"name": &gql.Field{
					Type: gql.String,
				},
				"country": &gql.Field{
					Type: gql.String,
				},
				"countryCode": &gql.Field{
					Type: gql.String,
				},
				"region": &gql.Field{
					Type: gql.String,
				},
				"regionCode": &gql.Field{
					Type: gql.String,
				},
				"latitude": &gql.Field{
					Type: gql.Float,
				},
				"longitude": &gql.Field{
					Type: gql.Float,
				},
				"population": &gql.Field{
					Type: gql.Int,
				},
			},
		},
	)

	var locationsResponseType = gql.NewObject(
		gql.ObjectConfig{
			Name: "LocationsResponse",
			Fields: gql.Fields{
				"data": &gql.Field{
					Type: gql.NewList(dataType),
				},
				"error": &gql.Field{
					Type: gql.String,
				},
			},
		},
	)

	// Schema
	fields := gql.Fields{
		"cities": &gql.Field{
			Type: locationsResponseType,
			Resolve: func(p gql.ResolveParams) (interface{}, error) {

				offset := rand.Intn(100)

				url := server.config.GeoDBAddress + "/cities"

				req, _ := http.NewRequest("GET", url, nil)

				req.Header.Add("X-RapidAPI-Key", server.config.GeoDBAPIKey)
				req.Header.Add("X-RapidAPI-Host", server.config.GeoDBAPIHost)

				params := req.URL.Query()
				params.Add("offset", fmt.Sprintf("%v", offset))
				params.Add("limit", fmt.Sprintf("%v", 5))
				params.Add("sort", "-population")
				params.Add("countryIds", "AT,CH,DE,FI,FR,GB,IT,NL,PL,SI")
				req.URL.RawQuery = params.Encode()

				res, _ := http.DefaultClient.Do(req)

				if res.StatusCode != 200 {
					log.Println("Slow down! Too many requests on GeoDB API.")
					r := LocationsResponse{Error: "Slow down! Too many requests on GeoDB API."}
					return r, nil
				}

				defer res.Body.Close()
				body, _ := ioutil.ReadAll(res.Body)

				var r LocationsResponse
				err := json.Unmarshal(body, &r)
				if err != nil {
					log.Panic("Cannot unmarshal LocationResponse")
				}

				return r, nil
			},
		},
	}
	rootQuery := gql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := gql.SchemaConfig{Query: gql.NewObject(rootQuery)}
	schema, err := gql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Parse graphQL request.
	var graphRequest graphRequest
	if err := json.NewDecoder(ctx.Request.Body).Decode(&graphRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "GraphQL request in wrong format."})
		ctx.Abort()
		return
	}

	params := gql.Params{
		Context:        ctx.Request.Context(),
		Schema:         schema,
		RequestString:  graphRequest.Query,
		VariableValues: graphRequest.Variables,
		OperationName:  graphRequest.Operation,
	}

	r := gql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}

	ctx.JSON(http.StatusOK, r)
}

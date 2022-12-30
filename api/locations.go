package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	gclient "github.com/machinebox/graphql"

	gql "github.com/graphql-go/graphql"
)

type graphRequest struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

func (server *Server) GetPopularLocations(ctx *gin.Context) {

	graphqlClient := gclient.NewClient(server.config.GeoDBAddress)

	// Schema
	fields := gql.Fields{
		"hello": &gql.Field{
			Type: gql.String,
			Resolve: func(p gql.ResolveParams) (interface{}, error) {
				graphqlRequest := gclient.NewRequest(`
					{countries(first:10) {
						edges {
							node {
								name
							}
						}
					}}
				`)

				graphqlRequest.Header.Set("X-RapidAPI-Key", server.config.GeoDBAPIKey)
				graphqlRequest.Header.Set("X-RapidAPI-Host", server.config.GeoDBAPIHost)

				var graphqlResponse interface{}
				if err := graphqlClient.Run(context.Background(), graphqlRequest, &graphqlResponse); err != nil {
					log.Fatalf("Error querying GeoDB, error: %v", err)
				}

				rJSON, _ := json.Marshal(graphqlResponse)

				return string(rJSON), nil
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

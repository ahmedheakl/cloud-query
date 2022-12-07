package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

type singleton struct {
	client *athena.Client
}

var athena_client *singleton

func GetOrCreateClient() *athena.Client {
	if athena_client == nil {
		cfg, _ := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
			o.Region = "me-south-1"
			return nil
		})
		athena_client = &singleton{
			client: athena.NewFromConfig(cfg),
		}
		return athena_client.client
	} else {
		return athena_client.client
	}
}

func GetResults(QueryID *string) ([]types.Row, error) {
	client := GetOrCreateClient()

	executionParams := &athena.GetQueryExecutionInput{
		QueryExecutionId: QueryID,
	}

	// poll query state, if success get results and return
	for {
		out, _ := client.GetQueryExecution(context.TODO(), executionParams)
		switch out.QueryExecution.Status.State {
		case types.QueryExecutionStateQueued, types.QueryExecutionStateRunning:
			time.Sleep(1 * time.Second)

		case types.QueryExecutionStateCancelled, types.QueryExecutionStateFailed:
			return nil, errors.New("query failed")

		case types.QueryExecutionStateSucceeded:
			resultsParams := &athena.GetQueryResultsInput{
				QueryExecutionId: QueryID,
			}

			data, _ := client.GetQueryResults(context.TODO(), resultsParams)
			return data.ResultSet.Rows, nil
		}
	}
}

func ExecuteQuery(query string) []byte {
	client := GetOrCreateClient()
	log.Printf("started query: %s", query)
	cntxt := &types.QueryExecutionContext{
		Catalog:  aws.String(athena_catalog),
		Database: aws.String(athena_database),
	}

	conf := &types.ResultConfiguration{
		OutputLocation: aws.String(athena_spill),
	}

	params := &athena.StartQueryExecutionInput{
		QueryString:           aws.String(query),
		ResultConfiguration:   conf,
		QueryExecutionContext: cntxt,
	}

	resp, err := client.StartQueryExecution(context.TODO(), params)
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}

	queryId := resp.QueryExecutionId
	log.Printf("query id: %s", *queryId)
	query_result, err := GetResults(queryId)
	if err != nil {
		log.Printf("query id %s result: %s", *queryId, err.Error())
		return []byte(fmt.Sprintf("Error executing query %s: ", err.Error()))
	}

	json_result, _ := json.Marshal(query_result)
	log.Printf("query id %s result: %s", *queryId, string(json_result))

	return (json_result)
}

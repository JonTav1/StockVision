package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type stock struct {
	Ticker       string  `json:"ticker" dynamodbav:"ticker"`
	Amount       int     `json:"amount" dynamodbav:"amount"`
	AveragePrice float64 `json:"averagePrice" dynamodbav:"averagePrice"`
	Username     string  `json:"username" dynamodbav:"username"`
}

var svc *dynamodb.DynamoDB

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	username, exists := request.PathParameters["username"]
	if !exists {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing path parameter 'username'",
		}, nil
	}
	stocks := getStocks(username)

	jsonResponse, err := json.Marshal(stocks)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonResponse),
	}
	return response, nil
}
func getStocks(username string) []stock {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("serverless-stock-app-stocks"),
		KeyConditionExpression: aws.String("#uname = :username"),
		ExpressionAttributeNames: map[string]*string{
			"#uname": aws.String("username"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(username),
			},
		},
	}

	// Perform the query operation
	result, err := svc.Query(input)
	if err != nil {
		log.Printf("Error querying stocks: %v", err)
		return nil
	}

	var stocks []stock
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &stocks)
	if err != nil {
		log.Printf("Error unmarshaling stocks: %v", err)
		return nil
	}

	return stocks
}
func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc = dynamodb.New(sess)
}
func main() {
	lambda.Start(handler)
}

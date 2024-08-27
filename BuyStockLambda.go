package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type stock struct {
	Ticker       string  `json:"ticker" dynamodbav:"ticker"`
	Amount       int     `json:"amount" dynamodbav:"amount"`
	AveragePrice float64 `json:"averagePrice" dynamodbav:"averagePrice"`
	Username     string  `json:"username" dynamodbav:"username"`
}

var svc *dynamodb.DynamoDB

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var stockChoice stock
	err := json.Unmarshal([]byte(request.Body), &stockChoice)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	stocks := getStocks(stockChoice.Username)
	stockDB := findStock(stockChoice.Ticker, stocks)
	stockDB.AveragePrice = (stockDB.AveragePrice*float64(stockDB.Amount) + stockChoice.AveragePrice*float64(stockChoice.Amount)) / float64(stockChoice.Amount+stockDB.Amount)

	stockDB.Amount = stockDB.Amount + stockChoice.Amount
	editStockDB(stockDB)

	jsonResponse, err := json.Marshal(stockDB)
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

	// Unmarshal the query result into a slice of Stock structs
	var stocks []stock
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &stocks)
	if err != nil {
		log.Printf("Error unmarshaling stocks: %v", err)
		return nil
	}

	return stocks
}
func findStock(stockTicker string, stocks []stock) stock {
	var stockChoice stock
	for _, stock := range stocks {
		if stock.Ticker == stockTicker {
			stockChoice = stock
		}
	}
	return stockChoice
}
func editStockDB(stock stock) {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				N: aws.String(strconv.Itoa(stock.Amount)),
			},
			":c": {
				N: aws.String(strconv.FormatFloat(stock.AveragePrice, 'f', -1, 64)),
			},
		},
		TableName: aws.String("serverless-stock-app-stocks"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(stock.Username),
			},
			"ticker": {
				S: aws.String(stock.Ticker),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set amount=:r, averagePrice=:c"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		log.Fatalf("Got error calling UpdateItem: %s", err)
	}

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

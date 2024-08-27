package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type User struct {
	Username  string `json:"username" dynamodbav:"username"`
	Password  string `json:"password" dynamodbav:"password"`
	FirstName string `json:"firstname" dynamodbav:"firstname"`
	LastName  string `json:"lastname" dynamodbav:"lastname"`
}
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var svc *dynamodb.DynamoDB

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var loginReq LoginRequest
	err := json.Unmarshal([]byte(request.Body), &loginReq)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	username := loginReq.Username
	password := loginReq.Password
	checkPassword(username, password)
	user := getUser(username)
	responseBody := User{
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
	}
	jsonResponse, err := json.Marshal(responseBody)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonResponse),
	}
	return response, nil
}

func getUser(username string) User {
	// Define the put input
	input := &dynamodb.GetItemInput{
		TableName: aws.String("serverless-stock-app-users"),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
		},
	}
	response, _ := svc.GetItem(input)

	var newuser User
	dynamodbattribute.UnmarshalMap(response.Item, &newuser)
	return newuser
}
func checkPassword(usernameTest string, passwordTest string) {
	newuser := getUser(usernameTest)
	if passwordTest != newuser.Password {
		panic("Password does not match")
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

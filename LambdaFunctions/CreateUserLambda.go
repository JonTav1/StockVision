package main

import (
	"encoding/json"
	"fmt"
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

var svc *dynamodb.DynamoDB

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req User
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	usernameExists(req.Username)
	user := makeUser(req.FirstName, req.LastName, req.Username, req.Password)
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

func makeUser(firstName string, lastName string, username string, password string) User { //creates a User object with arguments given, and stores it
	user := User{username, password, firstName, lastName}
	storeUser(user)
	return user
}
func usernameExists(usernameTest string) { //checks if the username exists for creation, if it already, does, throws error
	user := getUser(usernameTest)
	fmt.Println(user.Username)

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
func storeUser(user User) {
	// Marshal the user struct into a DynamoDB attribute value
	av, _ := dynamodbattribute.MarshalMap(user)
	fmt.Println(av)

	// Define the put input
	input := &dynamodb.PutItemInput{
		TableName:           aws.String("serverless-stock-app-users"),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(username)"),
	}
	// Perform the put operation
	_, _ = svc.PutItem(input)

	fmt.Println("Successfully stored user in DynamoDB")
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

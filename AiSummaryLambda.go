package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	trade_indicators "github.com/go-whale/trade-indicators"
	"github.com/sashabaranov/go-openai"
	"os"
	"strings"
	"time"
)

var prices []float64 = nil
var period int = 31
var volumeStock []uint64 = nil
var volumeCrypto []float64 = nil
var OPEN_AI_KEY string = ""
var APCA_API_KEY_ID string = ""
var APCA_API_SECRET_KEY string = ""

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ticker, tickerExists := request.PathParameters["ticker"]
	choice, choiceExists := request.PathParameters["choice"]
	if !tickerExists || !choiceExists {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing path parameter 'ticker' or 'choice'",
		}, nil
	}
	summary := stockSummary(ticker, choice)

	jsonResponse, err := json.Marshal(summary)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonResponse),
	}
	return response, nil
}

func stockSummary(ticker string, choice string) string {

	//
	stockSummary := Ai_request_helper(ticker, choice)

	return stockSummary
}
func Ai_request_helper(ticker string, choice string) string {
	stockFinder()
	prompt := "give me a paragraph summary of the stock" + ticker + "saying what they do, and any other releveant information. then " +
		"Give me a rating, from 1-100 on whether I should buy or sell the stock" + ticker + "based on the " +
		"past 31 days of prices, 14 day RSI, 26 day EMA, MACD based on the 12 and 26 day ema, and the volatility\n" +
		"explain in detail, give me a lengthy response."
	prices := getPrices(ticker, choice)
	pricesString := float64SliceToString(prices)
	prompt = prompt + "prices: " + pricesString + "\n"
	RSI := calculateRSI()
	rsiString := float64SliceToString(RSI)
	prompt = prompt + "prices: " + rsiString + "\n"
	EMA := calculateEMA()
	emaString := float64SliceToString(EMA)
	prompt = prompt + "prices: " + emaString + "\n"
	MACD := calculateMACD()
	macdString := float64SliceToString(MACD)
	prompt = prompt + "prices: " + macdString + "\n"
	Volatility := calculateVolatility()
	volatilityString := fmt.Sprintf("%f", Volatility) // s == "123.456000"
	prompt = prompt + "prices: " + volatilityString + "\n"
	response := Ai_request(prompt)
	return response
}

func stockFinder() {

	client := alpaca.NewClient(alpaca.ClientOpts{

		BaseURL: "https://api.alpaca.markets",
	})
	acct, err := client.GetAccount()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", *acct)

}
func cryptoDailyClose(period int, ticker string) ([]float64, []float64) {
	var prices []float64
	var volume []float64
	// No keys required for crypto data
	client := marketdata.NewClient(marketdata.ClientOpts{})

	request := marketdata.GetCryptoBarsRequest{
		TimeFrame: marketdata.OneDay,
		Start:     time.Now().AddDate(0, 0, (period*3)*-1),
		End:       time.Now().AddDate(0, 0, -1),
	}
	bars, err := client.GetCryptoBars(ticker, request)
	if err != nil {
		panic(err)
	}
	for _, bar := range bars {
		prices = append(prices, bar.Close)
		volume = append(volume, bar.Volume)
	}
	return prices, volume
}
func stockDailyClose(period int, ticker string) ([]float64, []uint64) {
	var prices []float64
	var volume []uint64
	client := marketdata.NewClient(marketdata.ClientOpts{})
	request := marketdata.GetBarsRequest{
		TimeFrame: marketdata.OneDay,
		Start:     time.Now().AddDate(0, 0, (period*3)*-1),
		End:       time.Now().AddDate(0, 0, -1),
		//AsOf:      "2022-06-10", // Leaving it empty yields the same results
	}
	bars, err := client.GetBars(ticker, request)
	if err != nil {
		panic(err)
	}

	for _, bar := range bars { //edited this so now i can receive data i want, and make it into ints (might want to make an array of ints)
		prices = append(prices, bar.Close)
		volume = append(volume, bar.Volume)

	}
	return prices, volume
}

func getPrices(ticker string, choice string) []float64 {
	option := 1
	if choice == "crypto" {
		option = 2
	}

	switch option {
	case 1:
		prices, volumeStock = stockDailyClose(period, ticker)
		return prices
	case 2:

		prices, volumeCrypto = cryptoDailyClose(period, ticker)
		return prices
	}
	fmt.Println("Invalid option. Please try again")
	return nil
}

func calculateRSI() []float64 { //returns 14 day rsi
	RSI, err := trade_indicators.CalculateRSI(prices, 14)
	if err != nil {
		panic(err.Error())
	}
	return RSI
}

func calculateEMA() []float64 { //returns 26 day ema
	EMA, err := trade_indicators.CalculateEMA(prices, 26)
	if err != nil {
		panic(err.Error())
	}
	return EMA
}

func calculateMACD() []float64 { // MACD = ema(12)-ema(26)
	MACD, _, err := trade_indicators.CalculateMACD(prices)
	if err != nil {
		panic(err.Error())
	}
	return MACD
}
func calculateVolatility() float64 {
	Volatility, err := trade_indicators.CalculateVolatility(prices, period)
	if err != nil {
		panic(err.Error())
	}
	return Volatility
}
func float64SliceToString(floats []float64) string {
	// Create a slice to hold the string representations of each float64
	strFloats := make([]string, len(floats))
	// Convert each float64 to a string and store it in the strFloats slice
	for i, v := range floats {
		strFloats[i] = fmt.Sprintf("%f", v)
	}
	// Join the slice of strings into a single string, separated by commas
	result := strings.Join(strFloats, ", ")
	return result
}

func Ai_request(prompt string) (response string) {

	c := openai.NewClient(OPEN_AI_KEY) //creates a client, which makes the requests to the servers
	// Create the message
	message := openai.ChatCompletionMessage{ //takes the prompt, and creates a message for the request
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	}
	// Create the chat completion request
	resp, err := c.CreateChatCompletion( //this makes the api request
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     "gpt-3.5-turbo",                         //model of ai to use
			MaxTokens: 1024,                                    //# of tokens
			Messages:  []openai.ChatCompletionMessage{message}, //message to send
		},
	)
	if err != nil { //if we get an error, print it
		panic(err)
	}
	// Print the response
	var responseContent strings.Builder //creates a string builder,to create a string of the ai responses
	for _, choice := range resp.Choices {
		responseContent.WriteString(choice.Message.Content) //appends each "choice" word from api to string
	}
	return responseContent.String()

}
func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a new SSM client
	svc := ssm.New(sess)

	// Retrieve the parameter value
	alpaca := "/key/alpaca"
	param, _ := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(alpaca),
		WithDecryption: aws.Bool(true),
	})
	APCA_API_KEY_ID = *param.Parameter.Value

	alpaca_secret := "/secretkey/alpaca"
	secretparam, _ := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(alpaca_secret),
		WithDecryption: aws.Bool(true),
	})
	APCA_API_SECRET_KEY = *secretparam.Parameter.Value

	openaikey := "/key/openai"
	aiparam, _ := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(openaikey),
		WithDecryption: aws.Bool(true),
	})
	OPEN_AI_KEY = *aiparam.Parameter.Value
	os.Setenv("APCA_API_KEY_ID", APCA_API_KEY_ID)
	os.Setenv("APCA_API_SECRET", APCA_API_SECRET_KEY)

}

func main() {
	lambda.Start(handler)
}

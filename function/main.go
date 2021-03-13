package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/jamesjj/line-bot-go/linelambda"
	"github.com/line/line-bot-sdk-go/linebot"
	"os"
)

func main() {
	runtime.Start(handleRequest)
}

func handleRequest(ctx context.Context, lambdaEvent events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Make lambda's event suitable for current Line Bot SDK
	httpReq, hrErr := linelambda.APIEventToHTTPRequest(lambdaEvent)
	if hrErr != nil {
		return linelambda.Finished("Event to HTTP Request", hrErr)
	}

	// Line Bot SDK
	bot, botErr := linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_CHANNEL_ACCESS_TOKEN"),
	)
	if botErr != nil {
		return linelambda.Finished("linebot.New", botErr)
	}

	// Parse web hook data to Line Events
	lineEvents, prErr := bot.ParseRequest(httpReq)
	if prErr != nil {
		return linelambda.Finished("bot.ParseRequest", prErr)
	}

	// Loop over the Line events and handle them
	for _, lineEvent := range lineEvents {

		// Handle Message events
		if lineEvent.Type == linebot.EventTypeMessage {

			// Handle Text messages
			switch message := lineEvent.Message.(type) {
			case *linebot.TextMessage:

				// Prepare some replies to each received message
				textReply1 := linebot.NewTextMessage(
					fmt.Sprintf(
						"You said: \"%s\"\n\nYour UserID is \"%s\"",
						message.Text,
						lineEvent.Source.UserID,
					),
				)
				textReply2 := linebot.NewTextMessage(
					fmt.Sprintf(
						"This is the whole `Message` object:\n%+v\n\nThis is the whole `Source` object:\n%+v\n",
						lineEvent.Message.(linebot.Message),
						lineEvent.Source,
					),
				)

				// Send the prepared reply messages
				_, rmErr := bot.ReplyMessage(lineEvent.ReplyToken, textReply1, textReply2).Do()
				if rmErr != nil {
					fmt.Printf("FAILED: bot.ReplyMessage: %v\n", rmErr)
				}

			}
		}
	}
	return linelambda.Finished("DONE", nil)
}

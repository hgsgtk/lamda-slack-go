package lambdahandler

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/nlopes/slack"
)

func Handler(ctx context.Context, event events.CloudWatchEvent) {
	fmt.Println("Start event handling")

	// FIXME osenvへの依存
	db, err := NewDBConn(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err != nil {
		fmt.Printf("Failed to get connection with database %#v", err)
		return
	}

	// Query database
	userNum, err := CountUser(db)
	if err != nil {
		fmt.Printf("Error by getting user count %#v\n", err)
		return
	}

	totalAmount, err := GetTotalOrder(db)
	if err != nil {
		fmt.Printf("Error by getting total amount %v\n", err)
		return
	}

	// slack
	// FIXME osenvへの依存
	slackChannel := os.Getenv("SLACK_CHANNEL")
	slackApi, err := NewSlackCli(os.Getenv("SLACK_ACCESS_TOKEN"))
	if err != nil {
		fmt.Printf("Error by creating slack client")
		return
	}

	channelID, timestamp, err := slackApi.PostMessage(
		slackChannel,
		slack.MsgOptionText("", false),
		slack.MsgOptionAttachments(slack.Attachment{
			Pretext: "本日のKPI達成状況はこちらです",
			Text:    fmt.Sprintf("新規獲得ユーザー: %d\n累計実績金額: %d", userNum, totalAmount),
		}))
	if err != nil {
		fmt.Printf("Error by sending slack message: %s (detail: %#v)\n", err.Error(), err)
		return
	}
	fmt.Printf("Success to post slack channel %s at %s", channelID, timestamp)
	return
}

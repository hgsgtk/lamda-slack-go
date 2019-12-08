package lambdahandler

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/nlopes/slack"
)

func Handler() error {
	// FIXME osenvへの依存
	db, err := NewDBConn(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err != nil {
		return errors.Wrap(err, "Failed to get connection with database")
	}

	// Query database
	userNum, err := CountUser(db)
	if err != nil {
		return errors.Wrap(err, "Failed by getting user count")
	}

	totalAmount, err := GetTotalOrder(db)
	if err != nil {
		return errors.Wrap(err, "Failed by getting total amount")
	}

	// slack
	// FIXME osenvへの依存
	timeoutStr := os.Getenv("SLACK_API_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return errors.Wrap(err, "Failed by configuration mistake")
	}

	slackApi, err := NewSlackCli(
		os.Getenv("SLACK_ACCESS_TOKEN"), timeout)
	if err != nil {
		return errors.Wrap(err, "Failed by creating slack client")
	}
	slackChannel := os.Getenv("SLACK_CHANNEL")
	_, _, err = slackApi.PostMessage(
		slackChannel,
		slack.MsgOptionText("", false),
		slack.MsgOptionAttachments(slack.Attachment{
			Pretext: "本日のKPI達成状況はこちらです",
			Text:    fmt.Sprintf("新規獲得ユーザー: %d\n累計実績金額: %d", userNum, totalAmount),
		}))
	if err != nil {
		return errors.Wrap(err, "Failed by sending slack message")
	}
	return nil
}

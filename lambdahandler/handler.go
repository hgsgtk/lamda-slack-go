package lambdahandler

import (
	"fmt"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

func Handler() error {
	dbCon, err := GetDBConfig()
	if err != nil {
		return errors.Wrap(err, "failed to get secret")
	}

	db, err := NewDBConn(dbCon)
	if err != nil {
		return errors.Wrap(err, "failed to get connection with database")
	}

	// Query database
	userNum, err := CountUser(db)
	if err != nil {
		return errors.Wrap(err, "failed by getting user count")
	}

	totalAmount, err := GetTotalOrder(db)
	if err != nil {
		return errors.Wrap(err, "failed by getting total amount")
	}

	// slack
	sConf, err := GetSlackConfig()
	if err != nil {
		return errors.Wrap(err, "failed by get secret")
	}

	slackApi, err := NewSlackCli(sConf.AccessToken, sConf.Timeout)
	if err != nil {
		return errors.Wrap(err, "failed by creating slack client")
	}
	_, _, err = slackApi.PostMessage(
		sConf.Channel,
		slack.MsgOptionText("", false),
		slack.MsgOptionAttachments(slack.Attachment{
			Pretext: "本日のKPI達成状況はこちらです",
			Text:    fmt.Sprintf("新規獲得ユーザー: %d\n累計実績金額: %d", userNum, totalAmount),
		}))
	if err != nil {
		return errors.Wrap(err, "failed by sending slack message")
	}
	return nil
}

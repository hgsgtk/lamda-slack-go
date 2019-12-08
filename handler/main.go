package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nlopes/slack"
)

func handler(ctx context.Context, event events.CloudWatchEvent) {
	fmt.Println("Start event handling")

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sqlMode := "TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY"
	loc := "Asia%2FTokyo"

	// Connection database
	ds := fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?sql_mode='%s'&parseTime=true&loc=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbName,
		sqlMode,
		loc)
	db, err := sql.Open("mysql", ds)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error by connection database %#v\n", err)
		return
	}

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "Error by ping database %#v\n", err)
		defer db.Close()
		return
	}
	fmt.Println("Success to connect database")

	// Query database
	// KPI 新規獲得ユーザー数
	q := `SELECT count(*) from users`

	stmt, err := db.Prepare(q)
	if err != nil {
		fmt.Printf("Error by preparing statement %#v\n", err)
		return
	}

	var userNum int
	if err := stmt.QueryRow().Scan(&userNum); err != nil {
		fmt.Printf("Error by quering and scanning database %#v\n", err)
		return
	}
	// KPI 累計実績金額
	q = `SELECT sum(amount) from orders`

	stmt, err = db.Prepare(q)
	if err != nil {
		fmt.Printf("Error by preparing statement %#v\n", err)
		return
	}

	var totalAmount int
	if err := stmt.QueryRow().Scan(&totalAmount); err != nil {
		fmt.Printf("Error by quering and scanning database %#v\n", err)
		return
	}

	// slack
	accessToken := os.Getenv("SLACK_ACCESS_TOKEN")
	slackChannel := os.Getenv("SLACK_CHANNEL")
	httpClient := http.Client{
		// タイムアウト設定をしておく
		Timeout: time.Duration(time.Second * 3),
	}
	slackApi := slack.New(accessToken, slack.OptionHTTPClient(&httpClient))

	if _, err := slackApi.AuthTest(); err != nil {
		fmt.Printf("Error by getting auth; %s detail %#v", err.Error(), err)
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

func main() {
	lambda.Start(handler)
}

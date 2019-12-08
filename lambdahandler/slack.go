package lambdahandler

import (
	"net/http"
	"time"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

func NewSlackCli(token string, timeout int) (*slack.Client, error) {
	httpClient := http.Client{
		// lambdaが死ぬ前にタイムアウトするように設定をしておく
		Timeout: time.Duration(time.Second * time.Duration(timeout)),
	}
	cli := slack.New(token, slack.OptionHTTPClient(&httpClient))

	if _, err := cli.AuthTest(); err != nil {
		return nil, errors.Wrap(err, "auth error")
	}

	return cli, nil
}

package monitor

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/lib/constant"
	"gitlab.senseauto.com/apcloud/app/datacollector-service/internal/model/dao"
	cmclient "gitlab.senseauto.com/apcloud/library/common-go/client"
	cmhttp "gitlab.senseauto.com/apcloud/library/common-go/client/http"
	cmconfig "gitlab.senseauto.com/apcloud/library/common-go/config"
	cmlog "gitlab.senseauto.com/apcloud/library/common-go/log"
)

const (
	keyKey       = "key"
	typeMarkdown = "markdown"
)

type AlertType int

const (
	AlertService AlertType = iota
	AlertImportData
)

var alertStr = []string{"AlertService", "AlertImportData"}

func (alert AlertType) String() string {
	return alertStr[alert]
}

func (alert AlertType) Index() int {
	return int(alert)
}

type AlertConfig struct {
	Path      string
	Key       string
	Templates []string
}

type AlertMsg struct {
	AlertType AlertType
	Service   string
	Module    string
	Time      string
	Msg       string
	Plus      map[string]string
}

type Alert struct {
	config AlertConfig
	ctx    context.Context
}

func (a *Alert) init(ctx context.Context) {
	logger := cmlog.Extract(ctx)

	if err := cmconfig.Global().UnmarshalKey(constant.CfgAlert, &a.config); err != nil {
		logger.Panic(err)
	}

	alertHTTPCfg := cmhttp.Config{}

	if err := cmconfig.Global().UnmarshalKey(constant.CfgHTTPAlert, &alertHTTPCfg); err != nil {
		logger.Panic(err)
	}

	cmclient.HTTP.Global(ctx, map[string]cmhttp.Config{
		constant.CfgHTTPAlertKey: alertHTTPCfg,
	})
}

func genServiceMsg(template string, msg *AlertMsg) string {
	result := fmt.Sprintf(template, msg.Service, msg.Module, msg.Time, msg.Msg)
	return result
}

func genImportMsg(template string, msg *AlertMsg) string {
	result := fmt.Sprintf(template, msg.Service, msg.Module, msg.Time, msg.Msg)

	for key, value := range msg.Plus {
		item := fmt.Sprintf(constant.FmtAlertPlus, key, value)
		result += item
	}

	return result
}

func (a *Alert) genAlertMsg(msg *AlertMsg) *dao.AlertRequest {
	request := dao.AlertRequest{
		Msgtype:  typeMarkdown,
		Markdown: dao.AlertMarkdown{}}

	switch msg.AlertType {
	case AlertService:
		request.Markdown.Content = genServiceMsg(a.config.Templates[AlertService], msg)
	case AlertImportData:
		request.Markdown.Content = genImportMsg(a.config.Templates[AlertImportData], msg)
	default:
	}

	return &request
}

func (a *Alert) sendAlert(msgs []*AlertMsg) {
	logger := cmlog.Extract(a.ctx)
	cli := cmclient.HTTP.GetClient(constant.CfgHTTPAlertKey)

	for _, msg := range msgs {
		req := a.genAlertMsg(msg)
		request := cli.R()
		resp, err := request.
			SetHeader("Content-Type", "application/json;charset=UTF-8").
			SetQueryParam(keyKey, a.config.Key).
			SetBody(*req).
			Post(a.config.Path)

		if err != nil {
			logger.Error("Send alert failed!")
			logger.Error(err)
		} else {
			respBody := &dao.AlertReponse{}
			if err := json.Unmarshal(resp.Body(), respBody); err != nil {
				logger.Errorf("unmarshal response body failed: %v", err)
			} else if respBody.Errcode != 0 {
				logger.Errorf("Send alert failed. msg : %s", respBody.Errmsg)
			}
		}
	}
}

func newAlert(ctx context.Context) *Alert {
	return &Alert{
		ctx:    ctx,
		config: AlertConfig{},
	}
}

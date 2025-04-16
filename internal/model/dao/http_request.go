package dao

type AlertReponse struct {
	Errcode int32  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type AlertMarkdown struct {
	Content string `json:"content"`
}

type AlertRequest struct {
	Msgtype  string        `json:"msgtype"`
	Markdown AlertMarkdown `json:"markdown"`
}

type Response[T any] struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type DataCollectorResponse struct {
}

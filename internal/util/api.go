package util

type APIResponse[T any] struct {
	Data       T      `json:"data,omitempty"`
	StatusCode int    `json:"code"` // 自定义错误码,需要判定特殊场景下使用. 大部分错误优先使用msg+http 异常错误码返回
	Msg        string `json:"msg"`
}

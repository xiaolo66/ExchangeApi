package ExchangeApi

import "fmt"

type ExError struct {
	Code    int                    // error code
	Message string                 // error message
	Data    map[string]interface{} // business data
}

func (e ExError) Error() string {
	return fmt.Sprintf("code: %v message: %v", e.Code, e.Message)
}

const (
	NotImplement = 10000 + iota
	UnHandleError //未被捕获的错误

	//exchange api business error
	ErrExchangeSystem = 20000 + iota //交易所系统出现错误
	ErrDataParse //解析错误
	ErrAuthFailed //鉴权错误
	ErrRequestParams
	ErrInsufficientFunds //资产不足错误
	ErrInvalidOrder //无效的订单
	ErrInvalidAddress //无效的地址
	ErrOrderNotFound //未发现订单
	ErrNotFoundMarket
	ErrChannelNotExist
	ErrInvalidDepth // recv dirty order book data

	//network error
	ErrDDoSProtection = 30000 + iota //发生限频错误
	ErrTimeout
	ErrBadRequest
	ErrBadResponse
)
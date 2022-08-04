package utils

import (
	"encoding/json"
	"strings"

	"rpc-load-test/utils/logutil"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gorilla/websocket"
)

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type JsonRpcMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Error   *jsonError      `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

type ErrMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ToMsgParams(msg ethereum.CallMsg) interface{} {
	params := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		params["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		params["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		params["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		params["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return params
}

func NewMsg(method string, params ...interface{}) (*JsonRpcMessage, error) {
	msg := &JsonRpcMessage{Version: "2.0", ID: []byte("1"), Method: method}
	if params != nil { // prevent sending "params" as null value
		var err error
		if msg.Params, err = json.Marshal(params); err != nil {
			return nil, err
		}
	}
	return msg, nil
}

func SendMsg(host string, msg *JsonRpcMessage, timeout int64, check string) (*Response, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	url := strings.Split(host, "://")
	req := Request{
		msg.Method,
		"POST",
		url[0],
		url[1],
		"",
		string(body),
		headers,
		timeout,
		check,
	}
	res, err := req.Call()
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func SendBatchMsg(host string, msg []*JsonRpcMessage, timeout int64, check string) (*Response, error) {
	body, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	//
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	url := strings.Split(host, "://")
	req := Request{
		"batch",
		"POST",
		url[0],
		url[1],
		"",
		string(body),
		headers,
		timeout,
		check,
	}
	res, err := req.Call()
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func SendWssMsg(host string, req []*JsonRpcMessage) (*Response, error) {
	webClient, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		return nil, err
	}
	defer webClient.Close()
	//
	err = webClient.WriteJSON(req)
	if err != nil {
		return nil, err
	}
	//
	for {
		_, message, err := webClient.ReadMessage()
		if err != nil {
			logutil.Error("Read:", err)
		}
		logutil.Debug("Recv:", message)
		return &Response{Body: string(message)}, nil
	}
}

func ToFilterParams(q ethereum.FilterQuery) (interface{}, error) {
	arg := map[string]interface{}{
		"address": q.Addresses,
		"topics":  q.Topics,
	}
	for _, topic := range q.Topics {
		if topic == nil {
			arg["topics"] = nil
			break
		}
	}
	if q.BlockHash != nil {
		arg["blockHash"] = *q.BlockHash
	}
	if q.FromBlock == nil {
		arg["fromBlock"] = nil
	} else {
		arg["fromBlock"] = hexutil.EncodeBig(q.FromBlock)
	}
	if q.ToBlock == nil {
		arg["toBlock"] = nil
	} else {
		arg["toBlock"] = hexutil.EncodeBig(q.ToBlock)
	}
	return arg, nil
}

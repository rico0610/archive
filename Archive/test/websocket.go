package test

import (
	"math/big"
	"math/rand"
	"time"

	"rpc-load-test/utils"
	"rpc-load-test/utils/logutil"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gorilla/websocket"
)

// Websocket_eth_subscribe_newHeads  use res record result, res length is subscribed count, use Response.Time record  subscribe return time
func Websocket_eth_subscribe_newHeads(cfg *utils.Config, method string, timeOut time.Time) ([]*utils.Response, error) {
	r := &utils.Response{Name: method, Code: 200}
	var res []*utils.Response
	ws := cfg.WsHost[rand.Intn(len(cfg.WsHost))]
	webClient, _, err := websocket.DefaultDialer.Dial(ws, cfg.Headers)
	if err != nil {
		logutil.Errorf("dial err: %v", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	defer webClient.Close()
	msg, err := utils.NewMsg("eth_subscribe", "newHeads")
	if err != nil {
		logutil.Error("NewMsg:", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	return subscribe(webClient, msg, res, method, 7, timeOut)

}

// Websocket_eth_subscribe_logs  use res record result, res length is subscribed count, use Response.Time record  subscribe return time
func Websocket_eth_subscribe_logs(cfg *utils.Config, method string, address []common.Address, topic0, topic1, topic2, topic3 []common.Hash, timeOut time.Time) ([]*utils.Response, error) {
	r := &utils.Response{Name: method, Code: 200}
	var res []*utils.Response
	webClient, _, err := websocket.DefaultDialer.Dial(cfg.WsHost[rand.Intn(len(cfg.WsHost))], cfg.Headers)
	if err != nil {
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	filterQuery := ethereum.FilterQuery{
		Addresses: address,
		Topics:    [][]common.Hash{topic0, topic1, topic2, topic3},
	}
	params, err := utils.ToFilterParams(filterQuery)
	if err != nil {
		logutil.Error("ToFilterParams:", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	defer webClient.Close()
	msg, err := utils.NewMsg("eth_subscribe", "logs", params)
	if err != nil {
		logutil.Error("NewMsg:", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	return subscribe(webClient, msg, res, method, 7, timeOut)
}

// Websocket_eth_subscribe_newPendingTransactions  use res record result, res length is subscribe return count, use Response.Time record subscribe return time
func Websocket_eth_subscribe_newPendingTransactions(cfg *utils.Config, method string, timeOut time.Time) ([]*utils.Response, error) {
	r := &utils.Response{Name: method, Code: 200}
	var res []*utils.Response
	webClient, _, err := websocket.DefaultDialer.Dial(cfg.WsHost[rand.Intn(len(cfg.WsHost))], cfg.Headers)
	if err != nil {
		logutil.Errorf("dial: %v", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	defer webClient.Close()
	msg, err := utils.NewMsg("eth_subscribe", "newPendingTransactions")
	if err != nil {
		logutil.Error("NewMsg:", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	return subscribe(webClient, msg, res, method, 30, timeOut)
}

func subscribe(webClient *websocket.Conn, msg *utils.JsonRpcMessage, res []*utils.Response, method string, readDeadline int, timeOut time.Time) ([]*utils.Response, error) {
	r := &utils.Response{Name: method, Code: 200}
	err := webClient.WriteJSON(msg)
	if err != nil {
		logutil.Errorf("WriteJSON err: %v", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	interval := time.Duration(readDeadline) * time.Second
	err = webClient.SetReadDeadline(time.Now().Add(interval))
	if err != nil {
		logutil.Errorf("SetReadDeadline err: %v", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	preTime := time.Now()
	for {
		_, message, err := webClient.ReadMessage()
		if err != nil {
			logutil.Errorf("read msg err:%v", err)
			r.Code = -1
			res = append(res, r)
			return res, err
		}
		err = webClient.SetReadDeadline(time.Now().Add(interval))
		if err != nil {
			logutil.Errorf("SetReadDeadline err: %v", err)
			r.Code = -1
			res = append(res, r)
			return res, err
		}
		logutil.Tracef("Received: %s.\n", message)
		// record time
		currTime := time.Now()
		duration := currTime.Sub(preTime).Milliseconds()
		preTime = currTime
		r.Time = duration
		res = append(res, r)
		if time.Now().Sub(timeOut) >= 0 {
			logutil.Debugf("===%s Subscription time reached ===", method)
			return res, nil
		}
		logutil.Debugf("===subscribe %s count: %d===", method, len(res))
	}
}

func subscribeAndSend(webClient *websocket.Conn, msg []*utils.JsonRpcMessage, nums int, res []*utils.Response, method string, readDeadline int, timeOut time.Time) ([]*utils.Response, error) {
	r := &utils.Response{Name: method, Code: 200}

	count := nums * len(msg)
	pretimes := make([]time.Time, count)
	for i := 0; i < count; i++ {
		index := i % len(msg)
		err := webClient.WriteJSON(msg[index])
		pretimes[i] = time.Now()
		if err != nil {
			logutil.Errorf("WriteJSON err: %v", err)
			r.Code = -1
			res = append(res, r)
			return res, err
		}
	}

	interval := time.Duration(readDeadline) * time.Second
	err := webClient.SetReadDeadline(time.Now().Add(interval))
	if err != nil {
		logutil.Errorf("SetReadDeadline err: %v", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	c := 0
	for {
		_, message, err := webClient.ReadMessage()
		if err != nil {
			logutil.Errorf("read msg err:%v", err)
			r.Code = -1
			res = append(res, r)
			return res, err
		}
		err = webClient.SetReadDeadline(time.Now().Add(interval))
		if err != nil {
			logutil.Errorf("SetReadDeadline err: %v", err)
			r.Code = -1
			res = append(res, r)
			return res, err
		}
		logutil.Tracef("Received: %s.\n", message)
		// record time
		currTime := time.Now()
		duration := currTime.Sub(pretimes[c]).Milliseconds()
		c++
		r.Time = duration
		res = append(res, r)
		if time.Now().Sub(timeOut) >= 0 || c >= count-1 {
			logutil.Debugf("===%s Subscription time reached ===", method)
			return res, nil
		}
		logutil.Debugf("===subscribe %s count: %d===", method, len(res))
	}
}

func Websocket_eth_subscribe_multi(cfg *utils.Config, method string, timeOut time.Time) ([]*utils.Response, error) {
	r := &utils.Response{Name: method, Code: 200}
	var res []*utils.Response
	webClient, _, err := websocket.DefaultDialer.Dial(cfg.WsHost[rand.Intn(len(cfg.WsHost))], cfg.Headers)
	if err != nil {
		logutil.Errorf("dial: %v", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	defer webClient.Close()
	rt := rand.Intn(len(cfg.TXHashList))
	txhash := common.HexToHash(cfg.TXHashList[rt])
	msg1, err := utils.NewMsg("eth_getTransactionReceipt", txhash)
	if err != nil {
		logutil.Error("NewMsg:", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	rh := rand.Intn(int(cfg.Range))
	height := cfg.Height - int64(rh)
	rp := rand.Intn(len(cfg.Bep20Addrs))
	path := []common.Address{cfg.WBNBAddr, cfg.Bep20Addrs[rp]}
	data, err := cfg.V2routerABI_.Pack("getAmountsOut", big.NewInt(10000), path)
	if err != nil {
		return nil, err
	}
	callMsg := ethereum.CallMsg{common.Address{}, &cfg.RouterAddr, uint64(0), nil, nil, nil, nil, data, nil}
	msg2, err := utils.NewMsg(ethCall, utils.ToMsgParams(callMsg), hexutil.EncodeBig(big.NewInt(height)))
	if err != nil {
		logutil.Error("NewMsg:", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	msg3, err := utils.NewMsg("eth_blockNumber")
	if err != nil {
		logutil.Error("NewMsg:", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	blockRange := rand.Intn(int(cfg.Range))
	ra := rand.Intn(len(cfg.Bep20Addrs) - 1)
	filterQuery := ethereum.FilterQuery{
		FromBlock: big.NewInt(cfg.Height - int64(blockRange)),
		ToBlock:   big.NewInt(cfg.Height),
		Addresses: cfg.Bep20Addrs[ra : ra+1],
		Topics:    [][]common.Hash{},
	}
	params, err := utils.ToFilterParams(filterQuery)
	if err != nil {
		return nil, err
	}
	msg4, err := utils.NewMsg("eth_getLogs", params)
	if err != nil {
		logutil.Error("NewMsg:", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	heights := cfg.Height - int64(rand.Intn(int(cfg.Range)))
	msg5, err := utils.NewMsg("eth_getBlockByNumber", (*hexutil.Big)(big.NewInt(heights)), true)
	if err != nil {
		logutil.Error("NewMsg:", err)
		r.Code = -1
		res = append(res, r)
		return res, err
	}
	msg := []*utils.JsonRpcMessage{msg1, msg2, msg3, msg4, msg5}
	return subscribeAndSend(webClient, msg, 100, res, method, 30, timeOut)
}

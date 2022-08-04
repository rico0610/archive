package test

import (
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"rpc-load-test/utils"
	"rpc-load-test/utils/logutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/gjson"
)

func Run_RPC_eth_call_getAmountsOut(cfg *utils.Config, method string) []*utils.Response {
	//
	r := rand.Intn(int(cfg.Range))
	height := cfg.Height - int64(r)
	r = rand.Intn(len(cfg.Bep20Addrs))
	path := []common.Address{cfg.WBNBAddr, cfg.Bep20Addrs[r]}
	res, err := RPC_eth_call_getAmountsOut(cfg, method, height, path, 10)
	if err != nil {
		logutil.Errorf("%s-getAmountsOut: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_call_balanceOf(cfg *utils.Config, method string) []*utils.Response {
	//
	r := rand.Intn(int(cfg.Range))
	height := cfg.Height - int64(r)
	r = rand.Intn(len(cfg.AccountList))
	bep20 := utils.T_cfg.BEP20[r]
	extAcc := cfg.AccountList[bep20][rand.Intn(len(cfg.AccountList[bep20]))]
	acc := common.HexToAddress(extAcc)
	c := common.HexToAddress(bep20)
	res, err := RPC_eth_call_balanceOf(cfg, method, height, acc, c, 10)
	if err != nil {
		logutil.Errorf("%s-balanceOf: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_getStorageAt(cfg *utils.Config, method string) []*utils.Response {
	//
	r := rand.Intn(int(cfg.Range))
	height := cfg.Height - int64(r)
	r = rand.Intn(len(cfg.BEP20))
	bep20 := cfg.BEP20[r]
	res, err := RPC_eth_getStorageAt(cfg, method, height, bep20, 10)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_getBalance(cfg *utils.Config, method string) []*utils.Response {
	r := rand.Intn(int(cfg.Range))
	height := cfg.Height - int64(r)
	r = rand.Intn(len(cfg.AccountList))
	bep20 := utils.T_cfg.BEP20[r]
	extAcc := cfg.AccountList[bep20][rand.Intn(len(cfg.AccountList[bep20]))]
	res, err := RPC_eth_getBalance(cfg, method, height, extAcc, 10)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_getTransactionCount(cfg *utils.Config, method string) []*utils.Response {
	r := rand.Intn(int(cfg.Range))
	height := cfg.Height - int64(r)
	r = rand.Intn(len(cfg.AccountList))
	bep20 := utils.T_cfg.BEP20[r]
	extAcc := cfg.AccountList[bep20][rand.Intn(len(cfg.AccountList[bep20]))]
	res, err := RPC_eth_getTransactionCount(cfg, method, height, extAcc, 10)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_getCode(cfg *utils.Config, method string) []*utils.Response {
	r := rand.Intn(int(cfg.Range))
	height := cfg.Height - int64(r)
	r = rand.Intn(len(cfg.BEP20))
	bep20 := cfg.BEP20[r]
	res, err := RPC_eth_getCode(cfg, method, height, bep20, 10)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_getProof(cfg *utils.Config, method string) []*utils.Response {
	r := rand.Intn(int(cfg.Range))
	height := cfg.Height - int64(r)
	r = rand.Intn(len(cfg.BEP20))
	bep20 := cfg.BEP20[r]
	res, err := RPC_eth_getProof(cfg, method, height, bep20, []string{cfg.BEP20StorageAt}, 10)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_estimateGas(cfg *utils.Config, method string) []*utils.Response {
	//
	r := rand.Intn(int(cfg.Range))
	height := cfg.Height - int64(r)
	r = rand.Intn(len(cfg.Bep20Addrs))
	addr := cfg.Bep20Addrs[r]
	res, err := RPC_eth_estimateGas(cfg, method, height, addr, cfg.RouterAddr, 10)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_debug_traceBlockByNumber(cfg *utils.Config, method string) []*utils.Response {
	//
	r := rand.Intn(int(cfg.Range))
	height := cfg.Height - int64(r)
	res, err := RPC_debug_traceBlockByNumber(cfg, method, height, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_debug_traceTransaction(cfg *utils.Config, method string) []*utils.Response {
	//
	r := rand.Intn(len(cfg.TXHashList))
	txhash := common.HexToHash(cfg.TXHashList[r])
	res, err := RPC_debug_traceTransaction(cfg, method, txhash, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_debug_traceCall(cfg *utils.Config, method string) []*utils.Response {
	//
	r := rand.Intn(int(cfg.Range))
	height := cfg.Height - int64(r)
	r = rand.Intn(len(cfg.Bep20Addrs))
	bep20Addr := cfg.Bep20Addrs[r]
	r = rand.Intn(len(cfg.AccountList))

	bep20 := utils.T_cfg.BEP20[r]
	extAcc := cfg.AccountList[bep20][rand.Intn(len(cfg.AccountList[bep20]))]
	extAccAddr := common.HexToAddress(extAcc)
	res, err := RPC_debug_traceCall(cfg, method, height, bep20Addr, extAccAddr, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

// Run_RPC_eth_getLogs use block range and contractAddress for test
func Run_RPC_eth_getLogs(cfg *utils.Config, method string) []*utils.Response {
	blockRange := rand.Intn(int(cfg.Range))
	res, err := RPC_eth_getLogs(cfg, method, big.NewInt(int64(blockRange)), big.NewInt(int64(blockRange)), nil, nil, nil, nil, nil, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_net_version(cfg *utils.Config, method string) []*utils.Response {
	res, err := RPC_net_version(cfg, method, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_blockNumber(cfg *utils.Config, method string) []*utils.Response {
	res, err := RPC_eth_blockNumber(cfg, method, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_NewFilter(cfg *utils.Config, method string) []*utils.Response {
	res, err := RPC_eth_newFilter(cfg, method, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_newBlockFilter(cfg *utils.Config, method string) []*utils.Response {
	res, err := RPC_eth_newBlockFilter(cfg, method, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_newPendingTransactionFilter(cfg *utils.Config, method string) []*utils.Response {
	res, err := RPC_eth_newPendingTransactionFilter(cfg, method, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_getFilterLogs(cfg *utils.Config, method string) []*utils.Response {
	r := rand.Intn(len(cfg.FilterLogList))
	filterid := cfg.FilterLogList[r]
	res, err := RPC_eth_getFilterLogs(cfg, method, filterid, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_getFilterLogsOnUpdate(cfg *utils.Config, method string) []*utils.Response {
	arg := map[string]interface{}{}
	msgf, err := utils.NewMsg("eth_newFilter", arg)
	if err != nil {
		return nil
	}
	resf, err := utils.SendMsg(cfg.FullNode, msgf, 300, "")
	if err != nil {
		return nil
	}
	Filters := gjson.Get(resf.Body, "result").String()
	res, err := RPC_eth_getFilterLogs(cfg, method, Filters, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_getFilterChanges(cfg *utils.Config, method string) []*utils.Response {
	r := rand.Intn(len(cfg.FilterList))
	filterid := cfg.FilterList[r]
	res, err := RPC_eth_getFilterLogs(cfg, method, filterid, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_getFilterChangesOnupdate(cfg *utils.Config, method string) []*utils.Response {
	msgf, err := utils.NewMsg("eth_newBlockFilter")
	if err != nil {
		return nil
	}
	resf, err := utils.SendMsg(cfg.FullNode, msgf, 300, "")
	if err != nil {
		return nil
	}
	Filters := gjson.Get(resf.Body, "result").String()
	time.Sleep(2 * time.Second)
	res, err := RPC_eth_getFilterLogs(cfg, method, Filters, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_uninstallFilter(cfg *utils.Config, method string) []*utils.Response {
	r := rand.Intn(len(cfg.FilterUninstallList))
	filterid := cfg.FilterUninstallList[r]
	res, err := RPC_eth_getFilterLogs(cfg, method, filterid, 300)
	if err != nil {
		//logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_uninstallFilterOnupdate(cfg *utils.Config, method string) []*utils.Response {
	msgf, err := utils.NewMsg("eth_newBlockFilter")
	if err != nil {
		return nil
	}
	resf, err := utils.SendMsg(cfg.FullNode, msgf, 300, "")
	if err != nil {
		return nil
	}
	Filters := gjson.Get(resf.Body, "result").String()
	res, err := RPC_eth_getFilterLogs(cfg, method, Filters, 300)
	if err != nil {
		//logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_getTransactionReceipt(cfg *utils.Config, method string) []*utils.Response {
	r := rand.Intn(len(cfg.TXHashList))
	txhash := common.HexToHash(cfg.TXHashList[r])
	res, err := RPC_eth_getTransactionReceipt(cfg, method, txhash, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_getBlockByNumber(cfg *utils.Config, method string) []*utils.Response {
	r := rand.Intn(int(cfg.Range))
	height := cfg.Height - int64(r)
	res, err := RPC_eth_getBlockByNumber(cfg, method, height, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_chainId(cfg *utils.Config, method string) []*utils.Response {
	res, err := RPC_eth_chainId(cfg, method, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_getTransactionByHash(cfg *utils.Config, method string) []*utils.Response {
	r := rand.Intn(len(cfg.TXHashList))
	txhash := common.HexToHash(cfg.TXHashList[r])
	res, err := RPC_eth_getTransactionByHash(cfg, method, txhash, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}
func Run_Websocket_eth_subscribe_newHeads(cfg *utils.Config, method string) []*utils.Response {

	dur := utils.T_cfg.StartTime.Add(time.Duration(cfg.DurationSec-1) * time.Second)

	res, err := Websocket_eth_subscribe_newHeads(cfg, method, dur)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
	}
	return res
}

func Run_Websocket_eth_subscribe_newPendingTransactions(cfg *utils.Config, method string) []*utils.Response {

	dur := utils.T_cfg.StartTime.Add(time.Duration(cfg.DurationSec-1) * time.Second)
	res, err := Websocket_eth_subscribe_newPendingTransactions(cfg, method, dur)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
	}
	return res
}

func Run_Websocket_eth_subscribe_logs(cfg *utils.Config, method string) []*utils.Response {
	//r := rand.Intn(len(cfg.Bep20Addrs) - 2)
	dur := utils.T_cfg.StartTime.Add(time.Duration(cfg.DurationSec-1) * time.Second)
	res, err := Websocket_eth_subscribe_logs(cfg, method, nil, nil, nil, nil, nil, dur)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
	}
	return res
}

func Run_Websocket_eth_subscribe_multi(cfg *utils.Config, method string) []*utils.Response {

	dur := utils.T_cfg.StartTime.Add(time.Duration(cfg.DurationSec-1) * time.Second)

	res, err := Websocket_eth_subscribe_multi(cfg, method, dur)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
	}
	return res
}

func Run_RPC_bor_getAuthor(cfg *utils.Config, method string) []*utils.Response {
	r := rand.Intn(int(cfg.Range)) + 1
	res, err := RPC_bor_getAuthor(cfg, method, r, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_bor_getSingalByHash(cfg *utils.Config, method string) []*utils.Response {
	r := rand.Intn(len(cfg.BlockHashList))
	txhash := common.HexToHash(cfg.BlockHashList[r])
	res, err := RPC_bor_getSingalByHash(cfg, method, txhash, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_Bor_getCurrent(cfg *utils.Config, method string) []*utils.Response {
	res, err := RPC_bor_getCurrent(cfg, method, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_Bor_getRootHash(cfg *utils.Config, method string) []*utils.Response {
	res, err := RPC_bor_getRootHash(cfg, method, cfg.Height, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_eth_getTransactionReceiptsByBlock(cfg *utils.Config, method string) []*utils.Response {
	r := rand.Intn(int(cfg.Range))
	res, err := RPC_eth_getTransactionReceiptsByBlock(cfg, method, strconv.Itoa(r), 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}
func Run_RPC_btach_top3(cfg *utils.Config, method string) []*utils.Response {
	res, err := RPC_btach_top3(cfg, method, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

func Run_RPC_GetTransactionDataAndReceipt(cfg *utils.Config, method string) []*utils.Response {
	//
	r := rand.Intn(len(cfg.TXHashList))
	txhash := common.HexToHash(cfg.TXHashList[r])
	res, err := RPC_eth_GetTransactionDataAndReceipt(cfg, method, txhash, 300)
	if err != nil {
		logutil.Errorf("%s: %s", method, err.Error())
		res = &utils.Response{method, -1, nil, "", 0}
	}
	return []*utils.Response{res}
}

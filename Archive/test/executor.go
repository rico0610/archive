package test

import (
	"math/rand"

	"rpc-load-test/utils"
)

func Run(cfg *utils.Config) []*utils.Response {
	//
	selected := ""
	var total int
	for _, v := range cfg.Scenarios {
		total += v
	}
	//
	r := rand.Intn(total)
	for k, v := range cfg.Scenarios {
		r -= v
		if r <= 0 {
			selected = k
			break
		}
	}
	//
	switch selected {
	case "eth_call_getAmountsOut":
		{
			return Run_RPC_eth_call_getAmountsOut(cfg, selected)
		}
	case "eth_call_balanceOf":
		{
			return Run_RPC_eth_call_balanceOf(cfg, selected)
		}
	case "eth_getStorageAt":
		{
			return Run_RPC_eth_getStorageAt(cfg, selected)
		}
	case "eth_getBalance":
		{
			return Run_RPC_eth_getBalance(cfg, selected)
		}
	case "eth_getTransactionCount":
		{
			return Run_RPC_eth_getTransactionCount(cfg, selected)
		}
	case "eth_getCode":
		{
			return Run_RPC_eth_getCode(cfg, selected)
		}
	case "eth_getProof":
		{
			return Run_RPC_eth_getProof(cfg, selected)
		}
	case "eth_estimateGas":
		{
			return Run_RPC_eth_estimateGas(cfg, selected)
		}
	case "debug_traceBlockByNumber":
		{
			return Run_RPC_debug_traceBlockByNumber(cfg, selected)
		}
	case "debug_traceTransaction":
		{
			return Run_RPC_debug_traceTransaction(cfg, selected)
		}
	case "debug_traceCall":
		{
			return Run_RPC_debug_traceCall(cfg, selected)
		}
	case "eth_getLogs":
		{
			return Run_RPC_eth_getLogs(cfg, selected)
		}
	case "net_version":
		{
			return Run_RPC_net_version(cfg, selected)
		}
	case "eth_blockNumber":
		{
			return Run_RPC_eth_blockNumber(cfg, selected)
		}
	case "eth_newFilter":
		{
			return Run_RPC_eth_NewFilter(cfg, selected)
		}
	case "eth_newBlockFilter":
		{
			return Run_RPC_eth_newBlockFilter(cfg, selected)
		}
	case "eth_newPendingTransactionFilter":
		{
			return Run_RPC_eth_newPendingTransactionFilter(cfg, selected)
		}
	case "eth_getFilterLogs":
		{
			return Run_RPC_eth_getFilterLogs(cfg, selected)
		}
	case "eth_getFilterLogsOnupdate":
		{
			return Run_RPC_eth_getFilterLogsOnUpdate(cfg, "eth_getFilterLogs")
		}
	case "eth_getFilterChanges":
		{
			return Run_RPC_getFilterChanges(cfg, selected)
		}
	case "eth_getFilterChangesOnupdate":
		{
			return Run_RPC_getFilterChangesOnupdate(cfg, "eth_getFilterChanges")
		}
	case "eth_uninstallFilter":
		{
			return Run_RPC_eth_uninstallFilter(cfg, selected)
		}
	case "eth_uninstallFilterOnupdate":
		{
			return Run_RPC_eth_uninstallFilterOnupdate(cfg, "eth_uninstallFilter")
		}
	case "eth_getTransactionReceipt":
		{
			return Run_RPC_eth_getTransactionReceipt(cfg, selected)
		}
	case "eth_getBlockByNumber":
		{
			return Run_RPC_eth_getBlockByNumber(cfg, selected)
		}
	case "eth_chainId":
		{
			return Run_RPC_eth_chainId(cfg, selected)
		}
	case "eth_getTransactionByHash":
		{
			return Run_RPC_eth_getTransactionByHash(cfg, selected)
		}
	case "eth_subscribe_newHeads":
		{
			return Run_Websocket_eth_subscribe_newHeads(cfg, selected)
		}
	case "eth_subscribe_logs":
		{
			return Run_Websocket_eth_subscribe_logs(cfg, selected)
		}
	case "eth_subscribe_newPendingTransactions":
		{
			return Run_Websocket_eth_subscribe_newPendingTransactions(cfg, selected)
		}
	case "eth_subscribe_multiple":
		{
			return Run_Websocket_eth_subscribe_multi(cfg, selected)
		}
	case "bor_getSignersAtHash":
		{
			return Run_RPC_bor_getSingalByHash(cfg, selected)
		}
	case "bor_getAuthor":
		{
			return Run_RPC_bor_getAuthor(cfg, selected)
		}
	case "bor_getCurrentValidators":
		{
			return Run_RPC_Bor_getCurrent(cfg, selected)
		}
	case "bor_getCurrentProposer":
		{
			return Run_RPC_Bor_getCurrent(cfg, selected)
		}
	case "bor_getRootHash":
		{
			return Run_RPC_Bor_getRootHash(cfg, selected)
		}
	case "eth_getTransactionReceiptsByBlock":
		{
			return Run_RPC_eth_getTransactionReceiptsByBlock(cfg, selected)
		}
	case "batch":
		{
			return Run_RPC_btach_top3(cfg, selected)
		}
	case "eth_getTransactionDataAndReceipt":
		{
			return Run_RPC_GetTransactionDataAndReceipt(cfg, selected)
		}
	case "nr_getTransactionDataAndReceipt":
		{
			return Run_RPC_GetTransactionDataAndReceipt(cfg, selected)
		}
	default:
		{
			panic("invalid scenario " + selected)
		}
	}
	return nil
}

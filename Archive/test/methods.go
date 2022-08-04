package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"

	"rpc-load-test/contracts/bep20"
	"rpc-load-test/utils"
	"rpc-load-test/utils/logutil"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const ethCall = "eth_call"

func RPC_eth_call_balanceOf(cfg *utils.Config, method string, height int64, account, contract common.Address, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s-balanceOf,height: %d, account: %s, contract: %s", method, height, account, contract)

	bep20ABI, err := abi.JSON(strings.NewReader(bep20.Bep20ABI))
	if err != nil {
		return nil, err
	}
	data, err := bep20ABI.Pack("balanceOf", account)
	if err != nil {
		return nil, err
	}
	callMsg := ethereum.CallMsg{
		From:     account,
		To:       &contract,
		Gas:      uint64(0),
		GasPrice: nil,
		Value:    nil,
		Data:     data,
	}
	msg, err := utils.NewMsg(ethCall, utils.ToMsgParams(callMsg), hexutil.EncodeBig(big.NewInt(height)))
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "0x0000")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	// change res method name
	res.Name = method
	return res, nil
}

func RPC_eth_call_getAmountsOut(cfg *utils.Config, method string, height int64, path []common.Address, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s-getAmountsOut: %d, %s, %s", method, height, path[0].Hex(), path[1].Hex())
	//
	data, err := cfg.V2routerABI_.Pack("getAmountsOut", big.NewInt(10000), path)
	if err != nil {
		return nil, err
	}
	callMsg := ethereum.CallMsg{common.Address{}, &cfg.RouterAddr, uint64(0), nil, nil, nil, nil, data, nil}
	msg, err := utils.NewMsg(ethCall, utils.ToMsgParams(callMsg), hexutil.EncodeBig(big.NewInt(height)))
	if err != nil {
		return nil, err
	}
	//todo: check it's not nil

	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "0x0000")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	// change res method name
	res.Name = method
	return res, nil
}

func RPC_eth_getStorageAt(cfg *utils.Config, method string, height int64, bep20 string, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s: %d, %s", method, height, bep20)
	//
	msg, err := utils.NewMsg(method, bep20, "0x0", hexutil.EncodeBig(big.NewInt(height)))
	if err != nil {
		return nil, err
	}
	//todo: check it's not nil
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "0x0000")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_getBalance(cfg *utils.Config, method string, height int64, extAcc string, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s: %d, %s", method, height, extAcc)
	//
	msg, err := utils.NewMsg(method, extAcc, hexutil.EncodeBig(big.NewInt(height)))
	if err != nil {
		return nil, err
	}
	//todo: check it's not nil
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_getTransactionCount(cfg *utils.Config, method string, height int64, extAcc string, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s: %d, %s", method, height, extAcc)
	//
	msg, err := utils.NewMsg(method, extAcc, hexutil.EncodeBig(big.NewInt(height)))
	if err != nil {
		return nil, err
	}

	//todo: check it's not nil
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_getCode(cfg *utils.Config, method string, height int64, bep20 string, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s: %d, %s", method, height, bep20)
	//
	msg, err := utils.NewMsg(method, bep20, hexutil.EncodeBig(big.NewInt(height)))
	if err != nil {
		return nil, err
	}
	//todo: check it's not nil
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_getProof(cfg *utils.Config, method string, height int64, bep20 string, addrStorageAt []string, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s: %d, %s, %s", method, height, bep20, addrStorageAt[0])
	//
	msg, err := utils.NewMsg(method, bep20, addrStorageAt, hexutil.EncodeBig(big.NewInt(height)))
	if err != nil {
		return nil, err
	}
	//todo: check it's not nil
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_estimateGas(cfg *utils.Config, method string, height int64, addr common.Address, spender common.Address, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s: %d, %s, %s", method, height, addr.Hex(), spender.Hex())
	//
	data, err := cfg.BEP20ABI_.Pack("approve", spender, big.NewInt(1e18))
	if err != nil {
		return nil, err
	}
	callMsg := ethereum.CallMsg{common.HexToAddress(cfg.Sender), &addr, 0, nil, nil, nil, nil, data, nil}

	msg, err := utils.NewMsg(method, utils.ToMsgParams(callMsg), hexutil.EncodeBig(big.NewInt(height)))
	if err != nil {
		return nil, err
	}
	//todo: check it's not nil
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_debug_traceBlockByNumber(cfg *utils.Config, method string, height int64, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s: %d", method, height)
	//
	msg, err := utils.NewMsg(method, hexutil.EncodeBig(big.NewInt(height)), map[string]interface{}{"tracer": "callTracer"})
	if err != nil {
		return nil, err
	}
	//todo: check it's not nil
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_debug_traceTransaction(cfg *utils.Config, method string, txhash common.Hash, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s: %s", method, txhash.Hex())
	//
	msg, err := utils.NewMsg(method, txhash, map[string]interface{}{"tracer": "callTracer"})
	if err != nil {
		return nil, err
	}
	//todo: check it's not nil
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_debug_traceCall(cfg *utils.Config, method string, height int64, bep20Addr common.Address, extAccAddr common.Address, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s: %d, %s", method, height, bep20Addr.Hex(), extAccAddr.Hex())

	data, err := cfg.BEP20ABI_.Pack("balanceOf", extAccAddr)
	if err != nil {
		return nil, err
	}
	callMsg := ethereum.CallMsg{
		From:      extAccAddr,
		To:        &bep20Addr,
		Gas:       uint64(0),
		GasPrice:  nil,
		GasFeeCap: nil,
		Value:     nil,
		Data:      data,
	}
	param := make(map[string]interface{})
	param["tracer"] = "callTracer"
	msg, err := utils.NewMsg(method, utils.ToMsgParams(callMsg), hexutil.EncodeBig(big.NewInt(height)), param)
	if err != nil {
		return nil, err
	}
	//todo: check it's not nil
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_getLogs(cfg *utils.Config, method string, fromBlock, toBlock *big.Int, contractAddress []common.Address, topic0, topic1, topic2, topic3 []common.Hash, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s: fromBlock: %d toBlock: %d , contractAddress: %v, %v,%v,%v, %v", method, fromBlock, toBlock, contractAddress, topic0, topic1, topic2, topic3)
	//
	filterQuery := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: contractAddress,
		Topics:    [][]common.Hash{topic0, topic1, topic2, topic3},
	}
	params, err := utils.ToFilterParams(filterQuery)
	if err != nil {
		return nil, err
	}
	msg, err := utils.NewMsg(method, params)
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func RPC_net_version(cfg *utils.Config, method string, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("method: %s: ", method)
	//
	msg, err := utils.NewMsg(method)
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_blockNumber(cfg *utils.Config, method string, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("method: %s: ", method)
	//
	msg, err := utils.NewMsg(method)
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_newFilter(cfg *utils.Config, method string, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("method: %s: ", method)
	//
	arg := map[string]interface{}{}
	msg, err := utils.NewMsg(method, arg)
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_getFilterLogs(cfg *utils.Config, method string, id string, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("method: %s: ", method)
	//
	msg, err := utils.NewMsg(method, id)
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_newBlockFilter(cfg *utils.Config, method string, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("method: %s: ", method)
	//
	msg, err := utils.NewMsg(method)
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_newPendingTransactionFilter(cfg *utils.Config, method string, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("method: %s: ", method)
	//
	msg, err := utils.NewMsg(method)
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_getTransactionReceipt(cfg *utils.Config, method string, txhash common.Hash, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s: %s", method, txhash.Hex())
	//
	msg, err := utils.NewMsg(method, txhash)
	if err != nil {
		return nil, err
	}
	//todo: check it's not nil
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_getBlockByNumber(cfg *utils.Config, method string, height int64, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s: %d, %s", method, height)
	//
	msg, err := utils.NewMsg("eth_getBlockByNumber", (*hexutil.Big)(big.NewInt(height)), true)
	if err != nil {
		return nil, err
	}
	//todo: check it's not nil
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_chainId(cfg *utils.Config, method string, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("method: %s: ", method)
	//
	msg, err := utils.NewMsg(method)
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_getTransactionByHash(cfg *utils.Config, method string, txhash common.Hash, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("%s: %s", method, txhash.Hex())
	//
	msg, err := utils.NewMsg(method, txhash)
	if err != nil {
		return nil, err
	}
	//todo: check it's not nil
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func checkRPCErr(method string, res *utils.Response) error {
	if strings.Contains(res.Body, "\"result\":") {
		logutil.Debugf("%s: [rt: %d ms], %s",
			method, res.Time, res.Body)
	} else {
		var rpcErr utils.JsonRpcMessage
		json.Unmarshal([]byte(res.Body), &rpcErr)
		msg := fmt.Sprintf("%d: %s",
			rpcErr.Error.Code, rpcErr.Error.Message)
		return errors.New(msg)
	}
	return nil
}

func RPC_bor_getAuthor(cfg *utils.Config, method string, blocknumber int, timeout int64) (*utils.Response, error) {
	s := "0x" + strconv.Itoa(blocknumber)
	msg, err := utils.NewMsg(method, s)
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

//both eth_GetTransactionDataAndReceipt and nr_GetTransactionDataAndReceipt use the method
func RPC_eth_GetTransactionDataAndReceipt(cfg *utils.Config, method string, txhash common.Hash, timeout int64) (*utils.Response, error) {
	//
	logutil.Debugf("method: %s: ", method)
	//
	msg, err := utils.NewMsg(method, txhash.Hex())
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_bor_getSingalByHash(cfg *utils.Config, method string, txhash common.Hash, timeout int64) (*utils.Response, error) {
	//
	msg, err := utils.NewMsg(method, txhash)
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_bor_getCurrent(cfg *utils.Config, method string, timeout int64) (*utils.Response, error) {
	//
	msg, err := utils.NewMsg(method)
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_bor_getRootHash(cfg *utils.Config, method string, txhash int64, timeout int64) (*utils.Response, error) {
	//
	msg, err := utils.NewMsg(method, txhash-3, txhash)
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_eth_getTransactionReceiptsByBlock(cfg *utils.Config, method string, index string, timeout int64) (*utils.Response, error) {
	//
	msg, err := utils.NewMsg(method, "0x"+index)
	if err != nil {
		return nil, err
	}
	res, err := utils.SendMsg(cfg.Host[rand.Intn(len(cfg.Host))], msg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	//
	return res, nil
}

func RPC_btach_top3(cfg *utils.Config, method string, timeout int64) (*utils.Response, error) {
	r := rand.Intn(int(cfg.Range))
	height := cfg.Height - int64(r)
	msg1, err := utils.NewMsg("eth_getBlockByNumber", (*hexutil.Big)(big.NewInt(height)), true)
	if err != nil {
		return nil, err
	}
	r2 := rand.Intn(len(cfg.TXHashList))
	txhash := common.HexToHash(cfg.TXHashList[r2])
	msg2, err := utils.NewMsg("eth_getTransactionReceipt", txhash)
	if err != nil {
		return nil, err
	}

	filterQuery := ethereum.FilterQuery{}
	params, err := utils.ToFilterParams(filterQuery)
	if err != nil {
		return nil, err
	}
	msg3, err := utils.NewMsg(method, params)
	if err != nil {
		return nil, err
	}
	batchMsg := []*utils.JsonRpcMessage{msg1, msg2, msg3}
	res, err := utils.SendBatchMsg(cfg.Host[rand.Intn(len(cfg.Host))], batchMsg, timeout, "")
	if err != nil {
		return nil, err
	}
	err = checkRPCErr(method, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

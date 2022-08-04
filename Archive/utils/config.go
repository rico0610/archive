package utils

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"rpc-load-test/contracts/V2router"
	"rpc-load-test/contracts/bep20"
	"rpc-load-test/contracts/wbnb"
	"rpc-load-test/utils/logutil"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

type Config struct {
	IsStandalone  bool `yaml:"IsStandalone"`
	AWSLambdaNums int  `yaml:"AWSLambdaNums"`

	ScenariosPerSec int `yaml:"ScenariosPerSec"`
	DurationSec     int `yaml:"DurationSec"`
	GoroutineLimit  int `yaml:"GoroutineLimit"`
	LogLevel        int `yaml:"LogLevel"`

	TestNode string `yaml:"TestNode"`
	FullNode string `yaml:"FullNode"`

	ApiKey []string `yaml:"ApiKey"`

	WsUrl string `yaml:"WsUrl"`

	WsHost []string
	Host   []string

	Scenarios map[string]int `yaml:"Scenarios"`

	AccountLoadTill uint64 `yaml:"AccountLoadTill"`
	Height          int64  `yaml:"Height"`
	Range           int64  `yaml:"Range"`

	Router string   `yaml:"Router"`
	WBNB   string   `yaml:"WBNB"`
	BEP20  []string `yaml:"BEP20"`

	NumberOfTXHash  int    `yaml:"NumberOfTXHash"`
	NumberOfAccount int    `yaml:"NumberOfAccount"`
	AccountCSV      string `yaml:"AccountCSV"`
	TXHashCSV       string `yaml:"TXHashCSV"`

	Sender         string `yaml:"Sender"`
	BEP20StorageAt string `yaml:"BEP20StorageAt"`

	AggrInterval int         `yaml:"AggrInterval"`
	Headers      http.Header `yaml:"Headers"`

	// key is bep20 address , value is account address
	AccountList   map[string][]string
	TXHashList    []string
	StartTime     time.Time
	FullClient    *ethclient.Client
	BlockHashList []string

	RouterAddr common.Address
	WBNBAddr   common.Address
	Bep20Addrs []common.Address

	V2routerABI_ abi.ABI
	WBNBABI_     abi.ABI
	BEP20ABI_    abi.ABI

	FilterList          []string
	FilterLogList       []string
	FilterUninstallList []string
}

func init() {
	// these are http specific config for high load efficiency
	http.DefaultTransport.(*http.Transport).MaxIdleConns = 1000
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 1000
}

var T_cfg = &Config{}

func (cfg *Config) LoadYML(testnode, fullnode, wsUrl, apiKey string, qps, dur, lLevel, goroutineLimit int) error {
	//
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := filepath.Join(wd, "config-bsc-mainnet.yml")
	//
	configYML, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(configYML, cfg)
	if err != nil {
		return err
	}
	//
	if testnode != "" {
		cfg.TestNode = testnode
	}
	if fullnode != "" {
		cfg.FullNode = fullnode
	}
	if wsUrl != "" {
		cfg.WsUrl = wsUrl
	}
	if qps != 0 {
		cfg.ScenariosPerSec = qps
	}
	if dur != 0 {
		cfg.DurationSec = dur
	}
	if lLevel != 0 {
		cfg.LogLevel = lLevel
	}
	if goroutineLimit != 0 {
		cfg.GoroutineLimit = goroutineLimit
	}
	if apiKey != "" {
		keys := strings.Split(apiKey, ",")
		cfg.ApiKey = keys
	}
	for _, k := range cfg.ApiKey {
		cfg.Host = append(cfg.Host, cfg.TestNode+k)
		cfg.WsHost = append(cfg.WsHost, cfg.WsUrl+k)
		logutil.Infof("wsHost: %s", cfg.WsUrl+k)
		logutil.Infof("TestHost: %s", cfg.TestNode+k)
	}
	// support not have apiKey
	if len(cfg.ApiKey) == 0 {
		logutil.Info("--- no api key load test--")
		cfg.Host = append(cfg.Host, cfg.TestNode)
		cfg.WsHost = append(cfg.WsHost, cfg.WsUrl)
	}
	//
	logutil.InitLog(cfg.LogLevel, logutil.Stdout)
	//
	weighted := make(map[string]int)
	for k, v := range cfg.Scenarios {
		if v > 0 {
			weighted[k] = v
		}
	}
	cfg.Scenarios = weighted
	//
	if cfg.Height == 0 {
		cfg.FullClient, err = GetClient(cfg.FullNode)
		if err != nil {
			panic(err)
		}
		cfg.Height, err = GetHeight(cfg.FullClient)
		if err != nil {
			panic(err)
		}
	}
	//
	if cfg.Range == 0 {
		cfg.Range = 1
	}
	//
	cfg.RouterAddr = common.HexToAddress(cfg.Router)
	cfg.WBNBAddr = common.HexToAddress(cfg.WBNB)
	for _, v := range cfg.BEP20 {
		addr := common.HexToAddress(v)
		cfg.Bep20Addrs = append(cfg.Bep20Addrs, addr)
	}
	//
	cfg.V2routerABI_, err = abi.JSON(strings.NewReader(V2router.V2routerABI))
	if err != nil {
		return err
	}
	cfg.WBNBABI_, err = abi.JSON(strings.NewReader(wbnb.WbnbABI))
	if err != nil {
		return err
	}
	cfg.BEP20ABI_, err = abi.JSON(strings.NewReader(bep20.Bep20ABI))
	if err != nil {
		return err
	}
	//
	if cfg.AccountCSV == "" {
		cfg.AccountList = GetAccountsFromLogs(cfg.FullNode, nil, cfg.NumberOfAccount)
	} else {
		cfg.AccountList, err = LoadAccountCSV(cfg.AccountCSV)
	}
	if err != nil {
		panic(err)
	}
	//
	if cfg.TXHashCSV == "" {
		cfg.TXHashList, cfg.BlockHashList = getHashFromBlock(cfg.FullNode, cfg.NumberOfTXHash)
	} else {
		cfg.TXHashList, err = LoadCSV(cfg.TXHashCSV)
		if err != nil {
			panic(err)
		}
	}
	for s, _ := range cfg.Scenarios {
		if strings.Contains(s, "Filter") {
			cfg.FilterList = getFilterList(cfg.FullNode, cfg.NumberOfTXHash, 1)
			cfg.FilterLogList = getFilterList(cfg.FullNode, cfg.NumberOfTXHash, 0)
			//todo add  update
			cfg.FilterUninstallList = getFilterList(cfg.FullNode, cfg.NumberOfTXHash, 2)
		}
	}

	cfg.StartTime = time.Now()
	logutil.Infof("qps: %d , testNode: %s , dur: %d , goroutineLimit: %d , wsUrl: %s ", cfg.ScenariosPerSec, cfg.TestNode, cfg.DurationSec, cfg.GoroutineLimit, cfg.WsUrl)
	logutil.Infof("===init config success===")
	return nil
}

func LoadCSV(fpath string) ([]string, error) {
	//
	list := make([]string, 0, 1000)
	//
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	//
	fpath = filepath.Join(wd, fpath)
	dataCSV, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer dataCSV.Close()
	//
	scan := bufio.NewScanner(dataCSV)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		list = append(list, line)
	}
	return list, nil
}

func LoadAccountCSV(fpath string) (map[string][]string, error) {

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	//
	fpath = filepath.Join(wd, fpath)
	dataCSV, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer dataCSV.Close()
	//
	ac, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	contactAddress := ""
	accounts := map[string][]string{}
	for _, line := range strings.Split(string(ac), "\n") {
		data := strings.Split(line, ",")
		if len(data) == 1 {
			contactAddress = data[0]
		} else {
			accounts[contactAddress] = data
		}
	}
	return accounts, nil
}

func getHashFromBlock(rpc string, target int) ([]string, []string) {
	txHashs := make([]string, target)
	client, err := ethclient.Dial(rpc)
	hashlist := make([]string, target)
	if err != nil {
		panic(err)
	}
	currBlockNumber, err := client.BlockNumber(context.Background())

	if err != nil {
		logutil.Errorf("get currBlockNumber err: %v ", err)
	}
	path := fmt.Sprintf("txHash-%d-%d.csv", currBlockNumber, target)
	File, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logutil.Errorf("OpenFile err: %v ", err)
	}
	defer File.Close()
	logutil.Infof("====start load txHash, target: %d, store to file: %s=====", target, path)
	WriterCsv := csv.NewWriter(File)
	count := 0
	for {
		bls, err := client.BlockByNumber(context.Background(), big.NewInt(int64(currBlockNumber)))
		if err != nil {
			logutil.Errorf("sync block err: %s", err)
		}
		for _, trx := range bls.Transactions() {
			h := trx.Hash().String()
			logutil.Debug("get txHash: ", h)
			err1 := WriterCsv.Write([]string{h})
			if err1 != nil {
				logutil.Errorf("Write err: %v ", err)
			}
			txHashs[count] = h
			hashlist[count] = bls.Hash().String()
			count++
			if count >= target {
				WriterCsv.Flush()
				logutil.Info("=== load txHash success===")
				return txHashs, hashlist
			}
		}
		logutil.Debugf("sync block: %d , hash count: %d ", currBlockNumber, count)
		WriterCsv.Flush()
		currBlockNumber--
	}
}

func getFilterList(rpc string, target int, method int) []string {
	Filters := make([]string, target)
	var m string
	switch method {
	case 0:
		m = "eth_newFilter"
	case 1:
		m = "eth_newBlockFilter"
	case 2:
		m = "eth_newPendingTransactionFilter"
	default:
		m = "eth_newBlockFilter"
	}
	for i := 0; i < target; i++ {
		var msg *JsonRpcMessage
		var err error
		if m == "eth_newFilter" {
			arg := map[string]interface{}{}
			msg, err = NewMsg(m, arg)
		} else {
			msg, err = NewMsg(m)
		}
		if err != nil {
			return nil
		}
		res, err := SendMsg(rpc, msg, 300, "")
		if err != nil {
			return nil
		}
		Filters[i] = gjson.Get(res.Body, "result").String()

	}

	return Filters

}

func GetAccountsFromLogs(rpc string, addresses []common.Address, target int) map[string][]string {
	testAccounts := map[string][]string{}
	client, err := ethclient.Dial(rpc)
	if err != nil {
		panic(err)
	}

	currBlockNumber := uint64(0)
	if T_cfg.AccountLoadTill != 0 {
		currBlockNumber = T_cfg.AccountLoadTill
	} else {
		currBlockNumber, err = client.BlockNumber(context.Background())
		if err != nil {
			logutil.Errorf("get currBlockNumber err: %v ", err)
		}
	}
	path := fmt.Sprintf("account-%d-%d.csv", currBlockNumber, target)
	File, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		logutil.Errorf("OpenFile err: %v ", err)
	}
	defer File.Close()
	logutil.Infof("====start load testAccounts, count: %d, store to file: %s=====", target, path)
	WriterCsv := csv.NewWriter(File)
	index := int64(currBlockNumber)
	var bep20Address []common.Address
	if addresses != nil {
		bep20Address = addresses
	} else {
		bep20Address = T_cfg.Bep20Addrs
	}
	for {
		// filter bep20 topic ,range 10 block/mainnet  range 100(use Topics: transfer/approve )
		transferTopic := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
		approveTopic := common.HexToHash("0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c")
		filterQuery := ethereum.FilterQuery{
			FromBlock: big.NewInt(index - 500),
			ToBlock:   big.NewInt(index),
			Topics:    [][]common.Hash{{transferTopic, approveTopic}},
			Addresses: bep20Address,
		}
		index = index - 500
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		topics, err := client.FilterLogs(ctx, filterQuery)
		cancel()
		if err != nil {
			logutil.Errorf("fail to filter logs: %v", err)
			continue
		}
		logutil.Infof("index: %d,topics: %d", index, len(topics))
		for _, t := range topics {
			contractAddress := t.Address.Hex()
			// check topic length prevent panic
			if len(t.Topics) < 3 {
				continue
			}
			a1 := common.HexToAddress(t.Topics[1].Hex()).Hex()
			a2 := common.HexToAddress(t.Topics[2].Hex()).Hex()
			if !IsContain(testAccounts[contractAddress], a1) {
				testAccounts[contractAddress] = append(testAccounts[contractAddress], a1)
			}
			if !IsContain(testAccounts[contractAddress], a2) {
				testAccounts[contractAddress] = append(testAccounts[contractAddress], a2)
			}
		}

		// remove already bep20 address
		for c, a := range testAccounts {
			if len(a) >= target {
				for i, address := range bep20Address {
					if address.Hex() == c {
						bep20Address = append(bep20Address[:i], bep20Address[i+1:]...)
						WriterCsv.Write([]string{c})
						WriterCsv.Write(testAccounts[c])
						WriterCsv.Flush()

						logutil.Infof("%s load success: %d", c, len(a))
					}
				}
			}
		}

		if len(bep20Address) == 0 {
			logutil.Infof("===load accounts from logs success===")
			return testAccounts
		}
	}
}

func IsContain(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

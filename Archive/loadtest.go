package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"

	"rpc-load-test/test"
	"rpc-load-test/utils"
	"rpc-load-test/utils/logutil"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/panjf2000/ants/v2"

	"go.uber.org/ratelimit"
)

var (
	host    	   *string
	rpc			   *string
	blocks		   *int
	token		   *string
	users   	   *string
	qps            *int
	dur            *int
	lLevel         *int
	count          = 0
)

func init() {
	host = flag.String("host", "", "host")
	rpc = flag.String("rpc", "", "rpc")
	blocks = flag.Int("blocks", 10000, "blocks")
	token = flag.String("token", "", "erc20/bep20 hex address")
	users = flag.String("users", "", "csv of external accounts")
	qps = flag.Int("qps", 1, "qps")
	dur = flag.Int("dur", 10, "dur")
	lLevel = flag.Int("l", 1, "log level: 1=debug; 2=Info; 3=Error")
	flag.Parse()
	utils.T_cfg.IsStandalone = true
	utils.T_cfg.AggrInterval = 10
	utils.T_cfg.Host = []string{*host}
	utils.T_cfg.Scenarios = map[string]int{*rpc:100}
	utils.T_cfg.Range = int64(*blocks)
	utils.T_cfg.BEP20 = []string{*token}
	var err error
	if *users != "" {
		utils.T_cfg.AccountList, err = utils.LoadAccountCSV(*users)
		if err != nil {
			panic(err)
		}
	} else {
		addresses := []common.Address{common.HexToAddress(*token)}
		utils.T_cfg.AccountList = utils.GetAccountsFromLogs(*host, addresses, 10000)
	}
	logutil.Infof("Numbers of addresses: %d", len(utils.T_cfg.AccountList[*token]))
	utils.T_cfg.ScenariosPerSec = *qps
	utils.T_cfg.DurationSec = *dur
	utils.T_cfg.LogLevel = *lLevel
	logutil.InitLog(*lLevel, logutil.Stdout)
	utils.T_cfg.GoroutineLimit = *qps
	client, err := ethclient.Dial(*host)
	if err != nil {
		panic(err)
	}
	var height uint64
	height, err = client.BlockNumber(context.Background())
	if err != nil {
		panic(err)
	}
	utils.T_cfg.Height = int64(height)
}

func main() {
	//
	if utils.T_cfg.IsStandalone {
		HandleRequest()
		return
	}
	//
	lambda.Start(HandleRequest)
}

func HandleRequest() (string, error) {
	//
	rand.Seed(time.Now().UnixNano())
	//
	exec(utils.T_cfg)
	//
	return "Completed", nil
}

func setupTimer(dur time.Duration) *bool {
	t := time.NewTimer(dur)
	expired := false
	go func() {
		<-t.C
		expired = true
	}()
	return &expired
}

func exec(cfg *utils.Config) {
	//
	rlimiter := ratelimit.New(cfg.ScenariosPerSec)
	//
	resCurr := make(map[string][]int64)
	fCurr := make(map[string]int64)
	resAggr := make(map[string][]Result)
	//
	dSec := time.Duration(cfg.DurationSec) * time.Second
	expired := setupTimer(dSec)
	time.Sleep(1 * time.Second)
	//
	var wg sync.WaitGroup
	var m sync.Mutex
	// query fullnode height
	if cfg.FullClient != nil {
		go func() {
			for {
				time.Sleep(3 * time.Second)
				//
				if *expired {
					break
				}
				//
				height, err := utils.GetHeight(cfg.FullClient)
				if err != nil {
					logutil.Error("GetHeight: " + err.Error())
					continue
				}
				//
				m.Lock()
				cfg.Height = height
				m.Unlock()
				//
				logutil.Infof("height: %d", cfg.Height)
			}
		}()
	}
	// collect batch results
	var writer *bufio.Writer
	if cfg.IsStandalone {
		//
		timestamp := time.Now().UnixNano()
		fpath := fmt.Sprintf("./%d_%d.csv", cfg.ScenariosPerSec, timestamp)
		resfile, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		defer resfile.Close()
		writer = bufio.NewWriter(resfile)
		//
		go func() {
			for {
				time.Sleep(time.Duration(cfg.AggrInterval) * time.Second)
				//
				if *expired {
					break
				}
				//
				m.Lock()
				successes := resCurr
				resCurr = map[string][]int64{}
				failures := fCurr
				fCurr = map[string]int64{}
				m.Unlock()
				//
				collectResults(successes, failures, resAggr, writer)
			}
		}()
	}

	if cfg.GoroutineLimit == 0 {
		cfg.GoroutineLimit = math.MaxInt
	}
	pool, _ := ants.NewPool(cfg.GoroutineLimit)
	defer pool.Release()
	for {
		if *expired {
			break
		}

		rlimiter.Take()
		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()
			//
			res := test.Run(cfg)
			//
			m.Lock()
			for _, v := range res {
				if v.Code != -1 {
					resCurr[v.Name] = append(resCurr[v.Name], v.Time)
				} else {
					fCurr[v.Name]++
				}
			}
			m.Unlock()
		})
	}
	wg.Wait()
	//
	m.Lock()
	successes := resCurr
	resCurr = map[string][]int64{}
	failures := fCurr
	fCurr = map[string]int64{}
	m.Unlock()
	//
	collectResults(successes, failures, resAggr, writer)
	//
	if !cfg.IsStandalone {
		logutil.Info("Final:", resAggr)
		return
	}
	//
	for method, res := range resAggr {
		points := utils.NewPoints()
		var avg, l50, l90, l99, err, success, l int64
		for _, v := range res {
			if v.RtAVG > 0 {
				l++
			}
			points.AVG = append(points.AVG, float64(v.RtAVG))
			points.P50 = append(points.P50, float64(v.RtP50))
			points.P90 = append(points.P90, float64(v.RtP90))
			points.P99 = append(points.P99, float64(v.RtP99))
			avg = avg + v.RtAVG
			l50 = l50 + v.RtP50
			l90 = l90 + v.RtP90
			l99 = l99 + v.RtP99
			err = err + v.Failure
			success = success + v.Success
		}

		fmt.Println("RPC,QPS,AVG,L50,L90,L99,Error(%),Success,Error")
		//l := int64(len(points.AVG))
		errorRate := float32(err*10000/(err+success)) / float32(100)
		fmt.Printf("%s,%d,%d,%d,%d,%d,%f,%d,%d\n", method, utils.T_cfg.ScenariosPerSec, avg/l, l50/l, l90/l, l99/l, errorRate, success, err)

		points.RenderChart("./results", fmt.Sprintf("%d-%s.png", utils.T_cfg.ScenariosPerSec, method))
	}
}

type Result struct {
	Success int64
	Failure int64
	RtAVG   int64
	RtP50   int64
	RtP90   int64
	RtP99   int64
}

func collectResults(successes map[string][]int64, failures map[string]int64, resAggr map[string][]Result, writer *bufio.Writer) {
	//
	if len(successes) == 0 {
		for method, _ := range failures {
			result := Result{0, failures[method], 0, 0, 0, 0}
			resAggr[method] = append(resAggr[method], result)
			logutil.Infof("all is failed: %d", failures[method])
		}
	}
	for method, numbers := range successes {
		sort.Slice(numbers, func(i int, j int) bool {
			return numbers[i]-numbers[j] < 0
		})
		//
		var rtSum int64 = 0
		success := int64(len(numbers))
		//
		for _, v := range numbers {
			rtSum += v
		}
		var avg int64 = 0
		if success != 0 {
			avg = rtSum / success
		}
		l50 := numbers[int64(float64(success)*0.5)]
		l90 := numbers[int64(float64(success)*0.9)]
		l99 := numbers[int64(float64(success)*0.99)]
		//

		if writer != nil {
			line := fmt.Sprintf("%s,%d,%d,%d,%d,%d,%d\n",
				method, success, failures[method], avg, l50, l90, l99)
			writer.WriteString(line)
			writer.Flush()
		}
		logutil.Infof("run: %d, method: %s ; success: %d ; error: %d, (ms) avg: %d ; 50%%: %d ; 90%%: %d ; 99%%: %d",
			count, method, success, failures[method], avg, l50, l90, l99)

		result := Result{success, failures[method], avg, l50, l90, l99}
		resAggr[method] = append(resAggr[method], result)
	}
	logutil.Info("----print result----")
	count++
}

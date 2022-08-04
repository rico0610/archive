package utils

import (
	"context"
	"errors"
	"math/rand"
	"regexp"
	"time"

	"rpc-load-test/utils/logutil"

	"github.com/ethereum/go-ethereum/ethclient"
)

const base = "0123456789abcdefghijklmnopqrstuvwxyz"

func RandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func RandomString(bytes []byte) *string {
	for i, b := range bytes {
		bytes[i] = base[b%36]
	}
	str := string(bytes)
	return &str
}

func Search(text string, pattern string) ([]string, error) {
	regexp_, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	match := regexp_.FindStringSubmatch(text)
	if match == nil {
		return nil, errors.New("not found")
	}
	return match, nil
}

func GetClient(url string) (*ethclient.Client, error) {
	ctx, _ := context.WithTimeout(
		context.Background(), 10*time.Second)
	cli, err := ethclient.DialContext(ctx, url)
	logutil.Debugf("GetClient: %d, err: %v", cli, err)

	return cli, err
}

func GetHeight(cli *ethclient.Client) (int64, error) {
	ctx, _ := context.WithTimeout(
		context.Background(), 10*time.Second)
	height, err := cli.BlockNumber(ctx)
	logutil.Debugf("GetHeight: %d, err: %v", height, err)
	return int64(height), err
}

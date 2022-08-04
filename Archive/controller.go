package main

import (
	"fmt"
	"time"

	"rpc-load-test/utils"
	"rpc-load-test/utils/logutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func main() {
	//
	err := utils.T_cfg.LoadYML("", "", "", "", 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	//
	for i := 0; i < utils.T_cfg.AWSLambdaNums; i++ {
		callLambda(i)
		if i >= 500 {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func callLambda(index int) {
	fmt.Printf("#%d lambda function invoked at %s\n", index, time.Now().String())
	svc := lambda.New(session.New(), &aws.Config{Region: aws.String("us-east-1")})
	input := &lambda.InvokeInput{
		ClientContext: aws.String("loadtest"),
		FunctionName:  aws.String("loadtest-dev-loadtest"),
		// Event is for async mode, success status: 202
		// RequestResponse is for sync mode, success status: 200
		InvocationType: aws.String("Event"),
		LogType:        aws.String("Tail"),
		Payload:        nil,
	}
	res, err := svc.Invoke(input)
	if err != nil {
		logutil.Error(err)
	}
	logutil.Info(res)
}

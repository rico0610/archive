package main

import "rpc-load-test/utils"

func main() {
	points := utils.NewPoints()
	points.LoadDataFromCVS("./results", "5_1642081977980617000.csv", "eth_call")
	points.RenderChart("./results", "eth_call.png")
	//
	points = utils.NewPoints()
	points.LoadDataFromCVS("./results", "5_1642081977980617000.csv", "eth_getStorageAt")
	points.RenderChart("./results", "eth_getStorageAt.png")
	//
	points = utils.NewPoints()
	points.LoadDataFromCVS("./results", "5_1642081977980617000.csv", "eth_getBalance")
	points.RenderChart("./results", "eth_getBalance.png")
	//
	points = utils.NewPoints()
	points.LoadDataFromCVS("./results", "5_1642081977980617000.csv", "eth_getTransactionCount")
	points.RenderChart("./results", "eth_getTransactionCount.png")
	//
	points = utils.NewPoints()
	points.LoadDataFromCVS("./results", "5_1642081977980617000.csv", "eth_getCode")
	points.RenderChart("./results", "eth_getCode.png")
	//
	points = utils.NewPoints()
	points.LoadDataFromCVS("./results", "5_1642081977980617000.csv", "eth_getProof")
	points.RenderChart("./results", "eth_getProof.png")
	//
	points = utils.NewPoints()
	points.LoadDataFromCVS("./results", "5_1642081977980617000.csv", "eth_estimateGas")
	points.RenderChart("./results", "eth_estimateGas.png")
}

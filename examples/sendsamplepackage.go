package main

import (
	"fmt"
	"questocat/volc-sdk-x-golang/riskconsole"
)

const (
	Ak = "AK"
	Sk = "SK"
)

func main() {
	var instance = riskconsole.NewRiskConsoleInstance()
	instance.Client.SetAccessKey(Ak)
	instance.Client.SetSecretKey(Sk)
	request := riskconsole.PushTrafficRiskDataRequest{
		AppId:        579824,
		FilePath:     "/tmp/权重测试222.csv",
		BusinessType: "A1",
		Scene:        "1",
		DataType:     "1",
	}
	response, err := instance.PushTrafficRiskData(&request)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("result: %+v", response)
}

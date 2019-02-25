// example
package main

import (
	"github.com/tggo/antigate"
	"fmt"
	)

func main() {

	//baseURL , _ := url.Parse(client2.DefaultBaseURL)
	//client := client2.Client{
	//	Key: "ccbdc36ec3274cf8fe7b49fc2d8733e4",
	//	BaseURL: baseURL,
	//}
	//
	//client := client.NewClient(nil)
	//client.GetBalance(context.Background())
	// Configuration
	config := &antigate.Config{
		ClientKey:   "ccbdc36ec3274cf8fe7b49fc2d8733e4",
	}

	// Create a new client with the above configuration.
	client := antigate.NewClient(config)

	//balance, _ := client.Balance.Get()
	//fmt.Printf("%v", balance)


	task := antigate.Task{
		TaskType: "NoCaptchaTaskProxyless",
		WebsiteURL: "https://cabinet.sfs.gov.ua/cashregs/check",
		WebsiteKey: "6Lc3DE8UAAAAAIc2N3jarTo9_R_DuooXTFxJYPqa",

	}
	response, _ := client.Task.GetKeyForGoogle(task)
	//responseTaskId, err := client.Task.PutToWork( task )
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Printf("taskID: %d\n",responseTaskId)
	//
	//
	//
	//fmt.Println("wait 70 sec")
	//time.Sleep(70 * time.Second)
	//
	////responseTaskId := int64(504938809)
	//
	//fmt.Println("start")
	//response, err := client.Task.GetWork( responseTaskId)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//if response.Status == "processing" {
	//	fmt.Println("processing")
	//}

	fmt.Printf("response: %s\n", response)

}
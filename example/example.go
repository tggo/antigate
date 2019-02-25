// example
package main

import (
	"fmt"
	"github.com/tggo/antigate"
)

func main() {

	config := &antigate.Config{
		ClientKey: "ccbdc36ec3274cf8fe7b49fc2d8733e0",
	}

	client := antigate.NewClient(config)

	balance, _ := client.Balance.Get()
	fmt.Printf("%v", balance)

	task := antigate.Task{
		TaskType:   "NoCaptchaTaskProxyless",
		WebsiteURL: "https://google.com",
		WebsiteKey: "6Lc3DE8UAAAAAIc2N3jarTo9_R_DuooXFxJYPqa",
	}
	response, _ := client.Task.GetKeyForGoogle(task)

	fmt.Printf("response: %s\n", response)

}

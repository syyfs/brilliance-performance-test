package main

import (
	"fmt"
	"brilliance-performance-test/chaincodefunction"
	"github.com/spf13/viper"
	"time"
	"sync"
)

var (
	wg *sync.WaitGroup
)
var configPath = "./config/config.yaml"
var httpscertPath = "./httpscert/unionbank.crt"

func init() {
	fmt.Printf("InitYamlConfig configpath:%s \n", configPath)
	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("===== ReadInConfig faild , err is [%s]======\n", err)
	}
}


func main() {
	wg = new(sync.WaitGroup)
	// 定义一个chan,从chan中获取tocken
	token , err  := chaincodefunction.Login(httpscertPath)
	if err != nil {
		fmt.Errorf(" Invoke Login faild ! err is %s\n", err)
	}

	time.Sleep(time.Second * 5)

	fmt.Println("In main()")
	wg.Add(1)
	go invoke(token)

	wg.Add(1)
	go query(token)

	wg.Add(1)
	go querycc(token)

	wg.Wait()
	fmt.Println("At the end of main()")

}

func invoke (token string){
	defer wg.Done()
	var i int = 0
	invokeChan := make(chan int , 10)
	defer close(invokeChan)
	for {
		i++
		invokeChan <- i // 如果chan被写满了，那么chan就会阻塞当前协程,直到当前协程空间被释放出来
		go chaincodefunction.ChaincodeInvoke(i , token, httpscertPath , invokeChan, wg)

	}

}

func query(token string)  {
	defer wg.Done()
	var j int = 0
	queryChan := make(chan int , 10)
	defer close(queryChan)
	for {
		j++
		queryChan <- j
		go chaincodefunction.ChaincodeQuery(j,token, httpscertPath,queryChan,wg)
	}

}

func querycc(token string)  {
	defer wg.Done()
	var z int =0
	chanquery := make(chan int , 10)
	defer close(chanquery)

	for{
		z++
		chanquery <- z
		go chaincodefunction.QueryChaincode(z,token, httpscertPath,chanquery,wg)
	}
}






package main

import (
	"fmt"
	"brilliance-performance-test/chaincodefunction"
	"github.com/spf13/viper"
	"sync"
	"flag"
)

var (
	wg *sync.WaitGroup
)
var configPath = "./config/config.yaml"
var HttpscertPath = "./httpscert/unionbank.crt"

var channelName = flag.String("channel_id","mychannel1","channel name!")
var ccName = flag.String("cc_name","bclc1","channel name!")
var ccFcn = flag.String("cc_fcn","UnionChainSaveData","channel name!")

//var ccArgs = flag.String("cc_args","mychannel1","channel name!")

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
	flag.Parse()
	data := fmt.Sprintf("{\"Data\":{\"TxData\":{\"curId\":\"1\",\"curStep\":0,\"cryptoData\":\"12345\",\"rootId\":\"1\",\"recieveMember\":{\"member\":\"citic1\",\"key\":\"123\"},\"notifyMember\":[{\"member\":\"citic2\",\"key\":\"123\"}],\"LC_TX_CODE\":\"BCL0101\"}}}")
	fmt.Printf(" %c[%d;%d;%dm channelName:[%s];ccName:[%s];ccFcn:[%s];data:[%s]; %c[0m \n", 0x1B, 1, 40, 31, *channelName,*ccName,*ccFcn,data, 0x1B) // 红色

	// 定义一个chan,从chan中获取tocken
	//token , err  := chaincodefunction.Login(HttpscertPath)
	//if err != nil {
	//	fmt.Errorf(" Invoke Login faild ! err is %s\n", err)
	//}
	invokeChaincodeCfg := &chaincodefunction.InvokeChaincodeCfg{
		ChannelId: *channelName,
		CcName: *ccName,
		CcFcn:*ccFcn,
		CcArgs: []string{data},
	}
	fmt.Println("In main()")
	/** 同一个 用户不同客户端 同时一次最大支持100个invoke并发交易 **/
	wg.Add(1)
	go invoke(1, invokeChaincodeCfg)

	wg.Add(1)
	go invoke(2,invokeChaincodeCfg)

	/** 同一个 用户不同客户端 同时一次最大支持100个query并发交易 **/
	//wg.Add(1)
	//go query()
	//
	//wg.Add(1)
	//go query()
	//
	///** 同一个 用户不同客户端 同时一次最大支持100个querycc 并发交易【不进行共识写块交易】 **/
	//wg.Add(1)
	//go querycc()
	//
	//wg.Add(1)
	//go querycc()

	wg.Wait()
	fmt.Println("At the end of main()")

}

func invoke (croutine int ,cci *chaincodefunction.InvokeChaincodeCfg){
	defer wg.Done()
	token , err  := chaincodefunction.Login(HttpscertPath)
	if err != nil {
		fmt.Errorf(" Invoke Login faild ! err is %s\n", err)
	}
	var i int = 0
	invokeChan := make(chan int , 50)
	defer close(invokeChan)
	for {
		i++
		invokeChan <- i // 如果chan被写满了，那么chan就会阻塞当前协程,直到当前协程空间被释放出来
		go chaincodefunction.ChaincodeInvoke(cci ,croutine ,i , token, HttpscertPath , invokeChan)

	}
}

func query()  {
	defer wg.Done()

	token , err  := chaincodefunction.Login(HttpscertPath)
	if err != nil {
		fmt.Errorf(" Invoke Login faild ! err is %s\n", err)
	}

	var j int = 0
	queryChan := make(chan int , 100)
	defer close(queryChan)
	for {
		j++
		queryChan <- j
		go chaincodefunction.ChaincodeQuery(1, j,token, HttpscertPath,queryChan)
	}

}

func querycc()  {
	defer wg.Done()

	token , err  := chaincodefunction.Login(HttpscertPath)
	if err != nil {
		fmt.Errorf(" Invoke Login faild ! err is %s\n", err)
	}

	var z int =0
	chanquery := make(chan int , 10)
	defer close(chanquery)

	for{
		z++
		chanquery <- z
		go chaincodefunction.QueryChaincode(z,token, HttpscertPath,chanquery)
	}
}






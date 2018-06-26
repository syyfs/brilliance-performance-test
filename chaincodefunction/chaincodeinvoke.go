package chaincodefunction

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"bytes"
	"encoding/json"
	"brilliance-performance-test/model"
	"time"
	"brilliance-performance-test/utils"
	"crypto/x509"
	"crypto/tls"
	"errors"
	"brilliance/fast-deploy/common/log"
)

type InvokeChaincodeCfg struct {
	ChannelId string   `json:"channel_id"`
	CcName    string   `json:"cc_name"`
	CcFcn     string   `json:"cc_fcn"`
	CcArg     string   `json:"cc_arg"`
	CcArgs    []string `json:"cc_args"`
}

func Login(certPath string)  (token string , err error) {

	var resp *http.Response
	loginreq := &model.LoginReq{
		Username: "admin" ,
		Password: "brilliance" ,
	}

	loginreqbyte , err  := json.Marshal(loginreq)
	if err != nil {
		return "" , err
	}
	commonName := utils.GetCertCommonName()
	port := utils.GetClientPort()
	log.Warningf("commonName:%s, port:%s, cert:[%s]\n",commonName,port, certPath )
	if utils.GetHttpsEnable() {
		fmt.Printf("***** https *******\n")
		client , err := HttpsReqClient(certPath)
		if err != nil {
			return "" , err
		}

		url := fmt.Sprintf("https://%s%s/admin/login",commonName , port)
		fmt.Printf("**** url:[%s]****\n",url)
		resp , err = client.Post(url , "application/json" , bytes.NewBuffer(loginreqbyte))
	} else {

		fmt.Printf("***** http *******")
		client := &http.Client{}
		url := fmt.Sprintf("http://%s%s/admin/login",commonName , port)
		resp , err = client.Post(url, "application/json" , bytes.NewBuffer(loginreqbyte))
		//resp , err = http.Post("http://fast-deploy:8888/admin/login" , "application/json" , bytes.NewBuffer(loginreqbyte))


	}
	if err != nil {
		return "" , err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "" , err
	}

	loginresp := &model.LoginResponse{}
	err = json.Unmarshal(body , loginresp)
	if err != nil {
		return "" , err
	}
	return loginresp.Data.Token , nil
}

func ChaincodeInvoke(cci *InvokeChaincodeCfg,croutine int , i int , token string , cert string , invokeChan chan int){
	currentTime := time.Now().Unix()
	data := fmt.Sprintf("{\"Data\":{\"TxData\":{\"curId\":\"1\",\"curStep\":%d,\"cryptoData\":\"12345\",\"rootId\":\"1\",\"recieveMember\":{\"member\":\"citic1\",\"key\":\"123\"},\"notifyMember\":[{\"member\":\"citic2\",\"key\":\"123\"}],\"LC_TX_CODE\":\"BCL0101\"}}}",i)
	cci.CcArgs = []string{data}
	cfgBytes , err := json.Marshal(cci)
	if err != nil {
		fmt.Errorf(" ==== chaincodeinvoke marshal faild ! err is %s \n ", err)
	}

	err = execute(croutine ,i, cfgBytes , token, cert)
	if err != nil {
		panic(fmt.Errorf(" ==== execute clientDo  faild ! err is %s \n", err))
	}
	tm := time.Now().Unix() - currentTime
	fmt.Printf(" %c[%d;%d;%dm goroutine:[%d],invoke:第 [%v] 次 ; 所用时间 [%v]%c[0m \n", 0x1B, 1, 40, 32, croutine , i , tm, 0x1B) // 绿色
	//fmt.Printf("=======Invoke 第 [%v] 次 , 所用时间 [%v] =======\n", i , tm)
	<- invokeChan
}

func execute(croutine int ,i int , cfgBytes []byte,token string , cert string) error {
	if utils.GetHttpsEnable(){
		client , err := HttpsReqClient(cert)
		if err != nil {
			return fmt.Errorf(" ==== execute HttpsReqClient faild ! err is %s \n ", err)
		}
		url := fmt.Sprintf("https://%s%s/chaincode/execute", utils.GetCertCommonName() , utils.GetClientPort())
		err = clientDo(croutine ,i ,token , client , "POST" , url, cfgBytes  )
		if err != nil {
			panic(fmt.Errorf(" ==== execute clientDo  faild ! err is %s \n ", err))
		}
	} else {
		client := &http.Client{}
		//req, err := http.NewRequest("POST", "http://fast-deploy:8888/chaincode/execute", body)
		url := fmt.Sprintf("http://%s%s/chaincode/execute", utils.GetCertCommonName() , utils.GetClientPort())
		err := clientDo(croutine ,i ,token , client , "POST" , url, cfgBytes  )
		if err != nil {
			return fmt.Errorf(" ==== execute clientDo  faild ! err is %s \n ", err)
		}
	}
	return nil
}

func clientDo(croutine int ,i int,token string,client *http.Client ,method string, url string, cfgBytes []byte ) error {

	cibody := bytes.NewBuffer(cfgBytes)
	req, err := http.NewRequest(method, url , cibody)
	if err != nil {
		return	fmt.Errorf(" ==== [clientDo] NewRequest faild ! err is %s \n ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-Auth-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf(" ==== [clientDo] client.Do(req) faild ! err is %s \n ", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf(" ==== [clientDo] ReadAll faild ! err is %s \n ", err)
	}
	fmt.Printf("\n %c[1;40;32m coroutine:[%d]; invoke:第 [%v] 次 ; res:%s;  token:%s\n %c[0m\n\n", 0x1B, croutine,i ,string(body), token,0x1B)
	bodymap := make(map[string]interface{})
	err = json.Unmarshal(body,&bodymap)
	if err != nil {
		fmt.Printf("\n %c[1;40;32m  err:%s\n %c[0m\n\n", 0x1B, err,0x1B)
		return fmt.Errorf(" ==== [clientDo] ReadAll faild ! err is %s \n ", err)
	}
	fmt.Printf("\n %c[7;40;36m  coroutine:[%d]; invoke:第 [%v] 次 ;code=[%v] %c[0m\n\n", 0x1B,croutine,i, bodymap["code"],0x1B)
	val , ok := bodymap["code"]
	if !ok {
		panic(fmt.Errorf(" ==== [clientDo] response code faild ! coroutine:[%d] faild !\n ",croutine))
	}
	if val.(float64) != 200 {
		panic(fmt.Errorf(" ==== [clientDo] response code faild ! coroutine:[%d] faild !\n ",croutine))
	}
	defer resp.Body.Close()
	return nil
}

func ChaincodeQuery(croutine int,i int , token string, cert string,invokeChan chan int)  {
	fmt.Printf(" %c[%d;%d;%dm query token:[%s] %c[0m \n", 0x1B, 7, 41, 36, token, 0x1B) // 红色
	currentTime := time.Now().Unix()
	cci :=  &InvokeChaincodeCfg{
		ChannelId: "mychannel1",
		CcName: "mycc4",
		CcFcn: "query",
		CcArgs: []string{"1"},

	}
	cfgBytes , err := json.Marshal(cci)
	if err != nil {
		fmt.Errorf(" ==== chaincodeQuery marshal faild ! err is %s \n ", err)
	}

	err = execute(croutine,i,cfgBytes , token, cert)
	if err != nil {
		fmt.Errorf(" ==== execute clientDo  faild ! err is %s \n", err)
	}

	tm := time.Now().Unix() - currentTime
	fmt.Printf(" %c[%d;%d;%dm Query 第 [%v] 次, 所用时间 [%v] %c[0m \n", 0x1B, 7, 41, 36,  i , tm, 0x1B) // 红色
	//fmt.Printf("=======当前时间 [%v] , Query 第 [%v] 次, 所用时间 [%v] =======\n",  time.Now() , i , tm)
	<- invokeChan
}

func QueryChaincode(i int , token string, cert string,chanquery chan int){
	fmt.Printf(" %c[%d;%d;%dm querycc token:[%s] %c[0m \n", 0x1B, 7, 44, 37, token, 0x1B) // 红色
	currentTime := time.Now().Unix()
	cci :=  &InvokeChaincodeCfg{
		ChannelId: "mychannel1",
		CcName: "mycc4",
		CcFcn: "query",
		CcArgs: []string{"1"},
	}
	cfgBytes , err := json.Marshal(cci)
	if err != nil {
		fmt.Errorf(" ==== chaincodeQuery marshal faild ! err is %s \n ", err)
	}
	// 不会发送事件
	query(cfgBytes , token, cert)

	tm := time.Now().Unix() - currentTime
	fmt.Printf(" %c[%d;%d;%dm querycc QueryCC 第 [%v] 次, 所用时间 [%v]  %c[0m \n", 0x1B, 7, 44, 37, i , tm, 0x1B) // 红色
	//fmt.Printf("=======当前时间 [%v] , QueryCC 第 [%v] 次, 所用时间 [%v] =======\n",  time.Now() , i , tm)
	<- chanquery
}

func query (cfgBytes []byte,token string , cert string){
	if utils.GetHttpsEnable(){
		client , err := HttpsReqClient(cert)
		if err != nil {
			fmt.Errorf(" ==== query HttpsReqClient faild ! err is %s \n ", err)
		}
		url := fmt.Sprintf("https://%s%s/chaincode/query", utils.GetCertCommonName() , utils.GetClientPort())
		clientDo(1, 1,token , client , "POST" , url, cfgBytes  )

	} else {
		client := &http.Client{}
		//req, err := http.NewRequest("POST", "http://fast-deploy:8888/chaincode/execute", body)
		url := fmt.Sprintf("http://%s%s/chaincode/query", utils.GetCertCommonName() , utils.GetClientPort())
		clientDo(1, 1,token , client , "POST" , url, cfgBytes  )

	}
}

func HttpsPost(caCertPath string) {
	pool := x509.NewCertPool()

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}
	client := &http.Client{Transport: tr}
	data := "456"
	resp, err := client.Post("https://unionbank.com:8083/test", "application/json", bytes.NewReader([]byte(data)))
	if err != nil {
		fmt.Println("Get error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func HttpsReqClient(certPath string) (*http.Client ,  error){
	if certPath == "" {
		return nil , errors.New("Reqest https cert can not be empty!")
	}
	pool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile(certPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return nil , err
	}
	pool.AppendCertsFromPEM(caCrt)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: pool},
		DisableKeepAlives: true,
	}
	return &http.Client{Transport: tr} , nil
}


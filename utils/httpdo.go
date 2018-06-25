package utils

import (
	"io/ioutil"
	"fmt"
	"net/http"
	"bytes"
)

func HttpDo(method string , url string , tocken string  , bodyBytes []byte) ([]byte , error){
	client := &http.Client{}
	cibody := bytes.NewBuffer(bodyBytes)
	//req, err := http.NewRequest("POST", "http://fast-deploy:8888/chaincode/execute", body)
	req, err := http.NewRequest(method, url , cibody)
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-Auth-Token", tocken)

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))

	return body , nil
}

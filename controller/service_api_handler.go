package controller

import (
	"getContractDeployment/helper"
	"io"
	"net/http"
	"sync"
)

func callGetAPI(url string) []byte{
    req, err := http.NewRequest("GET", url, nil)
    helper.CheckError(err)

    res, err :=http.DefaultClient.Do(req)
    helper.CheckError(err)

    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    helper.CheckError(err)

    return body
}

func fetchAPIData(wg *sync.WaitGroup, url string ,bodyChan chan []byte){
    body := callGetAPI(url)

    bodyChan <- body
    wg.Done()
}
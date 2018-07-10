package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func coinsGetAccountInfo(token string) (*CoinsGetBody, error) {
	coinsAcctInfo := &CoinsGetBody{}

	req, err := http.NewRequest("GET", COINSURL_GETINFO, nil)
	if err != nil {
		fmt.Println("NewRequest GET fail")
		return nil, err
	}

	req.Header.Add("authorization", strings.Join([]string{"Bearer", token}, " "))
	req.Header.Add("content-type", "application/json")
	req.Header.Add("accept", "application/json")
	req.Header.Add("cache-control", "no-cache")

	log.Printf("Request Header = %#v \n", req.Header)

	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != 200 {
		log.Println("DefaultClient.Do fail")
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("ioutil.ReadAll fail")
		return nil, err
	}
	log.Printf("Response Body as Byte stream = %s\n", string(body))
	err = json.Unmarshal(body, coinsAcctInfo)
	if err != nil {
		log.Printf("Unmarshal fail")
		return nil, err
	}
	log.Printf("coinsAcctInfo  = %#v", coinsAcctInfo)

	return coinsAcctInfo, err
}

func coinsGetPBTCAccount(CryptoAccounts []CryptoAcct) (acct *CryptoAcct, err error) {

	for i := range CryptoAccounts {
		if CryptoAccounts[i].Currency == "PBTC" {
			return &CryptoAccounts[i], nil
		}
	}
	return nil, errors.New("[coinsGetAccountID] error No PBTC acct")
}

func coinsTransferDonation(accountid string, donatedamount float64, walletAddr string, token string) (string, error) {
	transbody := TransferBody{
		accountid,
		donatedamount,
		walletAddr,
	}
	fmt.Printf("Request Body = %#v \n", transbody)
	payload, _ := json.Marshal(transbody)
	payloadreader := bytes.NewReader(payload)

	req, err := http.NewRequest("POST", COINSURL_TRANSMONEY, payloadreader)
	if err != nil {
		log.Printf("Unmarshal fail")
		return "", err
	}

	req.Header.Add("authorization", strings.Join([]string{"Bearer", token}, " "))
	req.Header.Add("content-type", "application/json;charset=UTF-8")
	req.Header.Add("accept", "application/json")
	req.Header.Add("cache-control", "no-cache")
	fmt.Printf("Request Header = %#v \n", req.Header)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("DefaultClient fail: %+v\n", err)
		return "", err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("ReadAll fail: %+v\n", err)
		return "", err
	}

	log.Printf("res: %+v\n", res)
	log.Printf("body: %+v\n", string(body))
	return string(body), nil
}

func createQueryString(params *TransactionParams) (qstring string, err error) {
	//ex: {"selector":{"CharityID":"marble","Status":1}
	var selector = TransactionParamSelector{Selector: *params}
	serialized, err := json.Marshal(selector)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	qstring = string(serialized)
	return qstring, nil
}

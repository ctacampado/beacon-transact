package main

import (
	"encoding/json"
	"fmt"

	shim "github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func getTransactions(fargs CCFuncArgs) pb.Response {
	fmt.Println("starting getTransactions")

	var qparams = &TransactionParams{}
	err := json.Unmarshal([]byte(fargs.req.Params), qparams)
	if err != nil {
		return shim.Error("[getTransactions] Error unable to unmarshall msg: " + err.Error())
	}

	qstring, err := createQueryString(qparams)
	if err != nil {
		return shim.Error("[getTransactions] Error unable to create query string: " + err.Error())
	}

	resultsIterator, err := fargs.stub.GetQueryResult(qstring)
	fmt.Printf("- getQueryResultForQueryString resultsIterator:\n%+v\n", resultsIterator)
	defer resultsIterator.Close()
	if err != nil {
		return shim.Error("[getTransactions] Error unable to GetQueryResult: " + err.Error())
	}

	type qres struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	type qrsp struct {
		Elem []qres `json:"elem"`
	}

	var qresp = qrsp{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("[getTransactions] Error unable to get next item in iterator: " + err.Error())
		}

		q := qres{Key: queryResponse.Key, Value: string(queryResponse.Value)}
		qresp.Elem = append(qresp.Elem, q)
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%+v\n", qresp)
	fmt.Printf("- getQueryResultForQueryString querystring:\n%s\n", qstring)
	fmt.Printf("- getQueryResultForQueryString qparams:\n%+v\n", qparams)

	qr, err := json.Marshal(qresp)
	if err != nil {
		return shim.Error("[getTransactions] Error unable to Marshall qresp: " + err.Error())
	}

	fmt.Println("- end getTransactions")
	return shim.Success(qr)
}

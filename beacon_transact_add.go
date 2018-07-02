package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	shim "github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/satori/go.uuid"
)

func handleDonateTransaction(t *TxInfo, stub shim.ChaincodeStubInterface) pb.Response {
	// Call Coins API
	// Coins GET: retrieve account info
	log.Printf("CALL coinsGetAccountID\n")
	accts, err := coinsGetAccountInfo(t.DonationInfo.CoinsAPIToken)
	if err != nil {
		return shim.Error("[addTransaction] Error unable to retrieve coins account information: " + err.Error())
	}

	// Coins POST: Transfer donation via COINS
	pbtcAcct, err := coinsGetPBTCAccount(accts.Crypto_Accounts)
	if err != nil {
		return shim.Error("[addTransaction] Error unable to retrieve pbtc account: " + err.Error())
	}
	log.Printf("CALL coinsTransferDonation\n")
	fAmt, _ := strconv.ParseFloat(t.DonationInfo.Amount, 64)
	err = coinsTransferDonation(pbtcAcct.ID, fAmt, t.DonationInfo.WalletAddrDst, t.DonationInfo.CoinsAPIToken)
	if err != nil {
		return shim.Error("[addTransaction] Error unable to transfer donation: " + err.Error())
	}

	bytes, err := json.Marshal(t)
	if err != nil {
		log.Printf("[addTransaction] Could not marshal campaign info object: %+v\n", err)
		return shim.Error(err.Error())
	}

	err = stub.PutState(t.TxnID, bytes)
	if err != nil {
		log.Printf("[addTransaction] Error storing data in the ledger %+v\n", err)
		return shim.Error(err.Error())
	}

	return shim.Success(bytes)
}

func handleDisbursementTransaction(i int, t *TxInfo, stub shim.ChaincodeStubInterface) pb.Response {
	type CampaignParams struct {
		CharityID       string `json:"CharityID"`
		CampaignID      string `json:"CampaignID"`
		DisbursedAmount string `json:"DisbursedAmount"`
	}

	cparam := CampaignParams{
		CharityID:       t.DisbursementInfo[i].CharityID,
		CampaignID:      t.DisbursementInfo[i].CampaignID,
		DisbursedAmount: t.DisbursementInfo[i].Price,
	}

	bytes, err := json.Marshal(cparam)
	if err != nil {
		log.Printf("[handleDisbursementTransaction] Could not marshal campaign info object: %+v\n", err)
		return shim.Error(err.Error())
	}

	ccmsg := Message{
		AID:    t.AID,
		Type:   "modifyCampaign",
		Params: string(bytes),
	}

	args, err := json.Marshal(ccmsg)
	if err != nil {
		log.Printf("[handleDisbursementTransaction] Could not marshal campaign info object: %+v\n", err)
		return shim.Error(err.Error())
	}

	var argsarr [][]byte
	argsarr[0] = args

	return stub.InvokeChaincode("beacon_cmpgns", argsarr, "campaigns")
}

func addTransaction(fargs CCFuncArgs) pb.Response {
	log.Printf("starting addTransaction\n")

	u := uuid.Must(uuid.NewV4())
	var TxnID = u.String()

	t := &TxInfo{
		TxnID:   TxnID,
		TxnDate: string(time.Now().Format("2006-Jan-02")),
		Status:  "NEW",
	}

	err := json.Unmarshal([]byte(fargs.req.Params), &t)
	if err != nil {
		return shim.Error("[addTransaction] Error unable to unmarshall msg: " + err.Error())
	}
	log.Printf("[addTransaction ] transaction info: %+v\n", t)

	var rsp pb.Response
	if t.TxnType == "Donation" {
		return handleDonateTransaction(t, fargs.stub)
	} else if t.TxnType == "Disbursement" {
		for i, elem := range t.DisbursementInfo {
			log.Printf("[addTransaction ] elem info: %+v\n", elem)
			rsp = handleDisbursementTransaction(i, t, fargs.stub)
			if nil != rsp.GetPayload() {
				return rsp
			}
		}
	}

	log.Println("- end addTransaction")
	return shim.Error("[addTransaction] error unknown transaction type") //change nil to appropriate response
}

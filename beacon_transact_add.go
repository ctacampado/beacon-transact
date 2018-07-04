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
	fAmt, err := strconv.ParseFloat(t.DonationInfo.Amount, 64)
	if err != nil {
		return shim.Error("[addTransaction] Error converting Amount to type float from string: " + err.Error())
	}
	log.Printf("famt: %f\n", fAmt)

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

	type CampaignParams struct {
		CharityID     string `json:"CharityID"`
		CampaignID    string `json:"CampaignID"`
		DonatedAmount string `json:"DonatedAmount"`
	}
	cparam := CampaignParams{
		CharityID:     t.DonationInfo.CharityID,
		CampaignID:    t.DonationInfo.CampaignID,
		DonatedAmount: t.DonationInfo.Amount,
	}

	bytes, err = json.Marshal(cparam)
	if err != nil {
		log.Printf("[handleDonationTransaction] Could not marshal campaign info object: %+v\n", err)
		return shim.Error(err.Error())
	}

	ccmsg := Message{
		AID:    t.AID,
		Type:   "modifyCampaign",
		Params: string(bytes),
	}

	args, err := json.Marshal(ccmsg)
	if err != nil {
		log.Printf("[handleDonationTransaction] Could not marshal campaign info object: %+v\n", err)
		return shim.Error(err.Error())
	}

	argsarr := [][]byte{}
	argsarr = append(argsarr, []byte("modifyCampaign"))
	argsarr = append(argsarr, args)
	log.Printf("argsarr: %+v\n", argsarr)

	return stub.InvokeChaincode("cmpgnscc", argsarr, "campaigns")
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

	argsarr := [][]byte{}
	argsarr = append(argsarr, args)
	log.Printf("argsarr: %+v\n", argsarr)

	return stub.InvokeChaincode("cmpgnscc", argsarr, "campaigns")
}

func addTransaction(fargs CCFuncArgs) pb.Response {
	log.Printf("starting addTransaction\n")
	log.Printf("Param: %+v\n", fargs.req.Params)
	u := uuid.Must(uuid.NewV4())
	var TxnID = u.String()

	t := TxInfo{
		TxnID:   TxnID,
		TxnDate: string(time.Now().Format("2006-Jan-02")),
	}
	err := json.Unmarshal([]byte(fargs.req.Params), &t)
	if err != nil {
		return shim.Error("[addTransaction] Error unable to unmarshall msg: " + err.Error())
	}
	log.Printf("[addTransaction] DonationInfo: %+v\n", t.DonationInfo)
	log.Printf("[addTransaction] transaction info: %+v\n", t)

	var rsp pb.Response
	if t.TxnType == "Donation" {
		return handleDonateTransaction(&t, fargs.stub)
	} else if t.TxnType == "Disbursement" {
		for i, elem := range t.DisbursementInfo {
			log.Printf("[addTransaction ] elem info: %+v\n", elem)
			rsp = handleDisbursementTransaction(i, &t, fargs.stub)
			if nil != rsp.GetPayload() {
				return rsp
			}
		}
	}

	log.Println("- end addTransaction")
	return shim.Error("[addTransaction] error unknown transaction type") //change nil to appropriate response
}

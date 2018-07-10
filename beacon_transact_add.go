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

func handleDonateTransaction(t *TxInfo, fargs CCFuncArgs) error {
	// Call Coins API
	// Coins GET: retrieve account info
	log.Printf("CALL coinsGetAccountID\n")
	accts, err := coinsGetAccountInfo(t.DonationInfo.CoinsAPIToken)
	if err != nil {
		log.Printf("[addTransaction] Error unable to retrieve coins account information: " + err.Error())
		return err
	}

	// Coins POST: Transfer donation via COINS
	pbtcAcct, err := coinsGetPBTCAccount(accts.Crypto_Accounts)
	if err != nil {
		log.Printf("[addTransaction] Error unable to retrieve pbtc account: " + err.Error())
		return err
	}

	log.Printf("CALL coinsTransferDonation\n")
	fAmt, err := strconv.ParseFloat(t.DonationInfo.Amount, 64)
	if err != nil {
		log.Printf("[addTransaction] Error converting Amount to type float from string: " + err.Error())
		return err
	}
	log.Printf("famt: %f\n", fAmt)

	txnres, err := coinsTransferDonation(pbtcAcct.ID, fAmt, t.DonationInfo.WalletAddrDst, t.DonationInfo.CoinsAPIToken)
	if err != nil {
		log.Printf("[addTransaction] Error unable to transfer donation: " + err.Error())
		return err

	}
	t.TxnRes = txnres
	bytes, err := json.Marshal(t)
	if err != nil {
		log.Printf("[addTransaction] Could not marshal campaign info object: %+v\n", err)
		return err
	}

	err = fargs.stub.PutState(t.TxnID, bytes)
	if err != nil {
		log.Printf("[addTransaction] Error storing data in the ledger %+v\n", err)
		return err
	}

	fargs.msg.Data = string(bytes)
	log.Printf("fargs: %+v\n", fargs.msg)
	rspbytes, err := json.Marshal(fargs.msg)
	if err != nil {
		log.Printf("[handleDonation] Could not marshal fargs object: %+v\n", err)
		return err
	}
	log.Printf("rspbytes: %+v\n", rspbytes)
	log.Println("- end handleDonation")
	fargs.stub.SetEvent("donatetxn", rspbytes)
	return nil
}

func handleDisbursementTransaction(i int, t *TxInfo, fargs CCFuncArgs) error {
	bytes, err := json.Marshal(t.DisbursementInfo[i])
	if err != nil {
		log.Printf("[handleDisbursementTransaction] Could not marshal campaign info object: %+v\n", err)
		return err
	}

	err = fargs.stub.PutState(t.TxnID, bytes)
	if err != nil {
		log.Printf("[handleDisbursement] Error storing data in the ledger %+v\n", err)
		return err
	}
	fargs.msg.Data = string(bytes)
	log.Printf("fargs: %+v\n", fargs.msg)
	rspbytes, err := json.Marshal(fargs.msg)
	if err != nil {
		log.Printf("[handleDisbursement] Could not marshal fargs object: %+v\n", err)
		return err
	}
	log.Printf("rspbytes: %+v\n", rspbytes)
	log.Println("- end handleDisbursement")
	fargs.stub.SetEvent("disbursetxn", rspbytes)
	return nil
}

func addTransaction(fargs CCFuncArgs) pb.Response {
	log.Printf("starting addTransaction\n")
	log.Printf("Param: %+v\n", fargs.msg.Params)
	u := uuid.Must(uuid.NewV4())
	var TxnID = u.String()

	t := TxInfo{
		TxnID:   TxnID,
		TxnDate: string(time.Now().Format("2006-Jan-02")),
	}
	err := json.Unmarshal([]byte(fargs.msg.Params), &t)
	if err != nil {
		return shim.Error("[addTransaction] Error unable to unmarshall msg: " + err.Error())
	}
	log.Printf("[addTransaction] DonationInfo: %+v\n", t.DonationInfo)
	log.Printf("[addTransaction] transaction info: %+v\n", t)

	if t.DonationInfo.WalledAddrDst != "" {
		if t.TxnType == "Donation" {
			err = handleDonateTransaction(&t, fargs)
			if nil != err {
				return shim.Error("[addTransaction] error handleDonateTransaction") //change nil to appropriate response
			}
		} else if t.TxnType == "Disbursement" {
			for i, elem := range t.DisbursementInfo {
				log.Printf("[addTransaction ] elem info: %+v\n", elem)
				err = handleDisbursementTransaction(i, &t, fargs)
				if nil != err {
					return shim.Error("[addTransaction] error handleDisbursementTransaction") //change nil to appropriate response
				}
			}
		}
	}

	log.Println("- end addTransaction")
	return shim.Success(nil)
}

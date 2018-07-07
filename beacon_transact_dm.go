package main

import (
	shim "github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//--------------------------------------------------------------------------
//Start adding Chaincode-related Structures here

//CCFuncArgs common cc func args
type CCFuncArgs struct {
	function string
	msg      Message
	stub     shim.ChaincodeStubInterface
}

type ccfunc func(args CCFuncArgs) pb.Response

//Chaincode cc structure
type Chaincode struct {
	FMap map[string]ccfunc //ccfunc map
	Msg  Message           //data
}

//Message Charity Org Chain Code Message Structure
type Message struct {
	CID    string `json:"CID, omitempty"` //ClientID --for websocket push (event-based messaging readyness)
	AID    string `json:"AID"`            //ActorID (Donor ID/Charity Org ID/Auditor ID/etc.)
	Type   string `json:"type"`           //Chaincode Function
	Params string `json:"params"`         //Function Parameters
	Data   string `json:"data,omitempty"`
}

//End of Chaincode-related Structures
//--------------------------------------------------------------------------
//Start adding Query Parameter (Parm) Structures here

//TransactionParams Structure for Query Parameters
type TransactionParams struct {
	TxnID        string `json:"TxnID,omitempty"`
	TxnType      string `json:"TxnType,omitempty"`
	AID          string `json:"AID,omitempty"`
	TxnDate      string `json:"TxnDate,omitempty"`
	DonationInfo struct {
		WalletAddrSrc string `json:"WalletAddrSrc,omitempty"`
		WalletAddrDst string `json:"WalletAddrDst,omitempty"`
		CharityID     string `json:"CharityID,omitempty"`
		CampaignID    string `json:"CampaignID,omitempty"`
		Amount        string `json:"Amount,omitempty"`
		CoinsAPIToken string `json:"CoinsAPIToken,omitempty"`
	} `json:"DonationInfo,omitempty"`
	DisbursementInfo []DisbursementInfo `json:"DisbursementInfo,omitempty"`
}

//TransactionParamSelector Structure for Query Selector
type TransactionParamSelector struct {
	Selector TransactionParams `json:"selector"`
}

//End of Query Paramter Structures
//--------------------------------------------------------------------------
//Start adding Data Models here

type DisbursementInfo struct {
	CharityID      string `json:"CharityID"`
	CampaignID     string `json:"CampaignID"`
	Particular     string `json:"Particular"`
	QtyParticular  int    `json:"QtyParticular"`
	UnitParticular string `json:"Unitparticular"`
	Price          string `json:"Price"`
	Date           string `json:"Date"`
}

type TxInfo struct {
	TxnID            string             `json:"TxnID,omitempty"`
	TxnType          string             `json:"TxnType"`
	AID              string             `json:"AID"`
	TxnDate          string             `json:"TxnDate,omitempty"`
	TxnRes           string             `json:"TxnRes,omitempty"`
	DonationInfo     DonationInfo       `json:"DonationInfo,omitempty"`
	DisbursementInfo []DisbursementInfo `json:"DisbursementInfo,omitempty"`
}

type DonationInfo struct {
	WalletAddrSrc string `json:"WalletAddrSrc"`
	WalletAddrDst string `json:"WalletAddrDst"`
	CharityID     string `json:"CharityID"`
	CampaignID    string `json:"CampaignID"`
	Amount        string `json:"Amount"`
	CoinsAPIToken string `json:"CoinsAPIToken"`
}

//End of Data Models
//--------------------------------------------------------------------------

//Coins Data Models
const COINSURL_TRANSMONEY = "https://coins.ph/api/v3/transfers/"
const COINSURL_GETINFO = "https://coins.ph/api/v3/crypto-accounts/"

type CryptoAcct struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Currency        string `json:"currency"`
	Balance         string `json:"balance"`
	Pending_Balance string `json:"pending_balance"`
	Total_Received  string `json:"total_received"`
	Default_Address string `json:"default_address"`
	Is_Default      bool   `json:"is_default"`
}
type CryptoMeta struct {
	Total_Count   int    `json:"total_count"`
	Next_Page     string `json:"next_page"`
	Previous_Page string `json:"previous_page"`
}
type CoinsGetBody struct {
	Meta            CryptoMeta   `json:"meta"`
	Crypto_Accounts []CryptoAcct `json:"crypto-accounts"`
}
type TransferBody struct {
	Account        string  `json:"account"`
	Amount         float64 `json:"amount"`
	Target_Address string  `json:"target_address"`
}

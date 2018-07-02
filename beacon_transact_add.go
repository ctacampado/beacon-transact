package main

import (
	"encoding/json"
	"fmt"
	"log"

	shim "github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/satori/go.uuid"
)

func addTransaction(fargs CCFuncArgs) pb.Response {
	log.Printf("starting addTransaction\n")

	u := uuid.Must(uuid.NewV4())
	var campaignID = u.String()

	c := CampaignInfo{CampaignID: campaignID}

	err := json.Unmarshal([]byte(fargs.req.Params), &c)
	if err != nil {
		return shim.Error("[addTransaction] Error unable to unmarshall msg: " + err.Error())
	}

	c.Status = 1
	c.DonatedAmount = "0"
	c.TransAmount = "0"
	c.CampCompDate = "n/a"
	c.RatingFive = "0"
	c.RatingFour = "0"
	c.RatingThree = "0"
	c.RatingTwo = "0"
	c.RatingOne = "0"

	log.Printf("[addTransaction ] campaign info: %+v\n", c)

	bytes, err := json.Marshal(c)
	if err != nil {
		log.Printf("[addTransaction] Could not marshal campaign info object: %+v\n", err)
		return shim.Error(err.Error())
	}

	err = fargs.stub.PutState(campaignID, bytes)
	if err != nil {
		log.Printf("[addTransaction] Error storing data in the ledger %+v\n", err)
		return shim.Error(err.Error())
	}

	fmt.Println("- end addTransaction")
	return shim.Success(nil) //change nil to appropriate response
}

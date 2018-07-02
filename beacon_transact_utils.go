package main

import (
	"encoding/json"
	"fmt"
)

func createQueryString(params *CampaignQueryParams) (qstring string, err error) {
	//ex: {"selector":{"CharityID":"marble","Status":1}
	var selector = CampaignQuerySelector{Selector: *params}
	serialized, err := json.Marshal(selector)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	qstring = string(serialized)
	return qstring, nil
}

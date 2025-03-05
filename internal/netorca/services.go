// Copyright (c) HashiCorp, Inc.

package netorca

import (
	"fmt"
	"net/http"
)

type NetOrcaService struct {
	Id                    int
	Name                  string
	Owner                 NetOrcaOwner
	ApprovalRequired      bool
	AllowManualApproval   bool
	AllowManualCompletion bool
	Schema                interface{}
}

type NetOrcaOwner struct {
	Name string
	Id   int
}

func (c *NetOrcaClient) ServiceGet() NetOrcaService {
	url := fmt.Sprintf("%s/v1/orcabase/consumer/services", c.baseUrl)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", c.apiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return NetOrcaService{}
}

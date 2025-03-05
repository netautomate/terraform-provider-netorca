// Copyright (c) HashiCorp, Inc.

package netorca

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ChangeInstance struct {
	Id               int64                      `json:"id"`
	Url              string                     `json:"url"`
	State            string                     `json:"state"`
	Created          string                     `json:"created"`
	Modified         string                     `json:"modified"`
	Owner            ChangeInstanceOwner        `json:"owner"`
	ConsumerTeam     ChangeInstanceConsumerTeam `json:"consumer_team"`
	ServiceOwnerTeam interface{}                `json:"service_owner_team"`
	Submission       ChangeInstanceSubmission   `json:"submission"`
	ServiceItemField ServiceItem                `json:"service_item"`
}

type ChangeInstanceSubmission struct {
	Id       int64  `json:"id"`
	CommitId string `json:"commit_id"`
}

type ChangeInstanceConsumerTeam struct {
	Id       int64       `json:"id"`
	Name     string      `json:"name"`
	Metadata interface{} `json:"metadata"`
}

type ChangeInstanceOwner struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type NetOrcaChangeInstance struct {
	Count    int
	Next     string
	Previous string
	Results  []ChangeInstance
}

type ChangeInstanceQuery struct {
	Pov                string
	ApplicationId      int64
	ChangeType         string
	CommitId           string
	ConsumerTeamId     int64
	Limit              int64
	Offset             int64
	Ordering           string
	ServiceId          int64
	ServiceItemId      int64
	ServiceName        string
	ServiceOwnerTeamId int64
	State              string
	SubmissionId       int64
}

type ChangeInstanceUpdateRequest struct {
	State        string `json:"state"`
	Description  string `json:"description"`
	DeployedItem string `json:"deployed_item"`
}

type ChangeInstanceUpdateJson struct {
	State        string                 `json:"state"`
	Description  string                 `json:"description"`
	DeployedItem map[string]interface{} `json:"deployed_item"`
}

func (c *NetOrcaClient) ChangeInstancePatch(id int64, pov string, request ChangeInstanceUpdateRequest) error {
	url := fmt.Sprintf("%s/v1/orcabase/%s/change_instances/%d/", c.baseUrl, pov, id)
	var deployedItem map[string]interface{}

	err := json.Unmarshal([]byte(request.DeployedItem), &deployedItem)
	if err != nil {
		return err
	}

	content := ChangeInstanceUpdateJson{
		State:        request.State,
		Description:  request.Description,
		DeployedItem: deployedItem,
	}

	json, err := json.Marshal(content)
	if err != nil {
		return err
	}

	serv, err := http.NewRequest("PATCH", url, bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	serv.Header.Add("Authorization", c.GetApiKey())
	serv.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(serv)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("http code: %d\nresponse: %s\nurl: %s\nmethod: PATCH", resp.StatusCode, b, url)
	}

	return nil

}

func (c *NetOrcaClient) ChangeInstanceGet(q *ChangeInstanceQuery) (NetOrcaChangeInstance, error) {

	url := fmt.Sprintf("%s/v1/orcabase/%s/change_instances/", c.baseUrl, q.Pov)

	queryParameters := q.GetQueryParam()

	if queryParameters != "" {
		url = fmt.Sprintf("%s%s", url, queryParameters)
	}

	serv, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NetOrcaChangeInstance{}, err
	}

	serv.Header.Add("Authorization", c.GetApiKey())

	resp, err := c.client.Do(serv)
	if err != nil {
		return NetOrcaChangeInstance{}, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return NetOrcaChangeInstance{}, err
	}

	if resp.StatusCode != 200 {
		return NetOrcaChangeInstance{}, fmt.Errorf("http code: %d\n response: %s\nurl: %s\nmethod: GET", resp.StatusCode, b, url)
	}

	var changeInstances NetOrcaChangeInstance

	err = json.Unmarshal(b, &changeInstances)
	if err != nil {
		return NetOrcaChangeInstance{}, err
	}

	return changeInstances, nil
}

func (c *NetOrcaClient) ChangeInstanceGetById(id int64, pov string) (ChangeInstance, error) {

	url := fmt.Sprintf("%s/v1/orcabase/%s/change_instances/%d/", c.baseUrl, pov, id)

	serv, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ChangeInstance{}, err
	}

	serv.Header.Add("Authorization", c.GetApiKey())

	resp, err := c.client.Do(serv)
	if err != nil {
		return ChangeInstance{}, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChangeInstance{}, err
	}

	if resp.StatusCode != 200 {
		return ChangeInstance{}, fmt.Errorf("http code: %d\n response: %s\nurl: %s\nmethod: GET", resp.StatusCode, b, url)
	}

	var changeInstance ChangeInstance

	err = json.Unmarshal(b, &changeInstance)
	if err != nil {
		return ChangeInstance{}, err
	}

	return changeInstance, nil
}

// Returns a *ChangeInstanceQuery or nil and an error message if one of the type inferences aren't handled.
func NewChangeInstanceQuery(args map[string]interface{}) (*ChangeInstanceQuery, error) {
	c := ChangeInstanceQuery{}

	for k, v := range args {
		switch k {
		case "pov":
			if str, ok := v.(string); ok {
				c.Pov = str
			} else {
				return nil, fmt.Errorf("pov not passed as a string")
			}
		case "application_id":
			if i, ok := v.(int64); ok {
				c.ApplicationId = i
			} else {
				return nil, fmt.Errorf("application_id not passed as a uint64")
			}
		case "change_type":
			if str, ok := v.(string); ok {
				c.ChangeType = str
			} else {
				return nil, fmt.Errorf("change_type not passed as a string")
			}
		case "commit_id":
			if str, ok := v.(string); ok {
				c.CommitId = str
			} else {
				return nil, fmt.Errorf("commit_id not passed as a string")
			}
		case "consumer_team_id":
			if i, ok := v.(int64); ok {
				c.ConsumerTeamId = i
			} else {
				return nil, fmt.Errorf("consumer_team_id not passed as a uint64")
			}
		case "limit":
			if i, ok := v.(int64); ok {
				c.Limit = i
			} else {
				return nil, fmt.Errorf("limit not passed as a uint64")
			}
		case "offset":
			if i, ok := v.(int64); ok {
				c.Offset = i
			} else {
				return nil, fmt.Errorf("offset not passed as a uint64")
			}
		case "ordering":
			if str, ok := v.(string); ok {
				c.Ordering = str
			} else {
				return nil, fmt.Errorf("ordering not passed as a string")
			}
		case "service_id":
			if i, ok := v.(int64); ok {
				c.ServiceId = i
			} else {
				return nil, fmt.Errorf("service_id not passed as a uint64")
			}
		case "service_item_id":
			if i, ok := v.(int64); ok {
				c.ServiceItemId = i
			} else {
				return nil, fmt.Errorf("service_item_id not passed as a uint64")
			}
		case "service_name":
			if str, ok := v.(string); ok {
				c.ServiceName = str
			} else {
				return nil, fmt.Errorf("service_name not passed as a string")
			}
		case "service_owner_team_id":
			if i, ok := v.(int64); ok {
				c.ServiceOwnerTeamId = i
			} else {
				return nil, fmt.Errorf("service_owner_team_id not passed as a uint64")
			}
		case "state":
			if str, ok := v.(string); ok {
				c.State = str
			} else {
				return nil, fmt.Errorf("state not passed as a string")
			}
		case "submission_id":
			if i, ok := v.(int64); ok {
				c.SubmissionId = i
			} else {
				return nil, fmt.Errorf("submission_id not passed as a uint64")
			}
		}
	}

	return &c, nil
}

// Returns the formatted query parmaters for use with the http client.
// e.g. in the form of ?<field_name>=<field_value>&<field_name>=<field_value>
func (q ChangeInstanceQuery) GetQueryParam() string {
	queryParam := "?"

	if q.ServiceId != 0 {
		queryParam = fmt.Sprintf("%sservice_id=%d&", queryParam, q.ServiceId)
	}

	if q.ServiceItemId != 0 {
		queryParam = fmt.Sprintf("%sservice_item_id=%d&", queryParam, q.ServiceItemId)
	}

	if q.ServiceName != "" {
		queryParam = fmt.Sprintf("%sservice_name=%s&", queryParam, q.ServiceName)
	}

	if q.ServiceOwnerTeamId != 0 {
		queryParam = fmt.Sprintf("%sservice_owner_team_id=%d&", queryParam, q.ServiceOwnerTeamId)
	}

	if q.State != "" {
		queryParam = fmt.Sprintf("%sstate=%s&", queryParam, q.State)
	}

	if q.SubmissionId != 0 {
		queryParam = fmt.Sprintf("%ssubmission_id=%d&", queryParam, q.SubmissionId)
	}

	// Remove the trailing '&' if it exists
	if queryParam[len(queryParam)-1] == '&' {
		queryParam = queryParam[:len(queryParam)-1]
	}

	// If only '?' remains, return an empty string
	if queryParam == "?" {
		return ""
	}

	return queryParam
}

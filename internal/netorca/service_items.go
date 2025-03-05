// Copyright (c) HashiCorp, Inc.

package netorca

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

type ServiceItem struct {
	Id                int64                       `json:"id"`
	Url               string                      `json:"url"`
	Name              string                      `json:"name"`
	Created           string                      `json:"created"`
	Modified          string                      `json:"modified"`
	RuntimeState      string                      `json:"runtime_state"`
	ServiceName       string                      `json:"service_name"`
	ChangeState       string                      `json:"change_state"`
	Service           ServiceItemService          `json:"service"`
	Application       NetOrcaApplication          `json:"application"`
	DeployedItem      map[string]interface{}      `json:"deployed_item"`
	ConsumerTeam      ServiceItemConsumerTeam     `json:"consumer_team"`
	ServiceOwnerTeam  ServiceItemServiceOwnerTeam `json:"service_owner_team"`
	Declaration       map[string]interface{}      `json:"declaration"`
	Related           interface{}                 `json:"related"`
	HealthcheckStatus *int64                      `json:"healthcheck_status"`
}

type ServiceItemServiceOwnerTeam struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ServiceItemConsumerTeam struct {
	Id       int64       `json:"id"`
	Name     string      `json:"name"`
	Metadata interface{} `json:"metadata"`
}

type ServiceItemService struct {
	Id          int64            `json:"id"`
	Name        string           `json:"name"`
	Owner       ServiceItemOwner `json:"owner"`
	HealthCheck bool             `json:"healthcheck"`
}

type ServiceItemOwner struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ServiceItemQuery struct {
	Pov                string
	ApplicationId      int64
	ChangeState        string
	ConsumerTeamId     int64
	Limit              int64
	Name               string
	Offset             int64
	Ordering           string
	RuntimeState       string
	ServiceName        string
	ServiceId          int64
	ServiceOwnerId     int64
	ServiceOwnerTeamId int64
}

type NetOrcaApplication struct {
	Id       int64       `json:"id"`
	Name     string      `json:"name"`
	Metadata interface{} `json:"metadata"`
	Owner    int64       `json:"owner"`
}

type NetOrcaServiceItem struct {
	Count    int
	Next     string
	Previous string
	Results  []ServiceItem
}

func (c *NetOrcaClient) ServiceItemsGet(s *ServiceItemQuery) (NetOrcaServiceItem, error) {

	url := fmt.Sprintf("%s/v1/orcabase/%s/service_items/", c.baseUrl, s.Pov)

	queryParameters := s.GetQueryParam()
	if queryParameters != "" {
		url = fmt.Sprintf("%s%s", url, queryParameters)
	}

	serv, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NetOrcaServiceItem{}, err
	}

	serv.Header.Add("Authorization", c.GetApiKey())

	resp, err := c.client.Do(serv)
	if err != nil {
		return NetOrcaServiceItem{}, err
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return NetOrcaServiceItem{}, err
	}

	if resp.StatusCode != 200 {
		return NetOrcaServiceItem{}, fmt.Errorf("http code: %d\nresponse: %s\nurl: %s", resp.StatusCode, b, url)
	}

	var serviceItems NetOrcaServiceItem

	err = json.Unmarshal(b, &serviceItems)

	if err != nil {
		return NetOrcaServiceItem{}, err
	}

	return serviceItems, nil
}

// Returns a *ServiceItemQuery or nil and an error message if one of the type inferences aren't handled.
func NewServiceItemQuery(args map[string]interface{}) (*ServiceItemQuery, error) {
	s := ServiceItemQuery{}

	for k, v := range args {
		switch k {
		case "pov":
			if str, ok := v.(string); ok {
				s.Pov = str
			} else {
				return nil, fmt.Errorf("pov not passed as a string")
			}
		case "change_state":
			if str, ok := v.(string); ok {
				s.ChangeState = str
			} else {
				return nil, fmt.Errorf("change_state not passed as a string")
			}
		case "name":
			if str, ok := v.(string); ok {
				s.Name = str
			} else {
				return nil, fmt.Errorf("name not passed as a string")
			}
		case "application_id":
			i, ok := v.(int64)
			if ok && i >= 0 {
				s.ApplicationId = i
			} else {
				return nil, fmt.Errorf("application_id not passed as an uint64")
			}
		case "consumer_team_id":
			i, ok := v.(int64)
			if ok && i >= 0 {
				s.ConsumerTeamId = i
			} else {
				return nil, fmt.Errorf("consumer_team_id not passed as an uint64")
			}
		case "limit":
			i, ok := v.(int64)
			if ok && i >= 0 {
				s.Limit = i
			} else {
				return nil, fmt.Errorf("limit not passed as an uint64")
			}
		case "offset":
			i, ok := v.(int64)
			if ok && i >= 0 {
				s.Offset = i
			} else {
				return nil, fmt.Errorf("offset not passed as an uint64")
			}
		case "ordering":
			if str, ok := v.(string); ok {
				s.Ordering = str
			} else {
				return nil, fmt.Errorf("ordering not passed as a string")
			}
		case "runtime_state":
			if str, ok := v.(string); ok {
				s.RuntimeState = str
			} else {
				return nil, fmt.Errorf("runtime_state not passed as a string")
			}
		case "service_name":
			if str, ok := v.(string); ok {
				s.ServiceName = str
			} else {
				return nil, fmt.Errorf("service_name not passed as a string")
			}
		case "service_owner_id":
			i, ok := v.(int64)
			if ok && i >= 0 {
				s.ServiceOwnerId = i
			} else {
				return nil, fmt.Errorf("service_owner_id not passed as an uint64")
			}
		case "service_owner_team_id":
			i, ok := v.(int64)
			if ok && i >= 0 {
				s.ServiceOwnerTeamId = i
			} else {
				return nil, fmt.Errorf("service_owner_team_id not passed as an uint64")
			}
		}
	}

	return &s, nil
}

// Small helper function to check if all of the non-mandatory fields are zero values.
func isServiceItemQueryZeroValue(s ServiceItemQuery) bool {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		// First field (index 0) is POV which is a mandatory field.
		if i != 0 && !field.IsZero() {
			return false
		}
	}
	return true
}

// Returns the formatted query parmaters for use with the http client.
// e.g. in the form of ?<field_name>=<field_value>&<field_name>=<field_value>
// NOTE: The ordering matters here, the query params are crafted in the order evaluated
// in this function. For example if we have the following map used when calling NewServiceItemQuery()
//
//	map[string]interface{}{
//		"pov":                   "serviceowner",
//		"change_state":          "RUNNING",
//		"application_id":        int64(123),
//	}
//
// The query params would be rendered as: ?&application_id=123&change_state=RUNNING
// As ApplicationId is evaluated first below.
func (q *ServiceItemQuery) GetQueryParam() string {
	var queryParam string

	if isServiceItemQueryZeroValue(*q) {
		return ""
	}

	queryParam = "?"

	if q.ApplicationId != 0 {
		queryParam = fmt.Sprintf("%s&application_id=%d", queryParam, q.ApplicationId)
	}

	if q.ChangeState != "" {
		queryParam = fmt.Sprintf("%s&change_state=%s", queryParam, q.ChangeState)
	}

	if q.ConsumerTeamId != 0 {
		queryParam = fmt.Sprintf("%s&consumer_team_id=%d", queryParam, q.ConsumerTeamId)
	}

	if q.Limit != 0 {
		queryParam = fmt.Sprintf("%s&limit=%d", queryParam, q.Limit)
	}

	if q.Name != "" {
		queryParam = fmt.Sprintf("%s&name=%s", queryParam, q.Name)
	}

	if q.Offset != 0 {
		queryParam = fmt.Sprintf("%s&offset=%d", queryParam, q.Offset)
	}

	if q.Ordering != "" {
		queryParam = fmt.Sprintf("%s&ordering=%s", queryParam, q.Ordering)
	}

	if q.RuntimeState != "" {
		queryParam = fmt.Sprintf("%s&runtime_state=%s", queryParam, q.RuntimeState)
	}

	if q.ServiceName != "" {
		queryParam = fmt.Sprintf("%s&service_name=%s", queryParam, q.ServiceName)
	}

	if q.ServiceOwnerId != 0 {
		queryParam = fmt.Sprintf("%s&service_owner_id=%d", queryParam, q.ServiceOwnerId)
	}

	if q.ServiceOwnerTeamId != 0 {
		queryParam = fmt.Sprintf("%s&service_owner_team_id=%d", queryParam, q.ServiceOwnerTeamId)
	}

	return queryParam

}

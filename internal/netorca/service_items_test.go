// Copyright (c) HashiCorp, Inc.

package netorca

import (
	"reflect"
	"testing"
)

func TestNewServiceItemQuery(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		expected *ServiceItemQuery
		errMsg   string
	}{
		{
			name: "valid_case_fully_populated",
			args: map[string]interface{}{
				"pov":                   "consumer",
				"change_state":          "CHANGES_APPROVED",
				"name":                  "test",
				"application_id":        int64(123),
				"consumer_team_id":      int64(456),
				"limit":                 int64(10),
				"offset":                int64(20),
				"ordering":              "name",
				"runtime_state":         "IN_SERVICE",
				"service_owner_id":      int64(789),
				"service_owner_team_id": int64(101),
			},
			expected: &ServiceItemQuery{
				Pov:                "consumer",
				ChangeState:        "CHANGES_APPROVED",
				Name:               "test",
				ApplicationId:      int64(123),
				ConsumerTeamId:     int64(456),
				Limit:              int64(10),
				Offset:             int64(20),
				Ordering:           "name",
				RuntimeState:       "IN_SERVICE",
				ServiceOwnerId:     int64(789),
				ServiceOwnerTeamId: int64(101),
			},
			errMsg: "",
		},
		{
			name: "valid_case_partially_populated",
			args: map[string]interface{}{
				"pov": "serviceowner",
			},
			expected: &ServiceItemQuery{
				Pov:                "serviceowner",
				ChangeState:        "",
				Name:               "",
				ApplicationId:      int64(0),
				ConsumerTeamId:     int64(0),
				Limit:              int64(0),
				Offset:             int64(0),
				Ordering:           "",
				RuntimeState:       "",
				ServiceOwnerId:     int64(0),
				ServiceOwnerTeamId: int64(0),
			},
			errMsg: "",
		},
		{
			name: "invalid_type_pov",
			args: map[string]interface{}{
				"pov": int64(20),
			},
			expected: nil,
			errMsg:   "pov not passed as a string",
		},
		{
			name: "invalid_type_change_state",
			args: map[string]interface{}{
				"change_state": int64(20),
			},
			expected: nil,
			errMsg:   "change_state not passed as a string",
		},
		{
			name: "invalid_type_name",
			args: map[string]interface{}{
				"name": int64(20),
			},
			expected: nil,
			errMsg:   "name not passed as a string",
		},
		{
			name: "invalid_type_string_application_id",
			args: map[string]interface{}{
				"application_id": "string",
			},
			expected: nil,
			errMsg:   "application_id not passed as an uint64",
		},
		{
			name: "invalid_type_string_consumer_team_id",
			args: map[string]interface{}{
				"consumer_team_id": "string",
			},
			expected: nil,
			errMsg:   "consumer_team_id not passed as an uint64",
		},
		{
			name: "invalid_type_string_limit",
			args: map[string]interface{}{
				"limit": "string",
			},
			expected: nil,
			errMsg:   "limit not passed as an uint64",
		},
		{
			name: "invalid_type_string_offset",
			args: map[string]interface{}{
				"offset": "string",
			},
			expected: nil,
			errMsg:   "offset not passed as an uint64",
		},
		{
			name: "invalid_type_ordering",
			args: map[string]interface{}{
				"ordering": int64(1),
			},
			expected: nil,
			errMsg:   "ordering not passed as a string",
		},
		{
			name: "invalid_type_runtime_state",
			args: map[string]interface{}{
				"runtime_state": int64(1),
			},
			expected: nil,
			errMsg:   "runtime_state not passed as a string",
		},
		{
			name: "invalid_type_string_service_owner_id",
			args: map[string]interface{}{
				"service_owner_id": "",
			},
			expected: nil,
			errMsg:   "service_owner_id not passed as an uint64",
		},
		{
			name: "invalid_type_string_service_owner_team_id",
			args: map[string]interface{}{
				"service_owner_team_id": "",
			},
			expected: nil,
			errMsg:   "service_owner_team_id not passed as an uint64",
		},
		{
			name: "invalid_type_negative_application_id",
			args: map[string]interface{}{
				"application_id": int64(-100),
			},
			expected: nil,
			errMsg:   "application_id not passed as an uint64",
		},
		{
			name: "invalid_type_negative_consumer_team_id",
			args: map[string]interface{}{
				"consumer_team_id": int64(-100),
			},
			expected: nil,
			errMsg:   "consumer_team_id not passed as an uint64",
		},
		{
			name: "invalid_type_negative_limit",
			args: map[string]interface{}{
				"limit": int64(-100),
			},
			expected: nil,
			errMsg:   "limit not passed as an uint64",
		},
		{
			name: "invalid_type_negative_offset",
			args: map[string]interface{}{
				"offset": int64(-100),
			},
			expected: nil,
			errMsg:   "offset not passed as an uint64",
		},
		{
			name: "invalid_type_negative_service_owner_team_id",
			args: map[string]interface{}{
				"service_owner_team_id": int64(-100),
			},
			expected: nil,
			errMsg:   "service_owner_team_id not passed as an uint64",
		},
		{
			name: "invalid_type_string_service_owner_team_id",
			args: map[string]interface{}{
				"service_owner_team_id": int64(-100),
			},
			expected: nil,
			errMsg:   "service_owner_team_id not passed as an uint64",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewServiceItemQuery(test.args)

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected: %v, Got: %v", test.expected, result)
			}

			if err != nil && err.Error() != test.errMsg {
				t.Errorf("Expected error message: %s, Got: %s", test.errMsg, err.Error())
			}
		})
	}
}

func TestIsServiceItemQueryZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		args     ServiceItemQuery
		expected bool
	}{
		{
			name: "all_zero",
			args: ServiceItemQuery{
				Pov:                "",
				ChangeState:        "",
				Name:               "",
				ApplicationId:      int64(0),
				ConsumerTeamId:     int64(0),
				Limit:              int64(0),
				Offset:             int64(0),
				Ordering:           "",
				RuntimeState:       "",
				ServiceOwnerId:     int64(0),
				ServiceOwnerTeamId: int64(0),
			},
			expected: true,
		},
		{
			name: "all_zero_not_pov",
			args: ServiceItemQuery{
				Pov:                "serviceowner",
				ChangeState:        "",
				Name:               "",
				ApplicationId:      int64(0),
				ConsumerTeamId:     int64(0),
				Limit:              int64(0),
				Offset:             int64(0),
				Ordering:           "",
				RuntimeState:       "",
				ServiceOwnerId:     int64(0),
				ServiceOwnerTeamId: int64(0),
			},
			expected: true,
		},
		{
			name: "non_zero_1",
			args: ServiceItemQuery{
				Pov:                "serviceowner",
				ChangeState:        "test",
				Name:               "",
				ApplicationId:      int64(0),
				ConsumerTeamId:     int64(0),
				Limit:              int64(0),
				Offset:             int64(0),
				Ordering:           "",
				RuntimeState:       "",
				ServiceOwnerId:     int64(0),
				ServiceOwnerTeamId: int64(0),
			},
			expected: false,
		},
		{
			name: "non_zero_2",
			args: ServiceItemQuery{
				Pov:                "serviceowner",
				ChangeState:        "",
				Name:               "test",
				ApplicationId:      int64(0),
				ConsumerTeamId:     int64(0),
				Limit:              int64(0),
				Offset:             int64(0),
				Ordering:           "",
				RuntimeState:       "",
				ServiceOwnerId:     int64(0),
				ServiceOwnerTeamId: int64(0),
			},
			expected: false,
		},
		{
			name: "non_zero_3",
			args: ServiceItemQuery{
				Pov:                "serviceowner",
				ChangeState:        "",
				Name:               "",
				ApplicationId:      int64(1),
				ConsumerTeamId:     int64(0),
				Limit:              int64(0),
				Offset:             int64(0),
				Ordering:           "",
				RuntimeState:       "",
				ServiceOwnerId:     int64(0),
				ServiceOwnerTeamId: int64(0),
			},
			expected: false,
		},
		{
			name: "non_zero_4",
			args: ServiceItemQuery{
				Pov:                "serviceowner",
				ChangeState:        "",
				Name:               "",
				ApplicationId:      int64(0),
				ConsumerTeamId:     int64(1),
				Limit:              int64(0),
				Offset:             int64(0),
				Ordering:           "",
				RuntimeState:       "",
				ServiceOwnerId:     int64(0),
				ServiceOwnerTeamId: int64(0),
			},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isServiceItemQueryZeroValue(test.args)

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected: %v, Got: %v", test.expected, result)
			}

		})
	}
}

func TestGetQueryParam(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		expected string
	}{
		{
			name: "single_query_param_application_id",
			args: map[string]interface{}{
				"pov":                   "serviceowner",
				"change_state":          "",
				"name":                  "",
				"application_id":        int64(123),
				"consumer_team_id":      int64(0),
				"limit":                 int64(0),
				"offset":                int64(0),
				"ordering":              "",
				"runtime_state":         "",
				"service_owner_id":      int64(0),
				"service_owner_team_id": int64(0),
			},
			expected: "?&application_id=123",
		},
		{
			name: "single_query_param_change_state",
			args: map[string]interface{}{
				"pov":                   "serviceowner",
				"change_state":          "CHANGES_APPROVED",
				"name":                  "",
				"application_id":        int64(0),
				"consumer_team_id":      int64(0),
				"limit":                 int64(0),
				"offset":                int64(0),
				"ordering":              "",
				"runtime_state":         "",
				"service_owner_id":      int64(0),
				"service_owner_team_id": int64(0),
			},
			expected: "?&change_state=CHANGES_APPROVED",
		},
		{
			name: "single_query_param_consumer_team_id",
			args: map[string]interface{}{
				"pov":                   "serviceowner",
				"change_state":          "",
				"name":                  "",
				"application_id":        int64(0),
				"consumer_team_id":      int64(123),
				"limit":                 int64(0),
				"offset":                int64(0),
				"ordering":              "",
				"runtime_state":         "",
				"service_owner_id":      int64(0),
				"service_owner_team_id": int64(0),
			},
			expected: "?&consumer_team_id=123",
		},
		{
			name: "single_query_param_limit",
			args: map[string]interface{}{
				"pov":                   "serviceowner",
				"change_state":          "",
				"name":                  "",
				"application_id":        int64(0),
				"consumer_team_id":      int64(0),
				"limit":                 int64(123),
				"offset":                int64(0),
				"ordering":              "",
				"runtime_state":         "",
				"service_owner_id":      int64(0),
				"service_owner_team_id": int64(0),
			},
			expected: "?&limit=123",
		},
		{
			name: "single_query_param_offset",
			args: map[string]interface{}{
				"pov":                   "serviceowner",
				"change_state":          "",
				"name":                  "",
				"application_id":        int64(0),
				"consumer_team_id":      int64(0),
				"limit":                 int64(0),
				"offset":                int64(123),
				"ordering":              "",
				"runtime_state":         "",
				"service_owner_id":      int64(0),
				"service_owner_team_id": int64(0),
			},
			expected: "?&offset=123",
		},
		{
			name: "single_query_param_ordering",
			args: map[string]interface{}{
				"pov":                   "serviceowner",
				"change_state":          "",
				"name":                  "",
				"application_id":        int64(0),
				"consumer_team_id":      int64(0),
				"limit":                 int64(0),
				"offset":                int64(0),
				"ordering":              "name",
				"runtime_state":         "",
				"service_owner_id":      int64(0),
				"service_owner_team_id": int64(0),
			},
			expected: "?&ordering=name",
		},
		{
			name: "single_query_param_runtime_state",
			args: map[string]interface{}{
				"pov":                   "serviceowner",
				"change_state":          "",
				"name":                  "",
				"application_id":        int64(0),
				"consumer_team_id":      int64(0),
				"limit":                 int64(0),
				"offset":                int64(0),
				"ordering":              "",
				"runtime_state":         "RUNNING",
				"service_owner_id":      int64(0),
				"service_owner_team_id": int64(0),
			},
			expected: "?&runtime_state=RUNNING",
		},
		{
			name: "single_query_param_service_owner_id",
			args: map[string]interface{}{
				"pov":                   "serviceowner",
				"change_state":          "",
				"name":                  "",
				"application_id":        int64(0),
				"consumer_team_id":      int64(0),
				"limit":                 int64(0),
				"offset":                int64(0),
				"ordering":              "",
				"runtime_state":         "",
				"service_owner_id":      int64(123),
				"service_owner_team_id": int64(0),
			},
			expected: "?&service_owner_id=123",
		},
		{
			name: "single_query_param_service_owner_team_id",
			args: map[string]interface{}{
				"pov":                   "serviceowner",
				"change_state":          "",
				"name":                  "",
				"application_id":        int64(0),
				"consumer_team_id":      int64(0),
				"limit":                 int64(0),
				"offset":                int64(0),
				"ordering":              "",
				"runtime_state":         "",
				"service_owner_id":      int64(0),
				"service_owner_team_id": int64(123),
			},
			expected: "?&service_owner_team_id=123",
		},
		{
			name: "multi_query_param_1",
			args: map[string]interface{}{
				"pov":                   "serviceowner",
				"change_state":          "RUNNING",
				"name":                  "test-name",
				"application_id":        int64(0),
				"consumer_team_id":      int64(0),
				"limit":                 int64(0),
				"offset":                int64(0),
				"ordering":              "",
				"runtime_state":         "",
				"service_owner_id":      int64(0),
				"service_owner_team_id": int64(0),
			},
			expected: "?&change_state=RUNNING&name=test-name",
		},
		{
			name: "multi_query_param_2",
			args: map[string]interface{}{
				"pov":                   "serviceowner",
				"change_state":          "RUNNING",
				"name":                  "",
				"application_id":        int64(123),
				"consumer_team_id":      int64(0),
				"limit":                 int64(0),
				"offset":                int64(0),
				"ordering":              "",
				"runtime_state":         "",
				"service_owner_id":      int64(0),
				"service_owner_team_id": int64(0),
			},
			expected: "?&application_id=123&change_state=RUNNING",
		},
		{
			name: "multi_query_param_3",
			args: map[string]interface{}{
				"pov":                   "serviceowner",
				"change_state":          "RUNNING",
				"name":                  "",
				"application_id":        int64(123),
				"consumer_team_id":      int64(123),
				"limit":                 int64(0),
				"offset":                int64(0),
				"ordering":              "",
				"runtime_state":         "",
				"service_owner_id":      int64(0),
				"service_owner_team_id": int64(0),
			},
			expected: "?&application_id=123&change_state=RUNNING&consumer_team_id=123",
		},
		{
			name: "multi_query_param_4",
			args: map[string]interface{}{
				"pov":                   "serviceowner",
				"change_state":          "RUNNING",
				"name":                  "",
				"application_id":        int64(123),
				"consumer_team_id":      int64(123),
				"limit":                 int64(123),
				"offset":                int64(0),
				"ordering":              "",
				"runtime_state":         "",
				"service_owner_id":      int64(0),
				"service_owner_team_id": int64(0),
			},
			expected: "?&application_id=123&change_state=RUNNING&consumer_team_id=123&limit=123",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			queryObject, _ := NewServiceItemQuery(test.args)

			result := queryObject.GetQueryParam()

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected: %v, Got: %v", test.expected, result)
			}
		})
	}
}

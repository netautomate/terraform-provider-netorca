// Copyright (c) HashiCorp, Inc.

package netorca

import (
	"reflect"
	"testing"
)

func TestChangeInstances(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		expected *ChangeInstanceQuery
		errMsg   string
	}{
		{
			name: "valid_case_fully_populated",
			args: map[string]interface{}{
				"pov":                   "consumer",
				"application_id":        int64(123),
				"consumer_team_id":      int64(456),
				"limit":                 int64(10),
				"offset":                int64(20),
				"ordering":              "name",
				"service_owner_team_id": int64(101),
				"change_type":           "change_type",
				"commit_id":             "commit_id",
				"service_id":            int64(1),
				"service_item_id":       int64(2),
				"service_name":          "service_name",
				"state":                 "state",
				"submission_id":         int64(3),
			},
			expected: &ChangeInstanceQuery{
				Pov:                "consumer",
				ApplicationId:      int64(123),
				ChangeType:         "change_type",
				CommitId:           "commit_id",
				ConsumerTeamId:     int64(456),
				Limit:              int64(10),
				Offset:             int64(20),
				Ordering:           "name",
				ServiceId:          int64(1),
				ServiceItemId:      int64(2),
				ServiceName:        "service_name",
				ServiceOwnerTeamId: int64(101),
				State:              "state",
				SubmissionId:       int64(3),
			},
			errMsg: "",
		},
		{
			name: "valid_case_partially_populated",
			args: map[string]interface{}{
				"pov": "serviceowner",
			},
			expected: &ChangeInstanceQuery{
				Pov:                "serviceowner",
				ApplicationId:      int64(0),
				ConsumerTeamId:     int64(0),
				Limit:              int64(0),
				Offset:             int64(0),
				Ordering:           "",
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
			name: "invalid_type_application_id",
			args: map[string]interface{}{
				"application_id": "not_an_int",
			},
			expected: nil,
			errMsg:   "application_id not passed as a uint64",
		},
		{
			name: "invalid_type_consumer_team_id",
			args: map[string]interface{}{
				"consumer_team_id": "not_an_int",
			},
			expected: nil,
			errMsg:   "consumer_team_id not passed as a uint64",
		},
		{
			name: "invalid_type_limit",
			args: map[string]interface{}{
				"limit": "not_an_int",
			},
			expected: nil,
			errMsg:   "limit not passed as a uint64",
		},
		{
			name: "invalid_type_offset",
			args: map[string]interface{}{
				"offset": "not_an_int",
			},
			expected: nil,
			errMsg:   "offset not passed as a uint64",
		},
		{
			name: "invalid_type_ordering",
			args: map[string]interface{}{
				"ordering": 123,
			},
			expected: nil,
			errMsg:   "ordering not passed as a string",
		},
		{
			name: "invalid_type_service_owner_team_id",
			args: map[string]interface{}{
				"service_owner_team_id": "not_an_int",
			},
			expected: nil,
			errMsg:   "service_owner_team_id not passed as a uint64",
		},
		{
			name: "invalid_type_change_type",
			args: map[string]interface{}{
				"change_type": 123,
			},
			expected: nil,
			errMsg:   "change_type not passed as a string",
		},
		{
			name: "invalid_type_commit_id",
			args: map[string]interface{}{
				"commit_id": 123,
			},
			expected: nil,
			errMsg:   "commit_id not passed as a string",
		},
		{
			name: "invalid_type_service_id",
			args: map[string]interface{}{
				"service_id": "not_an_int",
			},
			expected: nil,
			errMsg:   "service_id not passed as a uint64",
		},
		{
			name: "invalid_type_service_item_id",
			args: map[string]interface{}{
				"service_item_id": "not_an_int",
			},
			expected: nil,
			errMsg:   "service_item_id not passed as a uint64",
		},
		{
			name: "invalid_type_service_name",
			args: map[string]interface{}{
				"service_name": 123,
			},
			expected: nil,
			errMsg:   "service_name not passed as a string",
		},
		{
			name: "invalid_type_state",
			args: map[string]interface{}{
				"state": 123,
			},
			expected: nil,
			errMsg:   "state not passed as a string",
		},
		{
			name: "invalid_type_submission_id",
			args: map[string]interface{}{
				"submission_id": "not_an_int",
			},
			expected: nil,
			errMsg:   "submission_id not passed as a uint64",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := NewChangeInstanceQuery(test.args)

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected: %v, Got: %v", test.expected, result)
			}

			if err != nil && err.Error() != test.errMsg {
				t.Errorf("Expected error message: %s, Got: %s", test.errMsg, err.Error())
			}
		})
	}
}

func TestGetQueryParams(t *testing.T) {
	tests := []struct {
		name     string
		query    ChangeInstanceQuery
		expected string
	}{
		{
			name: "all_fields_populated",
			query: ChangeInstanceQuery{
				ServiceId:          1,
				ServiceItemId:      2,
				ServiceName:        "test_service",
				ServiceOwnerTeamId: 3,
				State:              "active",
				SubmissionId:       4,
			},
			expected: "?service_id=1&service_item_id=2&service_name=test_service&service_owner_team_id=3&state=active&submission_id=4",
		},
		{
			name: "some_fields_populated",
			query: ChangeInstanceQuery{
				ServiceId:    1,
				ServiceName:  "test_service",
				State:        "active",
				SubmissionId: 4,
			},
			expected: "?service_id=1&service_name=test_service&state=active&submission_id=4",
		},
		{
			name: "some_fields_populated_without_service_id",
			query: ChangeInstanceQuery{
				ServiceName: "test_service",
				State:       "active",
			},
			expected: "?service_name=test_service&state=active",
		},
		{
			name:     "no_fields_populated",
			query:    ChangeInstanceQuery{},
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.query.GetQueryParam()
			if result != test.expected {
				t.Errorf("Expected: %v, Got: %v", test.expected, result)
			}
		})
	}
}

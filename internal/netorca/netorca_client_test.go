// Copyright (c) HashiCorp, Inc.

package netorca

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestNewNetorcaClient(t *testing.T) {
	t.Run("TestNewNetorcaClient returns correct initialized client", func(t *testing.T) {
		api_key := "123456"
		base_url := "https://example.com"

		c := NewClient(&base_url, &api_key, context.Background())

		if c.GetApiKey() != "Api-Key 123456" {
			t.Errorf("Expected %s, got %s", "123456", c.apiKey)
		}

	})
}

func TestChangeInstanceGetById(t *testing.T) {
	// Mock response from the server with actual netorca response
	mockResponse, err := os.ReadFile("testdata/change_instance_200.json")
	if err != nil {
		t.Fatalf("Failed to read mock response file: %v", err)
	}

	if mockResponse == nil {
		t.Fatalf("Failed to read mock response file: %v", err)
	}
	// mock response from server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponse)
	}))
	defer server.Close()

	apikey := "123456"
	client := NewClient(&server.URL, &apikey, context.Background())

	expected := ChangeInstance{}
	if err := json.Unmarshal(mockResponse, &expected); err != nil {
		t.Fatalf("Failed to unmarshal mock response: %v", err)
	}
	if reflect.DeepEqual(expected, ChangeInstance{}) {
		t.Fatalf("Failed to unmarshal mock response: %v", err)
	}

	id := int64(123)
	pov := "consumer"

	result, err := client.ChangeInstanceGetById(id, pov)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v, Got: %v", expected, result)
	}
}

func TestChangeInstanceGetByIdNotFound(t *testing.T) {
	// {"detail":"Not found."} is the response
	mockResponse := []byte(`{"detail":"Not found."}`)
	// mock response from server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write(mockResponse)

	}))
	defer server.Close()

	apikey := "123456"
	client := NewClient(&server.URL, &apikey, context.Background())

	id := int64(123)
	pov := "consumer"
	_, err := client.ChangeInstanceGetById(id, pov)

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	expected_err := fmt.Errorf("http code: %d\n response: %s\nurl: %s\nmethod: GET", 404, []byte(`{"detail":"Not found."}`), fmt.Sprintf("%s/v1/orcabase/consumer/change_instances/123/", server.URL))

	if err.Error() != expected_err.Error() {
		t.Fatalf("Expected error message: Not found., got %v", err.Error())
	}
}

func TestServiceItemGetList(t *testing.T) {
	mockResponse, err := os.ReadFile("testdata/service_items_200.json")
	if err != nil {
		t.Fatalf("Failed to read mock response file: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponse)
	}))
	defer server.Close()

	apikey := "123456"
	client := NewClient(&server.URL, &apikey, context.Background())

	expected := NetOrcaServiceItem{}
	if err := json.Unmarshal(mockResponse, &expected); err != nil {
		t.Fatalf("Failed to unmarshal mock response: %v", err)
	}
	if reflect.DeepEqual(expected, ServiceItem{}) {
		t.Fatalf("Failed to unmarshal mock response: %v", err)
	}
	q, err := NewServiceItemQuery(
		map[string]interface{}{
			"service_item_id": 123,
		},
	)

	if err != nil {
		t.Fatalf("Failed to create query: %v", err)
	}

	result, err := client.ServiceItemsGet(q)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected: %v, Got: %v", expected, result)
	}
}

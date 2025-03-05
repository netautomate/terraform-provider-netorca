// Copyright (c) HashiCorp, Inc.

package netorca

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type NetOrcaClient struct {
	baseUrl string
	client  *http.Client
	apiKey  string
}

func formatApiKey(apikey string) string {
	return fmt.Sprintf("Api-Key %s", apikey)
}

func getClientTransportDefaults() *http.Transport {
	return &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}

}

func (n NetOrcaClient) GetApiKey() string {
	return n.apiKey
}

func NewClient(url, apikey *string, ctx context.Context) *NetOrcaClient {
	tflog.Info(ctx, "Building NetOrca client")
	tr := getClientTransportDefaults()
	httpClient := http.Client{Transport: tr}

	return &NetOrcaClient{
		client:  &httpClient,
		baseUrl: *url,
		apiKey:  formatApiKey(*apikey),
	}
}

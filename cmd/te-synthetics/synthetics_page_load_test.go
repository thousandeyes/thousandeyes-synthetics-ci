// Copyright 2026 Cisco Systems, Inc. and its affiliates
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/thousandeyes/thousandeyes-synthetics-ci/internal/synthetics"
)

func TestSyntheticsPageLoad(t *testing.T) {
	const target = "https://target.invalid/"
	for _, tc := range []struct {
		name, thresholds, wantStatus string
		loadTime, domTime, numErrors int
		wantPolls                    int32
	}{
		{"waits for both agents and allows errors within budget", `{"maxPageLoadTime":250,"maxDomLoadTime":125,"maxErrors":1}`, "passed", 250, 125, 1, 3},
		{"fails the page load threshold", `{"maxPageLoadTime":249}`, "failed", 250, 125, 0, 1},
		{"fails the DOM load threshold", `{"maxDomLoadTime":124}`, "failed", 250, 125, 0, 1},
		{"fails on a resource error by default", `{}`, "failed", 250, 125, 1, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var creates, polls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != "Bearer synthetics-test-token" || r.URL.Query().Get("aid") != "123" {
					t.Error("request is missing authentication or account-group context")
				}
				switch {
				case r.Method == http.MethodPost && r.URL.Path == "/v7/tests/page-load/instant":
					creates.Add(1)
					if r.URL.Query().Get("expand") != "agent" {
						t.Error("creation must request expanded agent assignments")
					}
					var payload struct {
						URL      string `json:"url"`
						TestName string `json:"testName"`
						Agents   []struct {
							ID string `json:"agentId"`
						} `json:"agents"`
						BrowserLanguage string `json:"browserLanguage"`
					}
					if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
						t.Errorf("decode creation payload: %v", err)
					}
					if payload.URL != target || !strings.HasPrefix(payload.TestName, "ci-") || payload.BrowserLanguage != "en-US" ||
						len(payload.Agents) != 2 || payload.Agents[0].ID != "6" || payload.Agents[1].ID != "7" {
						t.Errorf("unexpected creation payload: %+v", payload)
					}
					syntheticJSON(t, w, http.StatusCreated, map[string]any{
						"testId": "1234", "type": "page-load",
						"agents": []any{
							map[string]any{"agentId": "6", "agentName": "First agent", "agentType": "cloud"},
							map[string]any{"agentId": "7", "agentName": "Second agent", "agentType": "enterprise"},
						},
						"_links": map[string]any{"testResults": []any{
							map[string]string{"href": "/v7/test-results/1234/http-server"},
							map[string]string{"href": "/v7/test-results/1234/page-load"},
						}},
					})
				case r.Method == http.MethodGet && r.URL.Path == "/v7/test-results/1234/page-load":
					poll := polls.Add(1)
					rows := make([]map[string]any, 0, 2)
					for _, id := range []string{"6", "7"} {
						row := map[string]any{
							"agent": map[string]any{"agentId": id}, "roundId": 1700000000,
							"responseTime": 34.35, "domLoadTime": tc.domTime, "pageLoadTime": tc.loadTime,
							"numErrors": tc.numErrors, "numObjects": 2, "totalSize": 16384,
						}
						if tc.wantPolls > 1 {
							// Each poll needs both agents to report the required fields.
							switch {
							case poll == 1 && id == "6":
								delete(row, "pageLoadTime")
							case poll == 2 && id == "7":
								delete(row, "numErrors")
							}
						}
						rows = append(rows, row)
					}
					syntheticJSON(t, w, http.StatusOK, map[string]any{"test": map[string]string{"testId": "1234"}, "results": rows})
				default:
					t.Errorf("unexpected API call: %s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
				}
			}))
			defer server.Close()

			suite := `{"version":1,"accountGroupId":"123","timeout":"90s","pollInterval":"1ms","tests":[{
				"name":"browser-check","type":"page-load",
				"config":{"url":"` + target + `","browserLanguage":"en-US","agents":[{"agentId":"6"},{"agentId":"7"}]},
				"thresholds":` + tc.thresholds + `
			}]}`
			junitPath := filepath.Join(t.TempDir(), "page-load.xml")
			stdout, stderr, err := executeSyntheticSuite(t, server.URL, suite, "--json", "--junit-report", junitPath)
			if (err == nil) != (tc.wantStatus == "passed") {
				t.Fatalf("expected status %s: %v\nstdout: %s\nstderr: %s", tc.wantStatus, err, stdout, stderr)
			}
			if creates.Load() != 1 || polls.Load() != tc.wantPolls {
				t.Fatalf("expected one creation and %d Page Load reads, got %d creations and %d reads", tc.wantPolls, creates.Load(), polls.Load())
			}
			var report synthetics.Report
			if err := json.Unmarshal([]byte(stdout), &report); err != nil {
				t.Fatalf("decode JSON report: %v\n%s", err, stdout)
			}
			if len(report.Tests) != 1 {
				t.Fatalf("expected one test result: %s", stdout)
			}
			test := report.Tests[0]
			if test.Status != tc.wantStatus || test.Type != "page-load" || test.Target != target || test.TestID != "1234" || len(test.Agents) != 2 {
				t.Fatalf("unexpected Page Load test result: %s", stdout)
			}
			wantMetrics := map[string]float64{
				"responseTime": 34.35, "domLoadTime": float64(tc.domTime), "pageLoadTime": float64(tc.loadTime),
				"numErrors": float64(tc.numErrors), "numObjects": 2, "totalSize": 16384,
			}
			for index, agent := range test.Agents {
				wantName := []string{"First agent", "Second agent"}[index]
				if agent.Status != tc.wantStatus || agent.AgentName != wantName || !reflect.DeepEqual(agent.Metrics, wantMetrics) {
					t.Errorf("unexpected Page Load measurements: %+v; expected %v", agent, wantMetrics)
				}
			}
			junit, err := os.ReadFile(junitPath)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(junit), `classname="page-load"`) || !strings.Contains(string(junit), "pageLoadTime") ||
				strings.Contains(string(junit), "<failure") != (tc.wantStatus == "failed") {
				t.Errorf("unexpected Page Load JUnit report: %s", junit)
			}
		})
	}
}

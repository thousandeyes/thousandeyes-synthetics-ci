// Copyright 2026 Cisco Systems, Inc. and its affiliates
// SPDX-License-Identifier: Apache-2.0

package synthetics

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"
)

func TestRunRedactsTestNamesAcrossReports(t *testing.T) {
	names := []string{"check-fake-backend-token-and-fake-backend-token", "  東京 health  ", "fake-backend-token"}
	wantNames := []string{"check-[REDACTED]-and-[REDACTED]", "  東京 health  ", "[REDACTED]"}
	for _, scenario := range []struct {
		name, mode, rule string
		wantError        bool
		statuses         []string
	}{
		{"completed", "measurements", "blocking", false, []string{"passed", "passed", "passed"}},
		{"explicitly-skipped", "malformed", "skipped", false, []string{"skipped", "skipped", "skipped"}},
		{"error-and-unstarted", "malformed", "blocking", true, []string{"error", "skipped", "skipped"}},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			client := helperBackendClient(t, scenario.mode)
			t.Setenv("TE_SYNTHETICS_HELPER_MEASUREMENTS", `"responseCode":200`)
			config := Config{Version: 1, AccountGroupID: "123", PollInterval: "1ms"}
			for _, name := range names {
				config.Tests = append(config.Tests, TestConfig{
					Name: name, Type: "http-server", ExecutionRule: scenario.rule,
					Config: map[string]json.RawMessage{
						"url":    json.RawMessage(`"https://target.invalid"`),
						"agents": json.RawMessage(`[{"agentId":"6"}]`),
					},
				})
			}
			report, err := Run(context.Background(), config, Options{BaseURL: client.base.String(), Token: client.token, APICLIPath: client.apiCLIPath})
			if (err != nil) != scenario.wantError || report.Failed() != scenario.wantError {
				t.Fatalf("execution outcome changed: failed=%v error=%v", report.Failed(), err)
			}
			var jsonOutput, junitOutput, terminal bytes.Buffer
			if err := WriteJSON(&jsonOutput, report); err != nil {
				t.Fatal(err)
			}
			if err := WriteJUnit(&junitOutput, report); err != nil {
				t.Fatal(err)
			}
			if err := WriteSummary(&terminal, report); err != nil {
				t.Fatal(err)
			}
			var decoded Report
			if err := json.Unmarshal(jsonOutput.Bytes(), &decoded); err != nil {
				t.Fatal(err)
			}
			var junit struct {
				Cases []struct {
					Name string `xml:"name,attr"`
				} `xml:"testcase"`
			}
			if err := xml.Unmarshal(junitOutput.Bytes(), &junit); err != nil {
				t.Fatal(err)
			}
			if len(report.Tests) != len(names) || len(decoded.Tests) != len(names) || len(junit.Cases) != len(names) {
				t.Fatal("report lost configured tests")
			}
			for index, want := range wantNames {
				if report.Tests[index].Name != want || decoded.Tests[index].Name != want || !strings.HasPrefix(junit.Cases[index].Name, want) || !strings.Contains(terminal.String(), want) {
					t.Errorf("test %d did not retain the expected safe name %q in every report", index, want)
				}
				if report.Tests[index].Status != scenario.statuses[index] || decoded.Tests[index].Status != scenario.statuses[index] {
					t.Errorf("test %d changed status", index)
				}
				if config.Tests[index].Name != names[index] {
					t.Errorf("test %d changed the input configuration", index)
				}
			}
			for format, output := range map[string]string{"JSON": jsonOutput.String(), "JUnit": junitOutput.String(), "terminal": terminal.String()} {
				if strings.Contains(output, client.token) {
					t.Errorf("%s exposed the API token in a configured test name", format)
				}
			}
		})
	}
}

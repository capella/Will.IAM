// +build integration

package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	helpers "github.com/ghostec/Will.IAM/testing"
)

func beforeEachAMHandlers(t *testing.T) {
	t.Helper()
	storage := helpers.GetStorage(t)
	rels := []string{"permissions", "role_bindings", "service_accounts", "roles"}
	for _, rel := range rels {
		if _, err := storage.PG.DB.Exec(
			fmt.Sprintf("DELETE FROM %s", rel),
		); err != nil {
			panic(err)
		}
	}
}

func TestAMListHandler(t *testing.T) {
	beforeEachAMHandlers(t)
	rootSA := helpers.CreateRootServiceAccount(t)
	// saUC := helpers.GetServiceAccountsUseCase(t)
	app := helpers.GetApp(t)
	type testCase struct {
		reqPath        string
		expectedOutput []string
	}
	testCases := []testCase{
		testCase{
			reqPath:        "/am",
			expectedOutput: []string{"Will.IAM"},
		},
		testCase{
			reqPath:        "/am?prefix=Will.I",
			expectedOutput: []string{"Will.IAM"},
		},
		testCase{
			reqPath:        "/am?prefix=wi",
			expectedOutput: []string{},
		},
		testCase{
			reqPath:        "/am?prefix=x",
			expectedOutput: []string{},
		},
		testCase{
			reqPath: "/am?prefix=Will.IAM::",
			expectedOutput: []string{
				"Will.IAM::CreateRoles", "Will.IAM::EditRole", "Will.IAM::ListRoles",
				"Will.IAM::CreateServiceAccounts", "Will.IAM::EditServiceAccount",
				"Will.IAM::ListServiceAccounts",
			},
		},
		testCase{
			reqPath: "/am?prefix=Will.IAM::Edit",
			expectedOutput: []string{
				"Will.IAM::EditRole", "Will.IAM::EditServiceAccount",
			},
		},
		testCase{
			reqPath:        "/am?prefix=Will.IAM::EditR",
			expectedOutput: []string{"Will.IAM::EditRole"},
		},
	}
	for _, tt := range testCases {
		req, _ := http.NewRequest("GET", tt.reqPath, nil)
		req.Header.Set("Authorization", fmt.Sprintf(
			"KeyPair %s:%s", rootSA.KeyID, rootSA.KeySecret,
		))
		rec := helpers.DoRequest(t, req, app.GetRouter())
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status 200. Got %d", rec.Code)
		}
		suggs := []map[string]interface{}{}
		err := json.Unmarshal([]byte(rec.Body.String()), &suggs)
		if err != nil {
			t.Errorf("Unexpected error %s", err.Error())
			return
		}
		if len(suggs) != len(tt.expectedOutput) {
			t.Errorf(
				"Expected suggestions to have length %d. Got %d",
				len(tt.expectedOutput), len(suggs),
			)
			return
		}
		for i := range suggs {
			if suggs[i]["prefix"].(string) != tt.expectedOutput[i] {
				t.Errorf(
					"Expected suggs[%d] to be %s. Got %s",
					i, suggs[i]["prefix"].(string), tt.expectedOutput[i],
				)
				return
			}
		}
	}
}

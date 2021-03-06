// Copyright 2017 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const (
	tmpDirPrefix = "test-infra_package-coverage-"
)

var (
	tmpDir string
)

func TestParseReport(t *testing.T) {
	exampleReport := "?   \tpilot/cmd\t[no test files]\nok  \tpilot/model\t1.3s\tcoverage: 90.2% of statements"
	reportFile := filepath.Join(tmpDir, "report")
	if err := ioutil.WriteFile(reportFile, []byte(exampleReport), 0644); err != nil {
		t.Errorf("Failed to write example report file, %v", err)
	}

	c := &codecovChecker{
		codeCoverage: make(map[string]float64),
		report:       reportFile,
	}

	if err := c.parseReport(); err != nil {
		t.Errorf("Failed to parse report, %v", err)
	} else {
		if len(c.codeCoverage) != 1 && c.codeCoverage["pilot/model"] != 90.2 {
			t.Error("Wrong result from parseReport()")
		}
	}
}

func TestSatisfiedRequirement(t *testing.T) {
	exampleRequirement := "pilot/model\t90"
	requirementFile := filepath.Join(tmpDir, "requirement2")
	if err := ioutil.WriteFile(requirementFile, []byte(exampleRequirement), 0644); err != nil {
		t.Errorf("Failed to write example requirement file, %v", err)
	}

	c := &codecovChecker{
		codeCoverage: map[string]float64{
			"pilot/model": 90.2,
		},
		requirement: requirementFile,
	}

	if err := c.checkRequirement(); err != nil {
		t.Errorf("Failed to check requirement, %v", err)
	} else {
		if len(c.failedPackage) != 0 {
			t.Error("Wrong result from checkRequirement()")
		}
	}
}

func TestMissRequirement(t *testing.T) {
	exampleRequirement := "pilot/model\t92.3"
	requirementFile := filepath.Join(tmpDir, "requirement3")
	if err := ioutil.WriteFile(requirementFile, []byte(exampleRequirement), 0644); err != nil {
		if err := ioutil.WriteFile(requirementFile, []byte(exampleRequirement), 0644); err != nil {
			t.Errorf("Failed to write example requirement file, %v", err)
		}
	}

	c := &codecovChecker{
		codeCoverage: map[string]float64{
			"pilot/model": 90.2,
		},
		requirement: requirementFile,
	}

	if err := c.checkRequirement(); err != nil {
		t.Errorf("Failed to check requirement, %v", err)
	} else {
		if len(c.failedPackage) != 1 {
			t.Error("Wrong result from checkRequirement()")
		}
	}
}

func TestPassCheck(t *testing.T) {
	exampleReport := "?   \tpilot/cmd\t[no test files]\nok  \tpilot/model\t1.3s\tcoverage: 90.2% of statements"
	reportFile := filepath.Join(tmpDir, "report4")
	if err := ioutil.WriteFile(reportFile, []byte(exampleReport), 0644); err != nil {
		t.Errorf("Failed to write example report file, %v", err)
	}

	exampleRequirement := "pilot/model\t89"
	requirementFile := filepath.Join(tmpDir, "requirement4")
	if err := ioutil.WriteFile(requirementFile, []byte(exampleRequirement), 0644); err != nil {
		if err := ioutil.WriteFile(requirementFile, []byte(exampleRequirement), 0644); err != nil {
			t.Errorf("Failed to write example requirement file, %v", err)
		}
	}

	c := &codecovChecker{
		codeCoverage: make(map[string]float64),
		report:       reportFile,
		requirement:  requirementFile,
	}

	// No other error code, code only show gcs upload failed which is expected
	if code := c.checkPackageCoverage(); code != 3 {
		t.Errorf("Unexpected return code, expected: %d, actual: %d", 3, code)
	}
}

func TestFailedCheck(t *testing.T) {
	exampleReport := "?   \tpilot/cmd\t[no test files]\nok  \tpilot/model\t1.3s\tcoverage: 90.2% of statements"
	reportFile := filepath.Join(tmpDir, "report5")
	if err := ioutil.WriteFile(reportFile, []byte(exampleReport), 0644); err != nil {
		t.Errorf("Failed to write example report file, %v", err)
	}

	exampleRequirement := "pilot/model\t93"
	requirementFile := filepath.Join(tmpDir, "requirement5")
	if err := ioutil.WriteFile(requirementFile, []byte(exampleRequirement), 0644); err != nil {
		if err := ioutil.WriteFile(requirementFile, []byte(exampleRequirement), 0644); err != nil {
			t.Errorf("Failed to write example requirement file, %v", err)
		}
	}

	c := &codecovChecker{
		codeCoverage: make(map[string]float64),
		report:       reportFile,
		requirement:  requirementFile,
	}

	if code := c.checkPackageCoverage(); code != 2 {
		t.Errorf("Unexpected return code, expected: %d, actual: %d", 2, code)
	}
}

func TestMain(m *testing.M) {
	var err error
	if tmpDir, err = ioutil.TempDir("", tmpDirPrefix); err != nil {
		log.Printf("Failed to create tmp directory: %s, %s", tmpDir, err)
		os.Exit(4)
	}

	exitCode := m.Run()

	if err := os.RemoveAll(tmpDir); err != nil {
		log.Printf("Failed to remove tmpDir %s", tmpDir)
	}

	os.Exit(exitCode)
}

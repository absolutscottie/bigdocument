package config

import (
	"testing"
)

type ConfigTestCase struct {
	name          string
	path          string
	expectedError string
	expectedName  string
	expectedAddr  string
	expectedPort  string
}

func TestLoadConfig(t *testing.T) {
	testCases := []ConfigTestCase{
		ConfigTestCase{
			name:          "file not found",
			path:          "doesnotexist.json",
			expectedError: "open doesnotexist.json: no such file or directory",
		},
		ConfigTestCase{
			name:          "invalid input",
			path:          "../../test/config/invalid-config.json",
			expectedError: "unexpected end of JSON input",
		},
		ConfigTestCase{
			name:          "valid input",
			path:          "../../test/config/valid-config.json",
			expectedError: "",
			expectedName:  "Test Configuration",
			expectedAddr:  "",
			expectedPort:  "8181",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			cfg, err := Load(testCase.path)
			if err != nil {
				if err.Error() != testCase.expectedError {
					t.Errorf("errors do not match: expected %v, got %v", err, testCase.expectedError)
				}
				return
			}

			if cfg.BindAddr != testCase.expectedAddr {
				t.Errorf("addr does not match")
			}

			if cfg.Port != testCase.expectedPort {
				t.Errorf("port does not matach")
			}

			if cfg.Name != testCase.expectedName {
				t.Errorf("name does not match")
			}
		})
	}
}

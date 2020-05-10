/*
 * Copyright © 2019 – 2020 Red Hat Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package podman

import (
	"bytes"
	"encoding/json"

	"github.com/HarryMichal/go-version"
	"github.com/containers/toolbox/pkg/shell"
	"github.com/sirupsen/logrus"
)

var (
	LogLevel = logrus.ErrorLevel
)

// CheckVersion compares provided version with the version of Podman.
//
// Takes in one string parameter that should be in the format that is used for versioning (eg. 1.0.0, 2.5.1-dev).
//
// Returns true if the Podman version is equal to or higher than the required version.
func CheckVersion(requiredVersion string) bool {
	podmanVersion, _ := GetVersion()

	podmanVersion = version.Normalize(podmanVersion)
	requiredVersion = version.Normalize(requiredVersion)

	return version.CompareSimple(podmanVersion, requiredVersion) >= 0
}

// GetContainers is a wrapper function around `podman ps --format json` command.
//
// Parameter args accepts an array of strings to be passed to the wrapped command (eg. ["-a", "--filter", "123"]).
//
// Returned value is a slice of dynamically unmarshalled json, so it needs to be treated properly.
//
// If a problem happens during execution, first argument is nil and second argument holds the error message.
func GetContainers(args ...string) ([]map[string]interface{}, error) {
	var stdout bytes.Buffer

	logLevelString := LogLevel.String()
	args = append([]string{"--log-level", logLevelString, "ps", "--format", "json"}, args...)

	if err := shell.Run("podman", nil, &stdout, nil, args...); err != nil {
		return nil, err
	}

	output := stdout.Bytes()
	var containers []map[string]interface{}

	if err := json.Unmarshal(output, &containers); err != nil {
		return nil, err
	}

	return containers, nil
}

// GetImages is a wrapper function around `podman images --format json` command.
//
// Parameter args accepts an array of strings to be passed to the wrapped command (eg. ["-a", "--filter", "123"]).
//
// Returned value is a slice of dynamically unmarshalled json, so it needs to be treated properly.
//
// If a problem happens during execution, first argument is nil and second argument holds the error message.
func GetImages(args ...string) ([]map[string]interface{}, error) {
	var stdout bytes.Buffer

	logLevelString := LogLevel.String()
	args = append([]string{"--log-level", logLevelString, "images", "--format", "json"}, args...)
	if err := shell.Run("podman", nil, &stdout, nil, args...); err != nil {
		return nil, err
	}

	output := stdout.Bytes()
	var images []map[string]interface{}

	if err := json.Unmarshal(output, &images); err != nil {
		return nil, err
	}

	return images, nil
}

// GetVersion returns version of Podman in a string
func GetVersion() (string, error) {
	var stdout bytes.Buffer

	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "version", "--format", "json"}

	if err := shell.Run("podman", nil, &stdout, nil, args...); err != nil {
		return "", err
	}

	output := stdout.Bytes()
	var jsonoutput map[string]interface{}
	if err := json.Unmarshal(output, &jsonoutput); err != nil {
		return "", err
	}

	var podmanVersion string
	podmanClientInfoInterface := jsonoutput["Client"]
	switch podmanClientInfo := podmanClientInfoInterface.(type) {
	case nil:
		podmanVersion = jsonoutput["Version"].(string)
	case map[string]interface{}:
		podmanVersion = podmanClientInfo["Version"].(string)
	}
	return podmanVersion, nil
}

// Inspect is a wrapper around 'podman inspect' command
//
// Parameter 'typearg' takes in values 'container' or 'image' that is passed to the --type flag
func Inspect(typearg string, target string) (map[string]interface{}, error) {
	var stdout bytes.Buffer

	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "inspect", "--format", "json", "--type", typearg, target}

	if err := shell.Run("podman", nil, &stdout, nil, args...); err != nil {
		return nil, err
	}

	output := stdout.Bytes()
	var info []map[string]interface{}

	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	return info[0], nil
}

func SetLogLevel(logLevel logrus.Level) {
	LogLevel = logLevel
}

func SystemMigrate(ociRuntimeRequired string) error {
	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "system", "migrate"}
	if ociRuntimeRequired != "" {
		args = append(args, []string{"--new-runtime", ociRuntimeRequired}...)
	}

	if err := shell.Run("podman", nil, nil, nil, args...); err != nil {
		return err
	}

	return nil
}

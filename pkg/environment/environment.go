/*

environment
  functions for getting data on environment conditions

*/

package environment

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sharingio/environment/pkg/common"
)

// GetEnvForPid returns a map containing environment variables and their values, given a process ID
func GetEnvForPid(pid int) (env map[string]string, err error) {
	env = map[string]string{}
	envFileRaw, err := common.ReadFile(fmt.Sprintf("/proc/%v/environ", pid))
	if err != nil {
		return map[string]string{}, err
	}
	envFile := strings.Replace(envFileRaw, "\000", "\n", -1)
	env, err = godotenv.Unmarshal(envFile)
	if err != nil {
		return map[string]string{}, err
	}
	return env, nil
}

// GetNamespace returns the name of the namespace that the process is running in
func GetNamespace() string {
	namespace, _ := common.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	return namespace
}

// GetPodLabels returns a map of the labels of the current Pod
// requires volume+volumeMount of labels using downwardAPI
func GetPodLabels() (labels map[string]string, err error) {
	labelsString, _ := common.ReadFile(common.GetPodLabelsFilePath())
	labels, err = godotenv.Unmarshal(labelsString)
	if err != nil {
		return map[string]string{}, err
	}
	return labels, nil
}

package environment

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sharingio/environment/pkg/common"
	"github.com/sharingio/environment/pkg/types"
)

func GetEnvForPid(pid int) (env map[string]string, err error) {
	env = map[string]string{}
	envFileRaw, err := common.ReadFile(fmt.Sprintf("/proc/%v/environ", pid))
	if err != nil {
		return map[string]string{}, err
	}
	envFile := strings.Replace(envFileRaw, "\000", "\n", -1)
	envParsed, err := godotenv.Unmarshal(envFile)
	env[string(types.EnvVarNameSharingioPairExposerDisabled)] =
		envParsed[string(types.EnvVarNameSharingioPairExposerDisabled)]
	env[string(types.EnvVarNameSharingioPairExposerHostname)] =
		envParsed[string(types.EnvVarNameSharingioPairExposerHostname)]
	env[string(types.EnvVarNameSharingioPairExposerAllowedPorts)] =
		envParsed[string(types.EnvVarNameSharingioPairExposerAllowedPorts)]
	env[string(types.EnvVarNameSharingioPairExposerDisabledPorts)] =
		envParsed[string(types.EnvVarNameSharingioPairExposerDisabledPorts)]
	// deprecated
	env[string(types.EnvVarNameSharingioPairSetHostname)] =
		envParsed[string(types.EnvVarNameSharingioPairSetHostname)]
	env[string(types.EnvVarNameSharingioPairExposerDisabled)] =
		envParsed[string(types.EnvVarNameSharingioPairExposerDisabled)]
	env[string(types.EnvVarNameSharingioPairIngressReconcilerAllowedPorts)] =
		envParsed[string(types.EnvVarNameSharingioPairIngressReconcilerAllowedPorts)]
	env[string(types.EnvVarNameSharingioPairIngressReconcilerDisabledPorts)] =
		envParsed[string(types.EnvVarNameSharingioPairIngressReconcilerDisabledPorts)]
	return env, nil
}

func GetNamespace() string {
	namespace, _ := common.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	return namespace
}

func GetPodLabels() (labels map[string]string, err error) {
	labelsString, _ := common.ReadFile(common.GetPodLabelsFilePath())
	labels, err = godotenv.Unmarshal(labelsString)
	if err != nil {
		return map[string]string{}, err
	}
	return labels, nil
}

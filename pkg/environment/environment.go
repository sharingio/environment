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
	env[string(types.EnvironmentVariableNameSharingioPairSetHostname)] = envParsed[string(types.EnvironmentVariableNameSharingioPairSetHostname)]
	env[string(types.EnvironmentVariableNameSharingioPairIngressReconcilerAllowedPorts)] = envParsed[string(types.EnvironmentVariableNameSharingioPairIngressReconcilerAllowedPorts)]
	env[string(types.EnvironmentVariableNameSharingioPairIngressReconcilerDisabledPorts)] = envParsed[string(types.EnvironmentVariableNameSharingioPairIngressReconcilerDisabledPorts)]
	return env, nil
}

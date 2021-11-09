package types

import (
	"github.com/cakturk/go-netstat/netstat"
)

type Process struct {
	Name          string            `json:"name"`
	Pid           int               `json:"pid"`
	Uid           uint32            `json:"uid"`
	LocalAddr     *netstat.SockAddr `json:"localAddr"`
	Env           map[string]string `json:"env"`
	Hostname      string            `json:"hostname"`
	AllowedPorts  []int             `json:"allowedPorts"`
	DisabledPorts []int             `json:"disabledPorts"`
}

type EnvironmentVariableName string

const (
	EnvironmentVariableNameSharingioPairSetHostname                    EnvironmentVariableName = "SHARINGIO_PAIR_SET_HOSTNAME"
	EnvironmentVariableNameSharingioPairIngressReconcilerAllowedPorts  EnvironmentVariableName = "SHARINGIO_PAIR_INGRESS_RECONCILER_ALLOWED_PORTS"
	EnvironmentVariableNameSharingioPairIngressReconcilerDisabledPorts EnvironmentVariableName = "SHARINGIO_PAIR_INGRESS_RECONCILER_DISABLED_PORTS"
)

package types

import (
	"net"
)

const DotSharingDotIoExposerTemplatesFolderPath string = "/home/ii/.sharing.io/cluster-api/exposer/templates"

type Process struct {
	Name          string `json:"name"`
	Protocol      string `json:"protocol"`
	Pid           int    `json:"pid"`
	Uid           uint32 `json:"uid"`
	IP            net.IP `json:"ip"`
	Port          uint16 `json:"port"`
	Hostname      string `json:"hostname"`
	AllowedPorts  []int  `json:"allowedPorts"`
	DisabledPorts []int  `json:"disabledPorts"`
	Disabled      bool   `json:"disabled"`

	PodName      string            `json:"podName"`
	PodNamespace string            `json:"podNamespace"`
	PodLabels    map[string]string `json:"podLabels"`
	ServiceName  string            `json:"-"`
	ServicePort  uint16            `json:"-"`
	ExternalIP   string            `json:"-"`
	IngressHost  string            `json:"-"`
}

type EnvironmentVariableName string

// TODO shorten to EnvVarName ...
// TODO rename to SHARINGIO_PAIR_EXPOSER ...
const (
	EnvironmentVariableNameSharingioPairSetHostname                    EnvironmentVariableName = "SHARINGIO_PAIR_SET_HOSTNAME"
	EnvironmentVariableNameSharingioPairExposerDisabled                EnvironmentVariableName = "SHARINGIO_PAIR_EXPOSER_DISABLED"
	EnvironmentVariableNameSharingioPairIngressReconcilerAllowedPorts  EnvironmentVariableName = "SHARINGIO_PAIR_INGRESS_RECONCILER_ALLOWED_PORTS"
	EnvironmentVariableNameSharingioPairIngressReconcilerDisabledPorts EnvironmentVariableName = "SHARINGIO_PAIR_INGRESS_RECONCILER_DISABLED_PORTS"
)

type TemplateYamlFile string

const (
	TemplateYamlFileService              TemplateYamlFile = "service.yaml"
	TemplateYamlFileIngress              TemplateYamlFile = "ingress.yaml"
	TemplateYamlFileIngressV118OrEarlier TemplateYamlFile = "ingress-v1.18-or-earlier.yaml"
)

const IPAllInterfaces = "::"

/*

types
  package for types used through the targets

*/

package types

import (
	"net"
)

// Protocol is network protocols
type Protocol string

const (
	ProtocolTCP Protocol = "TCP"
	ProtocolUDP Protocol = "UDP"
)

// Process is a process that's listening on a port
type Process struct {
	Name          string   `json:"name"`
	Protocol      Protocol `json:"protocol"`
	Pid           int      `json:"pid"`
	Uid           uint32   `json:"uid"`
	IP            net.IP   `json:"ip"`
	Port          uint16   `json:"port"`
	Hostname      string   `json:"hostname"`
	AllowedPorts  []int    `json:"allowedPorts"`
	DisabledPorts []int    `json:"disabledPorts"`
	Disabled      bool     `json:"disabled"`

	ResourceName     string            `json:"resourceName"`
	PodName          string            `json:"podName"`
	PodNamespace     string            `json:"podNamespace"`
	PodLabels        map[string]string `json:"podLabels"`
	ServiceName      string            `json:"-"`
	ServicePort      uint16            `json:"-"`
	ExternalIP       string            `json:"-"`
	IngressHost      string            `json:"-"`
	IngressClassName string            `json:"-"`
}

// EnvVarName is allowed/disalloed for environment variables
type EnvVarName string

const (
	EnvVarNameSharingioPairExposerDisabled      EnvVarName = "SHARINGIO_PAIR_EXPOSER_DISABLED"
	EnvVarNameSharingioPairExposerHostname      EnvVarName = "SHARINGIO_PAIR_EXPOSER_HOSTNAME"
	EnvVarNameSharingioPairExposerAllowedPorts  EnvVarName = "SHARINGIO_PAIR_EXPOSER_ALLOWED_PORTS"
	EnvVarNameSharingioPairExposerDisabledPorts EnvVarName = "SHARINGIO_PAIR_EXPOSER_DISABLED_PORTS"

	// deprecated
	EnvVarNameSharingioPairSetHostname                    EnvVarName = "SHARINGIO_PAIR_SET_HOSTNAME"
	EnvVarNameSharingioPairIngressReconcilerAllowedPorts  EnvVarName = "SHARINGIO_PAIR_INGRESS_RECONCILER_ALLOWED_PORTS"
	EnvVarNameSharingioPairIngressReconcilerDisabledPorts EnvVarName = "SHARINGIO_PAIR_INGRESS_RECONCILER_DISABLED_PORTS"
)

const DotSharingDotIoExposerTemplatesFolderPath string = "/home/ii/.sharing.io/cluster-api/exposer/templates"

// TemplateYamlFile is the name for YAML files that are included
type TemplateYamlFile string

const (
	TemplateYamlFileService              TemplateYamlFile = "service.yaml"
	TemplateYamlFileIngress              TemplateYamlFile = "ingress.yaml"
	TemplateYamlFileIngressV118OrEarlier TemplateYamlFile = "ingress-v1.18-or-earlier.yaml"
)

// allowed listening IP address
const IPAllInterfaces = "0.0.0.0"
const IPv6AllInterfaces = "::"

type ResourceLabel string

const resourceLabelRTLD string = "io.sharing.pair"

const (
	ResourceLabelName    ResourceLabel = ResourceLabel(resourceLabelRTLD + "/" + "name")
	ResourceLabelPort    ResourceLabel = ResourceLabel(resourceLabelRTLD + "/" + "port")
	ResourceLabelUid     ResourceLabel = ResourceLabel(resourceLabelRTLD + "/" + "uid")
	ResourceLabelManaged ResourceLabel = ResourceLabel(resourceLabelRTLD + "/" + "managed")
)

// ResourceLabels is a map for the resource labels to track Kubernetes resources with
var ResourceLabels = map[string]string{string(ResourceLabelManaged): "true"}

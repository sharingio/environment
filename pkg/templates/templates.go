/*

templates
  package for parsing and rendering local YAML files

*/

package templates

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
	"time"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"sigs.k8s.io/yaml"

	"github.com/sharingio/environment/pkg/common"
	"github.com/sharingio/environment/pkg/types"
)

// Template is a struct for storing what's needed for rendering
type Template struct {
	process types.Process
}

// NewTemplate returns a populated Template, given a Process
func NewTemplate(process types.Process) (template *Template) {
	return &Template{
		process: process,
	}
}

// RenderYAMLWithProcess parses and evaluates a Go-Template style file, returning a string
func (t Template) RenderYAMLWithProcess(input string) (output string, err error) {
	tmpl, err := template.New(fmt.Sprintf("yaml-render-%v", time.Now().Unix())).Parse(input)
	if err != nil {
		return "", err
	}
	templatedBuffer := new(bytes.Buffer)
	err = tmpl.Execute(templatedBuffer, t.process)
	if err != nil {
		log.Printf("Error templating YAML: %v\n", err.Error())
		return output, fmt.Errorf("Error templating YAML: %v", err.Error())
	}
	output = templatedBuffer.String()
	return output, nil
}

// RenderService returns a Service, based on a template in .sharing.io
func (t Template) RenderService() (service v1.Service, err error) {
	serviceRawString, err := common.ReadFile(types.DotSharingDotIoExposerTemplatesFolderPath + "/" + string(types.TemplateYamlFileService))
	if err != nil {
		return v1.Service{}, err
	}
	serviceString, err := t.RenderYAMLWithProcess(serviceRawString)
	err = yaml.Unmarshal([]byte(serviceString), &service)
	if err != nil {
		return v1.Service{}, err
	}
	service.ObjectMeta.Labels = t.AddCommonLabels(service.ObjectMeta.Labels)
	return service, nil
}

// RenderService returns an Ingress, based on a template in .sharing.io
func (t Template) RenderIngress() (ingress networkingv1.Ingress, err error) {
	ingressRawString, err := common.ReadFile(types.DotSharingDotIoExposerTemplatesFolderPath + "/" + string(types.TemplateYamlFileIngress))
	if err != nil {
		return networkingv1.Ingress{}, err
	}
	ingressString, err := t.RenderYAMLWithProcess(ingressRawString)
	err = yaml.Unmarshal([]byte(ingressString), &ingress)
	if err != nil {
		return networkingv1.Ingress{}, err
	}
	ingress.ObjectMeta.Labels = t.AddCommonLabels(ingress.ObjectMeta.Labels)
	return ingress, nil
}

// RenderService returns an Ingress v1beta1, based on a template in .sharing.io
func (t Template) RenderIngressv1beta1() (ingress networkingv1beta1.Ingress, err error) {
	ingressRawString, err := common.ReadFile(types.DotSharingDotIoExposerTemplatesFolderPath + "/" + string(types.TemplateYamlFileIngressV118OrEarlier))
	if err != nil {
		return networkingv1beta1.Ingress{}, err
	}
	ingressString, err := t.RenderYAMLWithProcess(ingressRawString)
	err = yaml.Unmarshal([]byte(ingressString), &ingress)
	if err != nil {
		return networkingv1beta1.Ingress{}, err
	}
	ingress.ObjectMeta.Labels = t.AddCommonLabels(ingress.ObjectMeta.Labels)
	return ingress, nil
}

// AddCommonLabels adds exposer specific labels
func (t Template) AddCommonLabels(input map[string]string) map[string]string {
	if input == nil {
		input = map[string]string{}
	}
	input[string(types.ResourceLabelName)] = fmt.Sprintf("%v", t.process.Name)
	input[string(types.ResourceLabelPort)] = fmt.Sprintf("%v", t.process.Port)
	input[string(types.ResourceLabelUid)] = fmt.Sprintf("%v", t.process.Uid)
	input[string(types.ResourceLabelManaged)] = "true"
	return input
}

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

type Template struct {
	process types.Process
}

func NewTemplate(process types.Process) (template *Template) {
	return &Template{
		process: process,
	}
}

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
	return service, nil
}

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
	return ingress, nil
}

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
	return ingress, nil
}

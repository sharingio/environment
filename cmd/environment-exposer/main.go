package main

import (
	"fmt"
	"strconv"

	"github.com/sharingio/environment/pkg/common"
	"github.com/sharingio/environment/pkg/exposer"
	"github.com/sharingio/environment/pkg/templates"
	"github.com/sharingio/environment/pkg/types"
)

func main() {
	e, err := exposer.NewExposer()
	if err != nil {
		fmt.Println("Failed to get NewExposer", err)
		return
	}

	rm := exposer.NewResourceManager(e)
	rm.Namespace = common.GetAppEnvironmentNamespace()

	kVersion, err := e.Clientset.DiscoveryClient.ServerVersion()
	if err != nil {
		fmt.Println("Failed to Kubernetes ersion", err)
		return
	}
	kVersionMajor, err := strconv.Atoi(kVersion.Major)
	if err != nil {
		fmt.Println("Failed to convert Kubernetes major version to int", err)
		return
	}
	kVersionMinor, err := strconv.Atoi(kVersion.Minor)
	if err != nil {
		fmt.Println("Failed to convert Kubernetes minor version to int", err)
		return
	}

	listening, err := e.GetListening()
	if err != nil {
		fmt.Println("Failed to GetListening", err)
		return
	}

	for _, l := range listening {
		l.ServiceName = l.Name
		// ofset the number to ensure ports like 80 or 443 aren't overtaken if locally bound
		l.ServicePort = l.Port
		if l.ServicePort < 1000 {
			l.ServicePort = 10000 + l.Port
		}
		l.ExternalIP = common.GetAppExternalIP()
		l.IngressHost = fmt.Sprintf("%v.%v", l.ServiceName, e.IngressBaseDomain)

		tmpl := templates.NewTemplate(l)
		svc, err := tmpl.RenderService()
		if err != nil {
			fmt.Printf("Failed to render Service: %v\n", err)
		}

		ing, err := tmpl.RenderIngress()
		if err != nil {
			fmt.Printf("Failed to render Ingress: %v\n", err)
		}

		ingv1beta1, err := tmpl.RenderIngressv1beta1()
		if err != nil {
			fmt.Printf("Failed to render Ingress: %v\n", err)
		}
		fmt.Println("v1.Service", svc)
		err = rm.CreateOrUpdateService(&svc)
		if err != nil {
			fmt.Printf("Failed to create Service: %v\n", err)
		}

		if l.Protocol != types.ProtocolTCP {
			continue
		}
		if kVersionMajor == 1 && kVersionMinor > 18 {
			fmt.Println("networkingv1.Ingress", ing)
			err = rm.CreateOrUpdateIngress(&ing)
			if err != nil {
				fmt.Printf("Failed to create Ingress: %v\n", err)
			}
		}
		if kVersionMajor == 1 && kVersionMinor <= 18 {
			fmt.Println("networkingv1beta1.Ingress", ingv1beta1)
			err = rm.CreateOrUpdateIngressV1beta1(&ingv1beta1)
			if err != nil {
				fmt.Printf("Failed to create Ingress v1beta1: %v\n", err)
			}
		}
	}

	// TODO prune unused services
}

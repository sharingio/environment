package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/sharingio/environment/pkg/common"
	k "github.com/sharingio/environment/pkg/kubernetes"
	"github.com/sharingio/environment/pkg/templates"
	"github.com/sharingio/environment/pkg/types"

	"k8s.io/client-go/kubernetes"
)

type Exposer struct {
	exporterEndpoint string
	clientset        *kubernetes.Clientset
}

func NewExposer() (exposer *Exposer, err error) {
	clientset, err := k.NewClient()
	if err != nil {
		return nil, err
	}
	return &Exposer{
		exporterEndpoint: common.GetAppExporterEndpoint(),
		clientset:        clientset,
	}, nil
}

func (e *Exposer) GetListening() (listening []types.Process, err error) {
	req, err := http.NewRequest(http.MethodGet, e.exporterEndpoint+"/listening", nil)
	if err != nil {
		return []types.Process{}, err
	}
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []types.Process{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []types.Process{}, err
	}

	err = json.Unmarshal(body, &listening)
	if err != nil {
		return []types.Process{}, err
	}

	return listening, nil
}

func main() {
	exposer, err := NewExposer()
	if err != nil {
		fmt.Println("Failed to get NewExposer", err)
		return
	}
	kVersion, err := exposer.clientset.DiscoveryClient.ServerVersion()
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

	listening, err := exposer.GetListening()
	if err != nil {
		fmt.Println("Failed to GetListening", err)
		return
	}

	for _, l := range listening {
		l.ServiceName = l.Name
		// ofset the number to ensure ports like 80 or 443 aren't overtaken if locally bound
		if l.ServicePort > 1000 {
			l.ServicePort = 10000 + l.Port
		}
		l.ExternalIP = common.GetAppExternalIP()

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
		err = k.CreateOrUpdateService(svc)

		if l.Protocol != types.ProtocolTCP {
			continue
		}
		if kVersionMajor == 1 && kVersionMinor > 18 {
			fmt.Println("networkingv1.Ingress", ing)
			err = k.CreateOrUpdateIngress(ing)
		}
		if kVersionMajor == 1 && kVersionMinor <= 18 {
			fmt.Println("networkingv1beta1.Ingress", ingv1beta1)
			err = k.CreateOrUpdateIngressV1Beta1(ingv1beta1)
		}
	}

	// TODO prune unused services
}

/*

environment-exposer
  A controller for mananging Kubernetes Service and Ingress resources based on the the results of environment-exporter, in order to expose the ports inside Environment to the public internet.

*/

package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/sharingio/environment/pkg/common"
	"github.com/sharingio/environment/pkg/exposer"
	"github.com/sharingio/environment/pkg/templates"
	"github.com/sharingio/environment/pkg/types"

	"github.com/joho/godotenv"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

var (
	AppBuildVersion = "0.0.0"
	AppBuildHash    = "???"
	AppBuildDate    = "yyyy.mm.dd HH:MM"
)

func main() {
	log.Printf("launching sharingio/environment:exposer (%v, %v, %v)\n", AppBuildVersion, AppBuildHash, AppBuildDate)

	envFile := common.GetAppEnvFile()
	_ = godotenv.Load(envFile)
	e, err := exposer.NewExposer()
	if err != nil {
		log.Println("Failed to get NewExposer", err)
		return
	}

	rm := exposer.NewResourceManager(e)
	rm.Namespace = common.GetAppEnvironmentNamespace()

	kVersion, err := e.Clientset.DiscoveryClient.ServerVersion()
	if err != nil {
		log.Println("Failed to Kubernetes version", err)
		return
	}
	kVersionMajor, err := strconv.Atoi(kVersion.Major)
	if err != nil {
		log.Println("Failed to convert Kubernetes major version to int", err)
		return
	}
	kVersionMinor, err := strconv.Atoi(kVersion.Minor)
	if err != nil {
		log.Println("Failed to convert Kubernetes minor version to int", err)
		return
	}

	for {
		time.Sleep(time.Duration(e.ReconciliationInterval))

		listening, err := e.GetListening()
		if err != nil {
			log.Println("Failed to GetListening", err)
			continue
		}
		var listeningNames []string

	listenList:
		for _, l := range listening {
			l.ServiceName = l.Name
			if l.Hostname != "" {
				l.ServiceName = l.Hostname
			}

			// use {{ ServiceName }}-{{ Port }} as name if there's already a service with the name {{ ServiceName }}
			svcs, err := rm.GetServices()
			if err != nil {
				log.Printf("Failed to list Services: %v\n", err)
				continue listenList
			}
			serviceWithNameAlreadyExists := false
		existingServiceList:
			for _, svc := range svcs.Items {
				if svc.ObjectMeta.Labels[string(types.ResourceLabelName)] == fmt.Sprintf("%v", l.Name) &&
					svc.ObjectMeta.Labels[string(types.ResourceLabelPort)] != fmt.Sprintf("%v", l.Port) {
					serviceWithNameAlreadyExists = true
					break existingServiceList
				}
				break existingServiceList
			}
			if serviceWithNameAlreadyExists == true {
				l.ServiceName = fmt.Sprintf("%v-%v", l.ServiceName, l.Port)
			}

			// offset the number to ensure ports like 80 or 443 aren't overtaken if locally bound
			l.ServicePort = l.Port
			if l.ServicePort < 1000 {
				l.ServicePort = 10000 + l.Port
			}
			l.ExternalIP = common.GetAppExternalIP()
			l.IngressClassName = common.GetAppIngressClassName()
			l.IngressHost = fmt.Sprintf("%v.%v", l.ServiceName, e.IngressBaseDomain)

			// prefix resource name if provided
			l.ResourceName = l.ServiceName
			if e.ResourceNamePrefix != "" {
				l.ResourceName = fmt.Sprintf("%v-%v", e.ResourceNamePrefix, l.ServiceName)
			}
			listeningNames = append(listeningNames, l.ResourceName)

			tmpl := templates.NewTemplate(l)
			svc, err := tmpl.RenderService()
			if err != nil {
				log.Printf("Failed to render Service: %v\n", err)
				continue listenList
			}

			ing, err := tmpl.RenderIngress()
			if err != nil {
				log.Printf("Failed to render Ingress: %v\n", err)
				continue listenList
			}

			ingv1beta1, err := tmpl.RenderIngressv1beta1()
			if err != nil {
				log.Printf("Failed to render Ingress: %v\n", err)
				continue listenList
			}
			err = rm.CreateOrUpdateService(&svc)
			if err != nil && apierrors.IsAlreadyExists(err) == false {
				log.Printf("Failed to create Service '%v' in namespace '%v': %v\n", svc.ObjectMeta.Name, rm.Namespace, err)
			} else if apierrors.IsAlreadyExists(err) == false {
				log.Printf("Created v1.Service '%v' in namespace '%v'\n", svc.ObjectMeta.Name, rm.Namespace)
			}

			if l.Protocol != types.ProtocolTCP {
				continue listenList
			}
			if kVersionMajor == 1 && kVersionMinor > 18 {
				err = rm.CreateOrUpdateIngress(&ing)
				if err != nil && apierrors.IsAlreadyExists(err) == false {
					log.Printf("Failed to create Ingress in '%v' namespace '%v': %v\n", svc.ObjectMeta.Name, rm.Namespace, err)
				} else if apierrors.IsAlreadyExists(err) == false {
					log.Printf("Created networkingv1.Ingress '%v' in namespace '%v'\n", ing.ObjectMeta.Name, rm.Namespace)
				}
			} else if kVersionMajor == 1 && kVersionMinor <= 18 {
				err = rm.CreateOrUpdateIngressV1beta1(&ingv1beta1)
				if err != nil && apierrors.IsAlreadyExists(err) == false {
					log.Printf("Failed to create Ingress v1beta1 '%v' in namespace '%v': %v\n", svc.ObjectMeta.Name, rm.Namespace, err)
				} else if apierrors.IsAlreadyExists(err) == false {
					log.Printf("Created networkingv1beta1.Ingress '%v' in namespace '%v'\n", ingv1beta1.ObjectMeta.Name, rm.Namespace)
				}
			}
		}

		deleted, err := rm.PruneUnusedServices(listeningNames)
		if err != nil {
			log.Printf("Failed to prune unused Services in namespace '%v': %v\n", rm.Namespace, err)
		}
		if len(deleted) > 0 {
			log.Printf("Deleted Service(s) %v in namespace '%v'\n", strings.Join(deleted, ", "), rm.Namespace)
		}
		if kVersionMajor == 1 && kVersionMinor > 18 {
			deleted, err = rm.PruneUnusedIngresses(listeningNames)
			if err != nil {
				log.Printf("Failed to prune unused Ingress(es) in namespace '%v': %v\n", rm.Namespace, err)
			}
			if len(deleted) > 0 {
				log.Printf("Deleted Ingress(es) %v in namespace '%v'\n", strings.Join(deleted, ", "), rm.Namespace)
			}
		} else if kVersionMajor == 1 && kVersionMinor <= 18 {
			deleted, err = rm.PruneUnusedIngressesV1beta1(listeningNames)
			if err != nil {
				log.Printf("Failed to prune unused Ingresses v1beta1 in namespace: %v\n", rm.Namespace, err)
			}
			if len(deleted) > 0 {
				log.Printf("Deleted Ingresses v1beta1 %v in namespace '%v'\n", strings.Join(deleted, ", "), rm.Namespace)
			}
		}
	}
}

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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func main() {
	e, err := exposer.NewExposer()
	if err != nil {
		log.Println("Failed to get NewExposer", err)
		return
	}

	rm := exposer.NewResourceManager(e)
	rm.Namespace = common.GetAppEnvironmentNamespace()

	kVersion, err := e.Clientset.DiscoveryClient.ServerVersion()
	if err != nil {
		log.Println("Failed to Kubernetes ersion", err)
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
		listening, err := e.GetListening()
		if err != nil {
			log.Println("Failed to GetListening", err)
			return
		}
		var listeningNames []string

	listenList:
		for _, l := range listening {

			l.ServiceName = fmt.Sprintf("%v-%v", l.Name, l.Port)
			listeningNames = append(listeningNames, l.ServiceName)
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
				log.Printf("Failed to render Service: %v\n", err)
			}

			ing, err := tmpl.RenderIngress()
			if err != nil {
				log.Printf("Failed to render Ingress: %v\n", err)
			}

			ingv1beta1, err := tmpl.RenderIngressv1beta1()
			if err != nil {
				log.Printf("Failed to render Ingress: %v\n", err)
			}
			err = rm.CreateOrUpdateService(&svc)
			if err != nil && apierrors.IsAlreadyExists(err) == false {
				log.Printf("Failed to create Service: %v\n", err)
			} else if apierrors.IsAlreadyExists(err) == false {
				log.Printf("Created v1.Service '%v' in namespace '%v'\n", svc.ObjectMeta.Name, rm.Namespace)
			}

			if l.Protocol != types.ProtocolTCP {
				continue listenList
			}
			if kVersionMajor == 1 && kVersionMinor > 18 {
				err = rm.CreateOrUpdateIngress(&ing)
				if err != nil && apierrors.IsAlreadyExists(err) == false {
					log.Printf("Failed to create Ingress: %v\n", err)
				} else if apierrors.IsAlreadyExists(err) == false {
					log.Printf("Created networkingv1.Ingress '%v' in namespace '%v'\n", ing.ObjectMeta.Name, rm.Namespace)
				}
			} else if kVersionMajor == 1 && kVersionMinor <= 18 {
				err = rm.CreateOrUpdateIngressV1beta1(&ingv1beta1)
				if err != nil && apierrors.IsAlreadyExists(err) == false {
					log.Printf("Failed to create Ingress v1beta1: %v\n", err)
				} else if apierrors.IsAlreadyExists(err) == false {
					log.Printf("Created networkingv1beta1.Ingress '%v' in namespace '%v'\n", ingv1beta1.ObjectMeta.Name, rm.Namespace)
				}
			}
		}

		deleted, err := rm.PruneUnusedServices(listeningNames)
		if err != nil {
			log.Printf("Failed to prune unused Services: %v\n", err)
		}
		if len(deleted) > 0 {
			log.Printf("Deleted Services %v\n", strings.Join(deleted, ", "))
		}
		if kVersionMajor == 1 && kVersionMinor > 18 {
			deleted, err = rm.PruneUnusedIngresses(listeningNames)
			if err != nil {
				log.Printf("Failed to prune unused Ingresses: %v\n", err)
			}
			if len(deleted) > 0 {
				log.Printf("Deleted Ingresses %v\n", strings.Join(deleted, ", "))
			}
		} else if kVersionMajor == 1 && kVersionMinor <= 18 {
			deleted, err = rm.PruneUnusedIngressesV1beta1(listeningNames)
			if err != nil {
				log.Printf("Failed to prune unused Ingresses v1beta1: %v\n", err)
			}
			if len(deleted) > 0 {
				log.Printf("Deleted Ingresses v1beta1 %v\n", strings.Join(deleted, ", "))
			}
		}
		time.Sleep(time.Duration(e.ReconciliationInterval))
	}
}

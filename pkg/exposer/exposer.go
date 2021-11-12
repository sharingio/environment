package exposer

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/sharingio/environment/pkg/common"
	k "github.com/sharingio/environment/pkg/kubernetes"
	"github.com/sharingio/environment/pkg/types"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

var resourceLabelSelector = labels.SelectorFromSet(types.ResourceLabels).String()

type Exposer struct {
	ExporterEndpoint       string
	IngressBaseDomain      string
	ReconciliationInterval time.Duration
	Clientset              *kubernetes.Clientset
}

func NewExposer() (exposer *Exposer, err error) {
	exporterEndpoint := common.GetAppExporterEndpoint()
	ingressBaseDomain := common.GetAppEnvironmentBaseDomain()
	reconciliationInterval := common.GetAppReconciliationInterval()
	clientset, err := k.NewClient()
	if err != nil {
		return nil, err
	}
	return &Exposer{
		ExporterEndpoint:       exporterEndpoint,
		IngressBaseDomain:      ingressBaseDomain,
		ReconciliationInterval: reconciliationInterval,
		Clientset:              clientset,
	}, nil
}

func (e *Exposer) GetListening() (listening []types.Process, err error) {
	req, err := http.NewRequest(http.MethodGet, e.ExporterEndpoint+"/listening", nil)
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

type ResourceManager struct {
	Namespace string
	clientset *kubernetes.Clientset
}

func NewResourceManager(exposer *Exposer) (resourceManager ResourceManager) {
	return ResourceManager{
		clientset: exposer.Clientset,
	}
}

func (r ResourceManager) GetServices() (services *v1.ServiceList, err error) {
	services, err = r.clientset.CoreV1().Services(r.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: resourceLabelSelector})
	return services, err
}

func (r ResourceManager) GetIngresses() (ingresses *networkingv1.IngressList, err error) {
	ingresses, err = r.clientset.NetworkingV1().Ingresses(r.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: resourceLabelSelector})
	return ingresses, nil
}

func (r ResourceManager) GetIngressesV1beta1() (ingresses *networkingv1beta1.IngressList, err error) {
	ingresses, err = r.clientset.NetworkingV1beta1().Ingresses(r.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: resourceLabelSelector})
	return ingresses, nil
}

func (r ResourceManager) PruneUnusedServices(names []string) (deletedNames []string, err error) {
	services, err := r.GetServices()
	if err != nil {
		return []string{}, err
	}
	for _, service := range services.Items {
		nameFoundInServiceList := false
		for _, name := range names {
			if name == service.ObjectMeta.Name {
				nameFoundInServiceList = true
			}
		}
		if nameFoundInServiceList == false {
			deletedNames = append(deletedNames, service.ObjectMeta.Name)
		}
	}
	for _, name := range deletedNames {
		err = r.clientset.CoreV1().Services(r.Namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			return []string{}, err
		}
	}
	return deletedNames, nil
}

func (r ResourceManager) PruneUnusedIngresses(names []string) (deletedNames []string, err error) {
	ingresses, err := r.GetIngresses()
	if err != nil {
		return []string{}, err
	}
	for _, ingress := range ingresses.Items {
		nameFoundInIngressList := false
		for _, name := range names {
			if name == ingress.ObjectMeta.Name {
				nameFoundInIngressList = true
			}
		}
		if nameFoundInIngressList == false {
			deletedNames = append(deletedNames, ingress.ObjectMeta.Name)
		}
	}
	for _, name := range deletedNames {
		err = r.clientset.NetworkingV1().Ingresses(r.Namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			return []string{}, err
		}
	}
	return deletedNames, nil
}

func (r ResourceManager) PruneUnusedIngressesV1beta1(names []string) (deletedNames []string, err error) {
	ingresses, err := r.GetIngresses()
	if err != nil {
		return []string{}, err
	}
	for _, ingress := range ingresses.Items {
		nameFoundInIngressList := false
		for _, name := range names {
			if name == ingress.ObjectMeta.Name {
				nameFoundInIngressList = true
			}
		}
		if nameFoundInIngressList == false {
			deletedNames = append(deletedNames, ingress.ObjectMeta.Name)
		}
	}
	for _, name := range deletedNames {
		err = r.clientset.NetworkingV1beta1().Ingresses(r.Namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			return []string{}, err
		}
	}
	return deletedNames, nil
}

// TODO implement update
func (r ResourceManager) CreateOrUpdateService(service *v1.Service) (err error) {
	_, err = r.clientset.CoreV1().Services(r.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (r ResourceManager) CreateOrUpdateIngress(ingress *networkingv1.Ingress) (err error) {
	_, err = r.clientset.NetworkingV1().Ingresses(r.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (r ResourceManager) CreateOrUpdateIngressV1beta1(ingress *networkingv1beta1.Ingress) (err error) {
	_, err = r.clientset.NetworkingV1beta1().Ingresses(r.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

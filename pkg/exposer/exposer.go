package exposer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sharingio/environment/pkg/common"
	k "github.com/sharingio/environment/pkg/kubernetes"
	"github.com/sharingio/environment/pkg/types"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

var resourceLabelSelector = labels.SelectorFromSet(types.ResourceLabels).String()

type Exposer struct {
	ExporterEndpoint  string
	IngressBaseDomain string
	Clientset         *kubernetes.Clientset
}

func NewExposer() (exposer *Exposer, err error) {
	clientset, err := k.NewClient()
	if err != nil {
		return nil, err
	}
	exporterEndpoint := common.GetAppExporterEndpoint()
	ingressBaseDomain := common.GetAppEnvironmentBaseDomain()
	return &Exposer{
		Clientset:         clientset,
		ExporterEndpoint:  exporterEndpoint,
		IngressBaseDomain: ingressBaseDomain,
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

// NOTE listOptions include a labelSelector as "io.sharing.pair/managed=true"
// NOTE serviceName must be {name}-{port} (not ServicePort)

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

func (r ResourceManager) GetIngress() (ingresses *networkingv1.IngressList, err error) {
	ingresses, err = r.clientset.NetworkingV1().Ingresses(r.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: resourceLabelSelector})
	return ingresses, nil
}

func (r ResourceManager) GetIngressV1Beta1() (ingresses *networkingv1beta1.IngressList, err error) {
	ingresses, err = r.clientset.NetworkingV1beta1().Ingresses(r.Namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: resourceLabelSelector})
	return ingresses, nil
}

func (r ResourceManager) PruneUnusedServices() (err error) {
	return nil
}

func (r ResourceManager) PruneUnusedIngresses() (err error) {
	return nil
}

func (r ResourceManager) PruneUnusedIngressesV1Beta1() (err error) {
	return nil
}

// TODO implement update
func (r ResourceManager) CreateOrUpdateService(service *v1.Service) (err error) {
	_, err = r.clientset.CoreV1().Services(r.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil && apierrors.IsAlreadyExists(err) == false {
		return fmt.Errorf("Error creating Service, %v", err.Error())
	}
	return nil
}

func (r ResourceManager) CreateOrUpdateIngress(ingress *networkingv1.Ingress) (err error) {
	_, err = r.clientset.NetworkingV1().Ingresses(r.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil && apierrors.IsAlreadyExists(err) == false {
		return fmt.Errorf("Error creating Ingress, %v", err.Error())
	}
	return nil
}

func (r ResourceManager) CreateOrUpdateIngressV1beta1(ingress *networkingv1beta1.Ingress) (err error) {
	_, err = r.clientset.NetworkingV1beta1().Ingresses(r.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil && apierrors.IsAlreadyExists(err) == false {
		return fmt.Errorf("Error creating Ingress, %v", err.Error())
	}
	return nil
}

package kubernetes

import (
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// user agent for Kubernetes APIServer
var defaultKubeClientUserAgent = "sharingio/environment/cmd/environment-exposer"

// NewClient ...
// return a Kubernetes clientset
func NewClient() (clientset *kubernetes.Clientset, err error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	config.UserAgent = defaultKubeClientUserAgent
	config.QPS = 500
	config.Burst = 1000
	clientset, err = kubernetes.NewForConfig(config)
	return clientset, nil
}

// NOTE listOptions include a labelSelector as "io.sharing.pair/managed=true"
// NOTE serviceName must be {name}-{port} (not ServicePort)

func GetServices() (services []v1.Service, err error) {
	return services, nil
}

func GetIngresses() (ingresses []networkingv1.Ingress, err error) {
	return ingresses, nil
}

func GetIngressesV1Beta1() (ingresses []networkingv1beta1.Ingress, err error) {
	return ingresses, nil
}

func PruneUnusedServices() (err error) {
	return nil
}

func PruneUnusedIngresses() (err error) {
	return nil
}

func PruneUnusedIngressesV1Beta1() (err error) {
	return nil
}

func CreateOrUpdateService(service v1.Service) (err error) {
	return nil
}

func CreateOrUpdateIngress(ingress networkingv1.Ingress) (err error) {
	return nil
}

func CreateOrUpdateIngressV1Beta1(ingress networkingv1beta1.Ingress) (err error) {
	return nil
}

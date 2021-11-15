/*

kubernetes
  functions used for talking to Kubernetes

*/

package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// user agent for Kubernetes APIServer
var defaultKubeClientUserAgent = "sharingio/environment/cmd/environment-exposer"

// NewClient returns a Kubernetes clientset
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

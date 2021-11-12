package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// GetEnvOrDefault ...
// return env value or default to value
func GetEnvOrDefault(envName string, defaultValue string) (output string) {
	output = os.Getenv(envName)
	if output == "" {
		output = defaultValue
	}
	return output
}

// GetAppEnvFile ...
// location of an env file to load
func GetAppEnvFile() (output string) {
	return GetEnvOrDefault("APP_ENV_FILE", ".env")
}

// GetAppPort ...
// the port to bind to
func GetAppPort() (output string) {
	return GetEnvOrDefault("APP_PORT", "localhost:10093")
}

// GetPodLabelsFilePath ...
// the path to a downward API generated file containing the defined Pod labels
func GetPodLabelsFilePath() (output string) {
	return GetEnvOrDefault("APP_POD_LABELS_FILE_PATH", "/etc/podlabels/labels")
}

// GetAppExternalIP
// the externalIP for the cluster
// this is only used in exposer
func GetAppExternalIP() (output string) {
	return GetEnvOrDefault("APP_EXTERNAL_IP", GetEnvOrDefault("SHARINGIO_PAIR_LOAD_BALANCER_IP", ""))
}

// GetAppExporterEndpoint
// the HTTP URI for the exporter
// e.g: http://environment-exporter.default:10093
// this is only used in exposer
func GetAppExporterEndpoint() (output string) {
	return GetEnvOrDefault("APP_EXPORTER_ENDPOINT", "http://environment-exporter.default:10093")
}

// GetAppEnvironmentNamespace ...
// the namespace where Environment is deployed
// this is only used in the exposer
func GetAppEnvironmentNamespace() (output string) {
	return GetEnvOrDefault("APP_ENVIRONMENT_NAMESPACE", "default")
}

// GetAppEnvironmentBaseDomain ...
// the base domain to create ingresses with
// this is only used in the exposer
func GetAppEnvironmentBaseDomain() (output string) {
	return GetEnvOrDefault("APP_ENVIRONMENT_BASE_DOMAIN", GetEnvOrDefault("SHARINGIO_PAIR_BASE_DNS_NAME", "_.pair.sharing.io"))
}

// GetAppReconciliationInterval ...
// the base domain to create ingresses with
// this is only used in the exposer
func GetAppReconciliationInterval() (seconds time.Duration) {
	defaultSeconds := time.Duration(2 * time.Second)
	secondsString := GetEnvOrDefault("APP_RECONCILIATION_INTERVAL", fmt.Sprintf("%v", defaultSeconds.Seconds()))
	secondsInt, err := strconv.Atoi(secondsString)
	seconds = time.Duration(secondsInt) * time.Second
	if err != nil {
		log.Printf("Failed to convert string to int: %v\n", err)
	}
	if seconds <= 0 {
		seconds = defaultSeconds
	}
	log.Printf("Will sleep for %v seconds per iteration", seconds.Seconds())
	return seconds
}

// GetPodName ...
// the name of the Pod
// this is only used in the exporter
func GetPodName() (output string) {
	return GetEnvOrDefault("POD_NAME", GetEnvOrDefault("HOSTNAME", ""))
}

// Logging ...
// basic request logging middleware
func Logging(next http.Handler) http.Handler {
	// log all requests
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v %v %v %v", r.Method, r.URL, r.Proto, r.Response, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func ReadFile(name string) (contents string, err error) {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

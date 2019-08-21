package services

import (
	"fmt"

	"github.com/azohra/periscope/internal/pkg/periscope/tools"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // k8s auth
	"k8s.io/client-go/rest"
)

func getClientset() (clientset *kubernetes.Clientset, err error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err.Error())
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err.Error())
	}
	return clientset, err
}

// CreateDeployment create new deployment
func CreateDeployment(appName string, namespace string, containers []apiv1.Container, pullSecret string) (succeess bool) {

	clientset, err := getClientset()
	if err != nil {
		succeess = false
		return
	}

	fmt.Printf("Creating Deployments for %s\n", appName)
	deploymentsClient := clientset.AppsV1().Deployments(namespace)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: tools.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": appName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":       appName,
						"namespace": namespace,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: containers,
					ImagePullSecrets: []apiv1.LocalObjectReference{
						{
							Name: pullSecret,
						},
					},
				},
			},
		},
	}

	_, err = deploymentsClient.Create(deployment)
	if err != nil {
		fmt.Printf("Deployment %s launch failed\n", appName)
		fmt.Println(err.Error())
		succeess = false
	} else {
		fmt.Printf("Deployment %s launched successfully\n", appName)
		succeess = true
	}
	return
}

// CreateService create new service
func CreateService(appName string, namespace string, port int32) (succeess bool) {

	clientset, err := getClientset()
	if err != nil {
		succeess = false
		return
	}

	fmt.Printf("Creating Services for %s\n", appName)
	servicesClient := clientset.CoreV1().Services(namespace)
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appName,
			Namespace: namespace,
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{"app": appName},
			Ports: []apiv1.ServicePort{
				{
					Protocol: "TCP",
					Port:     port,
				},
			},
		},
	}

	fmt.Printf("Creating Services for %s\n", appName)
	_, err = servicesClient.Create(service)
	if err != nil {
		succeess = false
	} else {
		fmt.Printf("Service %s launched successfully!\n", appName)
		succeess = true
	}
	return
}

// CreateVirtualService create new virtualservice
func CreateVirtualService(appName string, namespace string, stateIDStr string, port int, versionStr string) (succeess bool) {

	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err.Error())
		succeess = false
		return
	}

	fmt.Printf("Creating VirtualService for %s\n", appName)
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Println(err.Error())
		succeess = false
		return
	}

	virtualServiceGVR := schema.GroupVersionResource{
		Group:    "networking.istio.io",
		Version:  "v1alpha3",
		Resource: "virtualservices",
	}

	gateways := []string{fmt.Sprintf("%s-gateway.%s", namespace, namespace)}
	hosts := []string{"*"}

	_, err = dynamicClient.Resource(virtualServiceGVR).Namespace(namespace).Create(
		templateVirtualService(
			"networking.istio.io/v1alpha3", // version
			"VirtualService",               // kind
			appName,                        // name
			namespace,                      // namespace
			gateways,                       // gateways
			hosts,                          // hosts
			port,                           // service port
			versionStr,                     // subset
			stateIDStr,                     // state
		), metav1.CreateOptions{}, "")

	if err != nil {
		fmt.Println(err.Error())
		succeess = false
	} else {
		succeess = true
	}
	return
}

// TODO: configurable templating
func templateVirtualService(version, kind, name string, namespace string, gateways []string, hosts []string, port int, subset string, stateID string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": version,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name":      name,
				"namespace": namespace,
			},
			"spec": map[string]interface{}{
				"hosts":    hosts,
				"gateways": gateways,

				// http
				"http": []map[string]interface{}{

					// - match
					map[string]interface{}{
						"match": []map[string]interface{}{

							map[string]interface{}{
								// - headers
								"headers": map[string]interface{}{
									"X-State-ID": map[string]interface{}{
										"exact": stateID,
									},
								},
							},
						},

						"route": []map[string]interface{}{
							// - destination
							map[string]interface{}{
								"destination": map[string]interface{}{
									"host": name,

									//"subset": "v1",
									"port": map[string]int{
										"number": port,
									},
								},
							},
						},
					},
				},
			},
		}}
}

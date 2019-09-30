package main

import (
	"flag"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)


var prefix = flag.String("prefix", "", "prefix for pod/container name")
var num = flag.Int("num", 1, "number of pods/containers to create")
var image = flag.String("image", "", "docker image registry location")
var namespace = flag.String("namespace", "default", "namespace, default is default")
var port = flag.Int("port", 0, "port to expose for the service")

func main() {

	// Get kubernetes client based off of kubeconfig
	k8s, err := GetK8Client()
	if err != nil {
		log.Fatalln("Cannot find kube config")
	}

	// Create the number of pods & services required
	for start := 1; start <= *num; start++ {
		name := *prefix + strconv.Itoa(start)
		log.Println("Creating items for: ", name)
		// Create pod with container
		pod := &apiv1.Pod{
			Spec: apiv1.PodSpec{
				Containers: []apiv1.Container{
					{
						Name: name,
						Image: *image,
					},
				},
			},
		}
		// Name and Label the pod
		pod.Name = name
		pod.Labels = map[string]string{
			"app": name,
		}
		// Submit pod to kubernetes
		pod, err = k8s.CoreV1().Pods(*namespace).Create(pod)
		if err != nil {
			log.Fatalln("Cannot create pod " + name + " ", err)
		}
		// Create the service
		svc := &apiv1.Service{
			Spec: apiv1.ServiceSpec{
				Selector: map[string]string{
					"app": name,
				},
				Ports: []apiv1.ServicePort{
					{
						Name: "web",
						Port: int32(*port),
					},
				},
				Type: apiv1.ServiceTypeLoadBalancer,
			},
		}
		svc.Name = name
		svc.Labels = svc.Spec.Selector
		// Submit service to kubernetes
		svc, err = k8s.CoreV1().Services(*namespace).Create(svc)
		if err != nil {
			log.Fatalln("Cannot create service " + name + " ", err)
		}

		// This will loop forever ... kill at command line if no ip found
		for len(svc.Status.LoadBalancer.Ingress) == 0 {
			time.Sleep(time.Second * 1)
			svc, err = k8s.CoreV1().Services(*namespace).Get(name, metav1.GetOptions{})
			if err != nil {
				log.Fatalln("Unable to get service: ", )
			}
		}

		for svc.Status.LoadBalancer.Ingress[0].IP == "" {
			log.Println("No IP found for service " + name + " trying again")
			svc, err = k8s.CoreV1().Services(*namespace).Get(name, metav1.GetOptions{})
			if err != nil {
				log.Fatalln("Unable to get service: ", )
			}
		}
		log.Println(name + " IP Address: ", svc.Status.LoadBalancer.Ingress[0].IP)
	}
}

func GetK8Client() (*kubernetes.Clientset,error) {

	var config *rest.Config
	var err error

	var kConfig *string
	if home := homeDir(); home != "" {
		filepath.Join()
		kConfig = flag.String("kConfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kConfig file")
	} else {
		kConfig = flag.String("kConfig", "", "absolute path to the kConfig file")
	}
	flag.Parse()
	config, err = clientcmd.BuildConfigFromFlags("", *kConfig)
	if err != nil {
		return nil, err
	}
	// Create client set
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return clientset,err
	}

	return clientset, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		log.Println("HOME: ", h)
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

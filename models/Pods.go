package models

import (
	"fmt"
	"context"
	"bytes"
	"io"
	"os"
	"github.com/jinzhu/gorm"
	"strings"

	// appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/api/resource"

	"k8s.io/client-go/kubernetes"

	"github.com/getsentry/sentry-go"
	// beeline "github.com/honeycombio/beeline-go"
	// "github.com/sirupsen/logrus"

	"k8s.io/client-go/tools/clientcmd"
    // "k8s.io/client-go/kubernetes/scheme"
    "k8s.io/client-go/tools/remotecommand"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/client-go/rest"

	"errors"

	// "k8s.io/client-go/util/homedir"
)

const debug = false

type Pod struct {
	gorm.Model
	Name   string `json:"name"`
	Namespace string `json: "namespace"`
	UserId uint
	Command string `json:"command"`

}

type Command struct {
}


func (pod *Pod) GetPodLogs() {

	// podVar := &Pod{}

	config, err := clientcmd.BuildConfigFromFlags("", *Kubeconfig)
	if err != nil {
		sentry.WithScope(func(scope *sentry.Scope){
        	scope.SetLevel(sentry.LevelFatal)
        	sentry.CaptureException(errors.New("Kubeconfig not found, thrown in BuildConfigFromFlags in ArkPCServer.go"))
        })
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		sentry.WithScope(func(scope *sentry.Scope){
        	scope.SetLevel(sentry.LevelFatal)
        	sentry.CaptureException(errors.New("Kubeconfig error, thrown in in NewForConfig in ArkPCServer.go"))
        })
		panic(err)
	}

    ctx := context.Background()


    req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(
	    pod.Name, 
	    &corev1.PodLogOptions{},
	)

	readCloser, err := req.Stream(ctx)
	if err != nil {
	        fmt.Println("Error2: ", err)
	} else {
	        buf := new(bytes.Buffer)
	        _, err = io.Copy(buf, readCloser)
	        fmt.Println("log : ", buf.String())
	}
}

// (string, string, error)
func (pod *Pod) SendCommand() {

	config, err := GetClientConfig()
	if err != nil {
		fmt.Println(fmt.Sprint(err))
	}

	clientset, err := GetClientsetFromConfig(config)
	if err != nil {
		fmt.Println(fmt.Sprint(err))
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		fmt.Println(fmt.Sprint(err))
	}

	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&corev1.PodExecOptions{
		Command:   strings.Fields(pod.Command),
		Container: "",
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, parameterCodec)

	if debug {
		fmt.Println("Request URL:", req.URL().String())
	}

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		fmt.Println(fmt.Sprint(err))
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    true,
	})
	if err != nil {
		fmt.Println(fmt.Sprint(err))
	}

	fmt.Println(stdout.String())

	// return stdout.String(), stderr.String(), nil

}

func GetClientConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()  
	if err != nil {
		if debug {
			fmt.Printf("Unable to create config. Error: %+v\n", err)
		}
		err1 := err
		// kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", *Kubeconfig)
		if err != nil {
			err = fmt.Errorf("InClusterConfig as well as BuildConfigFromFlags Failed. Error in InClusterConfig: %+v\nError in BuildConfigFromFlags: %+v", err1, err)
			return nil, err
		}
	}

	return config, nil
}

func GetClientset() (*kubernetes.Clientset, error) {
	config, err := GetClientConfig()
	if err != nil {
		return nil, err
	}

	return GetClientsetFromConfig(config)
}

func GetClientsetFromConfig(config *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		err = fmt.Errorf("failed creating clientset. Error: %+v", err)
		return nil, err
	}

	return clientset, nil
}


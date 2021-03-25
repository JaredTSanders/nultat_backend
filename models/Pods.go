package models

import (
	"fmt"
	"context"
	"bytes"
	"io"
	"github.com/jinzhu/gorm"

	// appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/apimachinery/pkg/api/resource"

	"k8s.io/client-go/kubernetes"

	"github.com/getsentry/sentry-go"
	// beeline "github.com/honeycombio/beeline-go"
	// "github.com/sirupsen/logrus"

	"k8s.io/client-go/tools/clientcmd"

	"errors"

	// "k8s.io/client-go/util/homedir"
)

type Pod struct {
	gorm.Model
	Name   string `json:"name"`
	Namespace string `json: "namespace"`
	UserId uint
}


func (pod *Pod) GetPodLogs() {

	// podVar := &Pod{}

	fmt.Println(pod.Namespace)
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

package models

import (
	"fmt"
	u "github.com/JaredTSanders/nultat_backend/utils"
	"context"

	"github.com/jinzhu/gorm"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"k8s.io/client-go/kubernetes"

	"github.com/getsentry/sentry-go"
	// beeline "github.com/honeycombio/beeline-go"
	// "github.com/sirupsen/logrus"

	"k8s.io/client-go/tools/clientcmd"
	uuid "github.com/satori/go.uuid"
	"strings"
	"errors"

	// "k8s.io/client-go/util/homedir"
)

type ArkPCServer struct {
	gorm.Model
	Name   string `json:"name"`
	UserId uint
	ServerID string
	MapName string `json:"map"`
	ServerPass string `gorm:"-" ; json:"spass"`
	AdminPass string `gorm:"-" ; json:"apass"`
	Backup int `json:"backup"`
	Update int `json: "update"`
}

func (arkPCServer *ArkPCServer) Validate() (map[string]interface{}, bool) {

	if arkPCServer.Name == "" {
		return u.Message(false, "ArkPCServer name should be on the payload"), false
	}

	// if arkPCServer.Phone == "" {
	// 	return u.Message(false, "Phone number should be on the payload"), false
	// }

	// if arkPCServer.UserId <= 0 {
	// 	return u.Message(false, "User is not recognized"), false
	// }

	//All the required parameters are present
	return u.Message(true, "success"), true
}

func (arkPCServer *ArkPCServer) Create() map[string]interface{} {

	// fmt.Printf(r.Va)


	if resp, ok := arkPCServer.Validate(); !ok {
		return resp
	}

	config, err := clientcmd.BuildConfigFromFlags("", *Kubeconfig)
	if err != nil {
		sentry.WithScope(func(scope *sentry.Scope){
        	scope.SetLevel(sentry.LevelFatal)
        	sentry.CaptureException(errors.New("Kubeconfig not found, thrown in BuildConfigFromFlags in ArkPCServer.go"))
        })
		panic(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		sentry.WithScope(func(scope *sentry.Scope){
        	scope.SetLevel(sentry.LevelFatal)
        	sentry.CaptureException(errors.New("Kubeconfig error, throwin in NewForConfig in ArkPCServer.go"))
        })
		panic(err)
	}

	str := fmt.Sprint(arkPCServer.UserId)
	fullName := strings.ToLower(arkPCServer.Name) + "-" +  str + "-" +  strings.ToLower(arkPCServer.MapName)

	kubeclient := client.AppsV1().Deployments("default")

	kc := client.CoreV1().PersistentVolumeClaims("default")


	pvc := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "ark-"+fullName,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.PersistentVolumeAccessMode("ReadWriteOnce"),
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					"storage": *resource.NewQuantity(104857600, resource.BinarySI),
				},
			},
			StorageClassName: ptrstring("longhorn"),
		},
	}

	// Manage resource
	_, err = kc.Create(context.TODO(), pvc, metav1.CreateOptions{})

	if err != nil {
		panic(err)
	}

	// Create resource object
	object := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fullName,
			Labels: map[string]string{
				"io.kompose.service": fullName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"io.kompose.service": fullName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"io.kompose.service": fullName,
					},
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						corev1.Volume{
							Name: "ark-"+fullName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "ark-"+fullName,
								},
							},
						},
					},
					Containers: []corev1.Container{
						corev1.Container{
							Name:  fullName,
							Image: "turzam/ark",
							Ports: []corev1.ContainerPort{
								corev1.ContainerPort{
									HostPort:      0,
									ContainerPort: 7778,
									Protocol:      corev1.Protocol("UDP"),
								},
								corev1.ContainerPort{
									HostPort:      0,
									ContainerPort: 7778,
								},
								corev1.ContainerPort{
									HostPort:      0,
									ContainerPort: 27015,
									Protocol:      corev1.Protocol("UDP"),
								},
								corev1.ContainerPort{
									HostPort:      0,
									ContainerPort: 27015,
								},
								corev1.ContainerPort{
									HostPort:      0,
									ContainerPort: 32330,
								},
							},
							Env: []corev1.EnvVar{
								corev1.EnvVar{
									Name:  "ADMINPASSWORD",
									Value: arkPCServer.AdminPass,
								},
								corev1.EnvVar{
									Name:  "BACKUPONSTART",
									Value: fmt.Sprint(arkPCServer.Backup),
								},
								corev1.EnvVar{
									Name:  "GID",
									Value: "1000",
								},
								corev1.EnvVar{
									Name:  "SERVERMAP",
									Value: arkPCServer.MapName,
								},
								corev1.EnvVar{
									Name: "SERVERPASSWORD",
									Value: arkPCServer.ServerPass,
								},
								corev1.EnvVar{
									Name:  "SESSIONNAME",
									Value: arkPCServer.Name,
								},
								corev1.EnvVar{
									Name:  "UID",
									Value: "1000",
								},
								corev1.EnvVar{
									Name:  "UPDATEONSTART",
									Value: fmt.Sprint(arkPCServer.Update),
								},
							},
							Resources: corev1.ResourceRequirements{},
							VolumeMounts: []corev1.VolumeMount{
								corev1.VolumeMount{
									Name:      "ark-"+fullName,
									MountPath: "/ark",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicy("Always"),
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.DeploymentStrategyType("Recreate"),
			},
			MinReadySeconds: 0,
		},
	}

	// Manage resource
	result, err := kubeclient.Create(context.TODO(), object, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}else{
		fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	}
	fmt.Println("Deployment Created successfully!")
	

	arkPCServer.ServerID = uuid.NewV4().String()

	GetDB().Create(arkPCServer)

	resp := u.Message(true, "success")
	resp["arkPCServer"] = arkPCServer
	return resp
}

func GetArkPCServer(id uint) *ArkPCServer {

	arkPCServer := &ArkPCServer{}
	err := GetDB().Table("arkPCServers").Where("id = ?", id).First(arkPCServer).Error
	if err != nil {
		return nil
	}
	return arkPCServer
}

func GetArkPCServerStatus(id uint) *ArkPCServer {
	arkPCServer := &ArkPCServer{}
	err := GetDB().Table("arkPCServers").Where("id = ?", id).First(arkPCServer).Error
	if err != nil {
		return nil
	}
	return arkPCServer
} 

// func GetArkPCServerShell(id uint) *ArkPCServer {
// 	arkPCServer := &ArkPCServer{}
// 	err := GetDB().Table("arkPCServers").Where("id = ?", id).First(arkPCServer).Error
// 	if err != nil {
// 		return nil
// 	}
// 	return arkPCServer
// }

// func SendArkPCServerShell(id uint)  *ArkPCServer {
// 	arkPCServer := &ArkPCServer{}
// 	err := err := GetDB().Table("arkPCServers").Where("id = ?", id).First(arkPCServer).Error
	
// }

// func GetArkPCServers(user uint) []*ArkPCServer {

// 	arkPCServers := make([]*ArkPCServer, 0)
// 	err := GetDB().Table("arkPCServers").Where("user_id = ?", user).Find(&arkPCServers).Error
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil
// 	}

// 	return arkPCServers
// }

func int32Ptr(i int32) *int32 { return &i }

func ptrstring(p string) *string { return &p }
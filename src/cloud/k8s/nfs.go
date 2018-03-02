package k8s

import (
	"github.com/astaxie/beego/logs"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)


// name example-nfs
// accessModes ReadWriteMany
func CreateNfsStorageServer(param StorageParam) {
	createServiceAccount(param)
	createNfsService(param)
	createNfsStatefulSet(param)
}

// 2018-02-22 12:28
// 创建nfs服务提供者
func createNfsService(param StorageParam) {
	service := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[string]interface{}{
			"labels": map[string]interface{}{
				"app": "nfs-provisioner",
			},
			"name": "nfs-provisioner",
		},
		"spec": map[string]interface{}{
			"selector": map[string]interface{}{
				"app": "nfs-provisioner",
			},

			"ports": []map[string]interface{}{
				map[string]interface{}{
					"name": "nfs",
					"port": 2049,
				}, map[string]interface{}{
					"name": "mountd",
					"port": 20048,
				}, map[string]interface{}{
					"name": "rpcbind",
					"port": 111,
				}, map[string]interface{}{
					"name":     "rpcbind-udp",
					"port":     111,
					"protocol": "UDP",
				},
			},
		},
	}
	classClient, err := GetYamlClient(param.ClusterName, "", "v1", "api")
	resource := &v12.APIResource{Name: "Services", Namespaced: true}
	classObj := unstructured.Unstructured{Object: service}
	d, err := classClient.Resource(resource, param.Namespace).Create(&classObj)
	if err != nil {
		logs.Error("创建Nfs Service失败", d, err)
	}
}

// 2018-01-20 15:12
// 创建nfs服务
func createNfsStatefulSet(param StorageParam) {
	statefulset := map[string]interface{}{
		"apiVersion": "apps/v1beta1",
		"kind":       "StatefulSet",
		"metadata": map[string]interface{}{
			"name": "nfs-provisioner",
		},
		"spec": map[string]interface{}{
			"replicas":           2,
			"serviceAccountName": "default",
			"serviceName":        "nfs-provisioner",
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						"pod.alpha.kubernetes.io/initialized": "true",
					},
					"labels": map[string]interface{}{
						"app": "nfs-provisioner",
					},
				},
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						map[string]interface{}{
							"args": []string{"-provisioner=example.com/nfs"},
							"env": []map[string]interface{}{
								map[string]interface{}{
									"name": "POD_IP",
									"valueFrom": map[string]interface{}{
										"fieldRef": map[string]interface{}{
											"fieldPath": "status.podIP",
										},
									},
								}, map[string]interface{}{
									"name":  "SERVICE_NAME",
									"value": "nfs-provisioner",
								}, map[string]interface{}{
									"name": "POD_NAMESPACE",
									"valueFrom": map[string]interface{}{
										"fieldRef": map[string]interface{}{
											"fieldPath": "metadata.namespace",
										},
									},
								},
							},
							"image":           "quay.io/kubernetes_incubator/nfs-provisioner:v1.0.8",
							"imagePullPolicy": "IfNotPresent",
							"name":            "nfs-provisioner",
							"ports": []map[string]interface{}{
								map[string]interface{}{
									"containerPort": 2049,
									"name":          "nfs",
								}, map[string]interface{}{
									"containerPort": 20048,
									"name":          "mountd",
								}, map[string]interface{}{
									"containerPort": 111,
									"name":          "rpcbind",
								}, map[string]interface{}{
									"containerPort": 111,
									"name":          "rpcbind-udp",
									"protocol":      "UDP",
								},
							},
							"securityContext": map[string]interface{}{
								"capabilities": map[string]interface{}{
									"add": []string{"DAC_READ_SEARCH", "SYS_RESOURCE"},
								},
							},
							"volumeMounts": []map[string]interface{}{
								map[string]interface{}{
									"mountPath": "/export",
									"name":      "export-volume",
								},
								map[string]interface{}{
									"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount/",
									"name":      "service-account",
								},
							},
						},
					},
					"terminationGracePeriodSeconds": 0,
					"volumes": []map[string]interface{}{
						map[string]interface{}{
							"hostPath": map[string]interface{}{
								"path": param.HostPath,
							},
							"name": "export-volume",
						},
						map[string]interface{}{
							"name": "service-account",
							"secret": map[string]interface{}{
								"defaultMode": 420,
								"secretName":  GetServiceAccount("default", param),
							},
						},
					},
				},
			},
		},
	}
	classClient, err := GetYamlClient(param.ClusterName, "apps", "v1beta1", "/apis")
	resource := &v12.APIResource{Name: "Statefulsets", Namespaced: true}
	classObj := unstructured.Unstructured{Object: statefulset}
	d, err := classClient.Resource(resource, param.Namespace).Create(&classObj)
	if err != nil {
		logs.Error("创建Nfs StatefulSet失败", d, err)
	}
}

// 2018-01-29 14:38
//  创建pvc
func createNfsPvc(param StorageParam) error {
	pvc := map[string]interface{}{
		"kind":       "PersistentVolumeClaim",
		"apiVersion": "v1",
		"metadata": map[string]interface{}{
			"name": param.PvcName,
			"annotations": map[string]interface{}{
				"volume.beta.kubernetes.io/storage-class": param.PvcName,
			},
		},
		"spec": map[string]interface{}{
			"accessModes": []string{param.AccessMode},
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"storage": param.Size + "Mi",
				},
			},
		},
	}
	return createPvc(param.ClusterName, pvc, param.Namespace)
}

// 创建StorageClass
// 2018-01-29 15:09
func createNfsStorageClass(param StorageParam) error {

	class := map[string]interface{}{
		"kind":       "StorageClass",
		"apiVersion": "storage.k8s.io/v1beta1",
		"metadata": map[string]interface{}{
			"name": param.PvcName,
		},
		"provisioner": "example.com/nfs",
	}
	err := createStorageClass(param.ClusterName, class)
	return err
}


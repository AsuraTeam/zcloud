package k8s

import (
	"k8s.io/client-go/kubernetes"
	"cloud/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/astaxie/beego/logs"
	"time"
)

var (
	glusterNamespace = util.Namespace("glusterfs", "glusterfs")
	glusterStorageClass = "zcloud-gluster-vol"
)

// 2018-02-21 13:50
// 创建heketi的serviceaccout
func createHeketiServiceAcction(client kubernetes.Clientset) {
	name := "heketi-service-account"
	CreateServiceAccount(client, glusterNamespace, name)
}

// 2018-02-21 14:00
// 创建heketi服务
func createHeketiService(clusterName string) {
	client, _ := GetYamlClient(clusterName, "", "v1", "api")
	name := "deploy-heketi"
	ports := []map[string]interface{}{
		map[string]interface{}{
			"name":       name,
			"port":       48080,
			"targetPort": 8080,
			"nodePort":   48080,
		},
	}
	labels := map[string]interface{}{
		"glusterfs":     "heketi-service",
		"deploy-heketi": "support",
	}
	param := ServiceParam{}
	param.Labels = labels
	param.ServiceName = name
	param.Namespace = glusterNamespace
	param.Cl2 = client
	createService(ports, param)
}

// 2018-02-21 14:13
// 创建rbac
func createHeketiRbac(clustername string) {
	conf := map[string]interface{}{
		"apiVersion": "rbac.authorization.k8s.io/v1",
		"kind":       "ClusterRoleBinding",
		"metadata": map[string]interface{}{
			"name": "heketi-gluster-admin",
		},
		"roleRef": map[string]interface{}{
			"apiGroup": "rbac.authorization.k8s.io",
			"kind":     "ClusterRole",
			"name":     "edit",
		},
		"subjects": []map[string]interface{}{
			map[string]interface{}{
				"kind": "ServiceAccount",
				"name": "heketi-service-account",
			},
		},
	}
	client, _ := GetYamlClient(clustername, "rbac.authorization.k8s.io", "v1", "/apis")
	resource := &metav1.APIResource{Name: "Clusterrolebindings", Namespaced: true}
	obj := unstructured.Unstructured{Object: conf}
	_, err := client.Resource(resource, glusterNamespace).Create(&obj)
	if err != nil {
		logs.Error("创建Glusterfs Heketi rbac失败", err)
	}
	client.Resource(resource, glusterNamespace).Create(&obj)
}

// 2018-02-21 14:40
// 创建私密文件
func createHeketiSecret(clustername string) {
	conf := map[string]interface{}{
		"apiVersion": "v1",
		"data": map[string]interface{}{
			"heketi.json": "ewogICJfcG9ydF9jb21tZW50IjogIkhla2V0aSBTZXJ2ZXIgUG9ydCBOdW1iZXIiLAogICJwb3J0IjogIjgwODAiLAoKICAiX3VzZV9hdXRoIjogIkVuYWJsZSBKV1QgYXV0aG9yaXphdGlvbi4gUGxlYXNlIGVuYWJsZSBmb3IgZGVwbG95bWVudCIsCiAgInVzZV9hdXRoIjogZmFsc2UsCgogICJfand0IjogIlByaXZhdGUga2V5cyBmb3IgYWNjZXNzIiwKICAiand0IjogewogICAgIl9hZG1pbiI6ICJBZG1pbiBoYXMgYWNjZXNzIHRvIGFsbCBBUElzIiwKICAgICJhZG1pbiI6IHsKICAgICAgImtleSI6ICJNeSBTZWNyZXQiCiAgICB9LAogICAgIl91c2VyIjogIlVzZXIgb25seSBoYXMgYWNjZXNzIHRvIC92b2x1bWVzIGVuZHBvaW50IiwKICAgICJ1c2VyIjogewogICAgICAia2V5IjogIk15IFNlY3JldCIKICAgIH0KICB9LAoKICAiX2dsdXN0ZXJmc19jb21tZW50IjogIkdsdXN0ZXJGUyBDb25maWd1cmF0aW9uIiwKICAiZ2x1c3RlcmZzIjogewogICAgIl9leGVjdXRvcl9jb21tZW50IjogIkV4ZWN1dGUgcGx1Z2luLiBQb3NzaWJsZSBjaG9pY2VzOiBtb2NrLCBrdWJlcm5ldGVzLCBzc2giLAogICAgImV4ZWN1dG9yIjogImt1YmVybmV0ZXMiLAoKICAgICJfZGJfY29tbWVudCI6ICJEYXRhYmFzZSBmaWxlIG5hbWUiLAogICAgImRiIjogIi92YXIvbGliL2hla2V0aS9oZWtldGkuZGIiLAoKICAgICJrdWJlZXhlYyI6IHsKICAgICAgInJlYmFsYW5jZV9vbl9leHBhbnNpb24iOiB0cnVlLAogICAgICAibmFtZXNwYWNlIiA6ICJnbHVzdGVyZnMtLWdsdXN0ZXJmcyIsCiAgICAgICJ0b2tlbiIgOiAiYWRtaW4sMTIzNDU2IgogICAgfSwKCiAgICAic3NoZXhlYyI6IHsKICAgICAgInJlYmFsYW5jZV9vbl9leHBhbnNpb24iOiB0cnVlLAogICAgICAia2V5ZmlsZSI6ICIvZXRjL2hla2V0aS9wcml2YXRlX2tleSIsCiAgICAgICJmc3RhYiI6ICIvZXRjL2ZzdGFiIiwKICAgICAgInBvcnQiOiAiMjIiLAogICAgICAidXNlciI6ICJyb290IiwKICAgICAgInN1ZG8iOiBmYWxzZQogICAgfQogIH0sCgogICJfYmFja3VwX2RiX3RvX2t1YmVfc2VjcmV0IjogIkJhY2t1cCB0aGUgaGVrZXRpIGRhdGFiYXNlIHRvIGEgS3ViZXJuZXRlcyBzZWNyZXQgd2hlbiBydW5uaW5nIGluIEt1YmVybmV0ZXMuIERlZmF1bHQgaXMgb2ZmLiIsCiAgImJhY2t1cF9kYl90b19rdWJlX3NlY3JldCI6IGZhbHNlCn0K",
		},
		"kind": "Secret",
		"metadata": map[string]interface{}{
			"name": "heketi-config-secret",
		},
		"type": "Opaque",
	}
	client, _ := GetYamlClient(clustername, "", "v1", "/api")
	resource := &metav1.APIResource{Name: "Secrets", Namespaced: true}
	obj := unstructured.Unstructured{Object: conf}
	_, err := client.Resource(resource, glusterNamespace).Create(&obj)
	if err != nil {
		logs.Error("创建Glusterfs Heketi secrets失败", err)
	}
	client.Resource(resource, glusterNamespace).Update(&obj)
}

// 2018-02-21 14:30
// 创建deployment
func createHeketiDeployment(clustername string) error {
	conf := map[string]interface{}{
		"kind":       "Deployment",
		"apiVersion": "extensions/v1beta1",
		"metadata": map[string]interface{}{
			"name": "deploy-heketi",
			"labels": map[string]interface{}{
				"glusterfs":     "heketi-deployment",
				"deploy-heketi": "deployment",
			},
			"annotations": map[string]interface{}{
				"description": "Defines how to deploy Heketi",
			},
		},
		"spec": map[string]interface{}{
			"replicas": 1,
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"name": "deploy-heketi",
					"labels": map[string]interface{}{
						"name":          "deploy-heketi",
						"glusterfs":     "heketi-pod",
						"deploy-heketi": "pod",
					},
				},
				"spec": map[string]interface{}{
					"serviceAccountName": "heketi-service-account",
					"containers": []map[string]interface{}{
						map[string]interface{}{
							"image":           "heketi/heketi:4",
							"imagePullPolicy": "Always",
							"name":            "deploy-heketi",
							"readinessProbe": map[string]interface{}{
								"timeoutSeconds":      3,
								"initialDelaySeconds": 3,
								"httpGet": map[string]interface{}{
									"path": "/hello",
									"port": 8080,
								},
							},
							"livenessProbe": map[string]interface{}{
								"timeoutSeconds":      3,
								"initialDelaySeconds": 30,
								"httpGet": map[string]interface{}{
									"path": "/hello",
									"port": 8080,
								},
							},
							"env": []map[string]interface{}{
								map[string]interface{}{
									"name":  "HEKETI_EXECUTOR",
									"value": "kubernetes",
								},
								{
									"name":  "HEKETI_FSTAB",
									"value": "/var/lib/heketi/fstab",
								},
								{
									"name":  "HEKETI_SNAPSHOT_LIMIT",
									"value": "14",
								},
								{
									"name":  "HEKETI_KUBE_GLUSTER_DAEMONSET",
									"value": "y",
								},
							},
							"ports": []map[string]interface{}{
								map[string]interface{}{
									"containerPort": 8080,
								},
							},
							"volumeMounts": []map[string]interface{}{
								map[string]interface{}{
									"name":      "db",
									"mountPath": "/var/lib/heketi",
								},
								{
									"name":      "config",
									"mountPath": "/etc/heketi",
								},
							},
						},
					},

					"volumes": []map[string]interface{}{
						map[string]interface{}{
							"name": "db",
						},
						{
							"name": "config",
							"secret": map[string]interface{}{
								"secretName": "heketi-config-secret",
							},
						},
					},
				},
			},
		},
	}
	client, _ := GetYamlClient(clustername, "extensions", "v1beta1", "/apis")
	resource := &metav1.APIResource{Name: "Deployments", Namespaced: true}
	obj := unstructured.Unstructured{Object: conf}
	_, err := client.Resource(resource, glusterNamespace).Create(&obj)
	if err != nil {
		logs.Error("创建Glusterfs Heketi deployment失败", err)
		return err
	}
	client.Resource(resource, glusterNamespace).Update(&obj)
	return err
}

// 2018-02-21 14:51
// 创建glusterfs部署
func createGlusterDeployment(clustername string) error{
	conf := map[string]interface{}{
		"kind":       "DaemonSet",
		"apiVersion": "extensions/v1beta1",
		"metadata": map[string]interface{}{
			"name": "glusterfs",
			"labels": map[string]interface{}{
				"glusterfs": "deployment",
			},
			"annotations": map[string]interface{}{
				"description": "GlusterFS Daemon Set",
				"tags":        "glusterfs",
			},
		},
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"name": "glusterfs",
					"labels": map[string]interface{}{
						"glusterfs-node": "daemonset",
					},
				},
				"spec": map[string]interface{}{
					"nodeSelector": map[string]interface{}{
						"storagenode": "glusterfs",
					},
					"hostNetwork": true,

					"containers": []map[string]interface{}{
						map[string]interface{}{
							"image":           "gluster/gluster-centos:latest",
							"imagePullPolicy": "IfNotPresent", // 只有不存在才下载
							"name":            "glusterfs",

							"readinessProbe": map[string]interface{}{
								"timeoutSeconds":      3,
								"initialDelaySeconds": 60,
								"exec": map[string]interface{}{
									"command": []string{
										"/bin/bash",
										"-c",
										"systemctl status glusterd.service",
									},
								},
							},
							"livenessProbe": map[string]interface{}{
								"timeoutSeconds":      3,
								"initialDelaySeconds": 60,
								"exec": map[string]interface{}{
									"command": []string{
										"/bin/bash",
										"-c",
										"systemctl status glusterd.service",
									},
								},
							},
							"securityContext": map[string]interface{}{
								"capabilities": map[string]interface{}{},
								"privileged":   true,
							},
							"volumeMounts": []map[string]interface{}{
								map[string]interface{}{
									"name":      "glusterfs-heketi",
									"mountPath": "/var/lib/heketi",
								},
								{
									"name":      "glusterfs-run",
									"mountPath": "/run",
								},
								{
									"name":      "glusterfs-lvm",
									"mountPath": "/run/lvm",
								},
								{
									"name":      "glusterfs-etc",
									"mountPath": "/etc/glusterfs",
								},
								{
									"name":      "glusterfs-logs",
									"mountPath": "/var/log/glusterfs",
								},
								{
									"name":      "glusterfs-config",
									"mountPath": "/var/lib/glusterd",
								},
								{
									"name":      "glusterfs-dev",
									"mountPath": "/dev",
								},
								{
									"name":      "glusterfs-cgroup",
									"mountPath": "/sys/fs/cgroup",
								},
								{
									"name":      "localtime",
									"mountPath": "/etc/localtime",
								},
							},

						},
					},
					"volumes": []map[string]interface{}{
						map[string]interface{}{
							"name": "glusterfs-heketi",
							"hostPath": map[string]interface{}{
								"path": "/var/lib/heketi",
							},
						},
						{
							"name": "glusterfs-run",
						},
						{
							"name": "glusterfs-lvm",
							"hostPath": map[string]interface{}{
								"path": "/run/lvm",
							},
						},
						{
							"name": "glusterfs-etc",
							"hostPath": map[string]interface{}{
								"path": "/etc/glusterfs",
							},
						},
						{
							"name": "glusterfs-logs",
							"hostPath": map[string]interface{}{
								"path": "/var/log/glusterfs",
							},
						},
						{
							"name": "glusterfs-config",
							"hostPath": map[string]interface{}{
								"path": "/var/lib/glusterd",
							},
						},
						{
							"name": "glusterfs-dev",
							"hostPath": map[string]interface{}{
								"path": "/dev",
							},
						},
						{
							"name": "glusterfs-cgroup",
							"hostPath": map[string]interface{}{
								"path": "/sys/fs/cgroup",
							},
						},
						{
							"name": "localtime",
							"hostPath": map[string]interface{}{
								"path": "/etc/localtime",
								"defaultMode":400,
							},
						},
					},
				},
			},
		},
	}
	client, _ := GetYamlClient(clustername, "extensions", "v1beta1", "/apis")
	resource := &metav1.APIResource{Name: "Daemonsets", Namespaced: true}
	obj := unstructured.Unstructured{Object: conf}
	_, err := client.Resource(resource, glusterNamespace).Create(&obj)
	if err != nil {
		logs.Error("创建Glusterfs Daemonsets失败", err)
	}
	return err
}

// 2018-02-21 16:12
// 更新glusterfs集群信息
func UpdateGlusterfsTopology(clustername string, client kubernetes.Clientset) {
	nodes := GetNodes(client, "storagenode=glusterfs")
	resouces := make([]map[string]interface{}, 0)
	for _, v := range nodes {
		node := map[string]interface{}{
			"node": map[string]interface{}{
				"hostnames": map[string]interface{}{
					"manage": []string{
						v.Name,
					},
					"storage": []string{
						v.Name,
					},
				},
				"zone": 1,
			},
			"devices": []string{
				"/dev/vdb",
			},
		}
		resouces = append(resouces, node)
	}
	cluster := map[string]interface{}{
		"clusters": []map[string]interface{}{
			map[string]interface{}{
				"nodes": resouces,
			},
		},
	}
	logs.Info(util.ObjToString(cluster))
}

/*
(4) distribute stripe volume 分布式条带卷
Brick server 数量是条带数的倍数，兼具 distribute 和 stripe 卷的特点。
分布式的条带卷，volume 中 brick 所包含的存储服务器数必须是 stripe 的倍数(>=2倍)，
兼顾分布式和条带式的功能。每个文件分布在四台共享服务器上，
通常用于大文件访问处理，最少需要 4 台服务器才能创建分布条带卷。
 */
// 2018-02-22 10:07
// 创建glusterfs提供者
func createGlusterfsStorageClass(clustername string) {
	class := map[string]interface{}{
		"apiVersion": "storage.k8s.io/v1beta1",
		"kind":       "StorageClass",
		"metadata": map[string]interface{}{
			"name": glusterStorageClass,
		},
		"parameters": map[string]interface{}{
			"resturl":         "http://127.0.0.1:48080",
			"restuser":        "", // 可选，authentication 的用户名
			"secretName":      "",  // 可选，authentication 的密码所在的 secret
			"secretNamespace": "", // 可选，authentication 的密码所在的 secret 所在的namespace
		},
		"provisioner":   "kubernetes.io/glusterfs",
		"reclaimPolicy": "Delete",
	}
	createStorageClass(clustername, class)
}

// 2018-02-21 15:03
// 需要物理硬件支持,每个机器有一块单独的硬盘
// 创建glusterfs集群
// 节点需要添加标签 storagenode=glusterfs
func CreateGlusterfs(param StorageParam) {
	client, _ := GetClient(param.ClusterName)
	CreateServiceAccount(client, glusterNamespace, "default")
	createGlusterDeployment(param.ClusterName)
	time.Sleep(time.Minute * 1)
	createHeketiServiceAcction(client)
	createHeketiRbac(param.ClusterName)
	createHeketiSecret(param.ClusterName)
	createHeketiServiceAcction(client)
	createHeketiDeployment(param.ClusterName)
	createHeketiService(param.ClusterName)
	createGlusterfsStorageClass(param.ClusterName)
}


// 2018-02-22 10:38
// 创建glusterfs pvc
func createGlusterfsPvc(param StorageParam) {
	pvc := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "PersistentVolumeClaim",
		"metadata": map[string]interface{}{
			"annotations": map[string]interface{}{
				"volume.beta.kubernetes.io/storage-class":   glusterStorageClass,
				"volume.beta.kubernetes.io/storage-provisioner": "kubernetes.io/glusterfs",
			},
			// glusterfs-claim
			"name":      param.PvcName,
			"namespace": param.Namespace,
		},
		"spec": map[string]interface{}{
			"accessModes": []string{
				param.AccessMode,
			},
			"resources": map[string]interface{}{
				"requests": map[string]interface{}{
					"storage": param.Size + "Mi",
				},
			},
		},
	}
	createPvc(param.ClusterName, pvc, param.Namespace)
}
package k8s

import (
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/astaxie/beego/logs"
	"strings"
	v1beta12 "k8s.io/api/apps/v1beta1"
	"fmt"
	"cloud/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strconv"
	"encoding/json"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// 获某个namespace下面的deployment信息
// {"metadata":{"name":"deploymentexample","namespace":"testservice--asdfasdfdasf","selfLink":"/apis/apps/v1beta1/namespaces/testservice--asdfasdfdasf/deployments/deploymentexample","uid":"c30cd8ea-f34e-11e7-8d1c-0894ef37b2d2","resourceVersion":"3433025","generation":1,"creationTimestamp":"2018-01-07T02:02:06Z","labels":{"release-version":"0","space":"testservice--asdfasdfdasf","uuid":"7e22a9b7fc32748be8527f6e2592ea67","zcloud-app":"testservice--asdfasdfdasf"},"annotations":{"deployment.kubernetes.io/revision":"1"}},"spec":{"replicas":3,"selector":{"matchLabels":{"app":"nginx"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"nginx"}},"spec":{"containers":[{"name":"nginx","image":"nginx:1.10","ports":[{"containerPort":80,"protocol":"TCP"}],"resources":{},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","securityContext":{},"schedulerName":"default-scheduler"}},"strategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":"25%","maxSurge":"25%"}},"revisionHistoryLimit":2,"progressDeadlineSeconds":600},"status":{"observedGeneration":1,"replicas":3,"updatedReplicas":3,"readyReplicas":3,"availableReplicas":3,"conditions":[{"type":"Available","status":"True","lastUpdateTime":"2018-01-07T02:02:09Z","lastTransitionTime":"2018-01-07T02:02:09Z","reason":"MinimumReplicasAvailable","message":"Deployment has minimum availability."},{"type":"Progressing","status":"True","lastUpdateTime":"2018-01-07T02:02:09Z","lastTransitionTime":"2018-01-07T02:02:06Z","reason":"NewReplicaSetAvailable","message":"ReplicaSet \"deploymentexample-845cfc7fb9\" has successfully progressed."}]}}
func GetDeployments(namespace string, clientset kubernetes.Clientset) []v1beta12.Deployment {
	opt := metav1.ListOptions{}
	deployments, err := clientset.AppsV1beta1().Deployments(namespace).List(opt)
	if err != nil {
		logs.Error("获取deployment 错误", err)
		return make([]v1beta12.Deployment, 0)
	}
	return deployments.Items
}

// 2018-02-19 15:08
// 获取Deployment信息
func GetDeployment(namespace string, clientset kubernetes.Clientset, name string) v1beta12.Deployment {
	deployment, err := clientset.AppsV1beta1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err == nil {
		return *deployment
	}
	logs.Error("获取deployment失败", err)
	return v1beta12.Deployment{}
}

// 2018-02-04 17:36
// 获取deploy更新前版本
func GetDeploymentsVersion(namespace string, name string, client kubernetes.Clientset) string {
	deploy, err := client.AppsV1beta1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return deploy.Annotations["deployment.kubernetes.io/revision"]
	}
	return ""
}


// 2018-02-04 19:55
// 更新镜像
func UpdateDeploymentImage(param RollingParam) (bool, error) {
	deploy, err := param.Client.AppsV1beta1().Deployments(param.Namespace).Get(param.Name, metav1.GetOptions{})
	if err != nil {
		logs.Error("UpdateDeploymentImage - 1 ", err)
		return false, err
	}

	secretsName := GetDockerImagePullName(strings.Split(param.Images, "/")[0])
	isExists := SecretIsExists(param.Client, param.Namespace, secretsName)
	if isExists {
		secrets := v1.LocalObjectReference{}
		secrets.Name = secretsName
		deploy.Spec.Template.Spec.ImagePullSecrets = []v1.LocalObjectReference{secrets}
		logs.Info("获取到安全密码", secrets)
	}
	deploy.Spec.Template.Spec.Containers[0].Image = param.Images

	// 更新参数配置
	if param.MinReadySeconds > 0 {
		deploy.Spec.MinReadySeconds = param.MinReadySeconds
	}
	if param.MaxSurge > 0 {
		deploy.Spec.Strategy.RollingUpdate.MaxSurge = &intstr.IntOrString{IntVal: param.MaxSurge}
	}
	logs.Error("UpdateDeploymentImage - 11")
	var m int64
	m = param.TerminationGracePeriodSeconds
	if m > 0 {
		deploy.Spec.Template.Spec.TerminationGracePeriodSeconds = &m
	}
	if param.MaxUnavailable > 0 {
		deploy.Spec.Strategy.RollingUpdate.MaxUnavailable = &intstr.IntOrString{IntVal: param.MaxUnavailable}
	}

	r, err := param.Client.AppsV1beta1().Deployments(param.Namespace).Update(deploy)
	if err != nil {
		logs.Error("UpdateDeploymentImage - 2 ", err)
		return false, err
	}
	logs.Error("UpdateDeploymentImage - 结果 ", util.ObjToString(r))
	return true, nil
}

// 2018-02-04 18:14
// 获取是否可以更新deployment
func GetDeploymentStatus(namespace string, name string, client kubernetes.Clientset) (bool, string) {
	deploy, err := client.AppsV1beta1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return true, ""
	}
	for _, v := range deploy.Status.Conditions {
		if v.Type == "Progressing" {
			if strings.Contains(v.Message, "has successfully progressed.") {
				return true, v.Message
			} else {
				return false, v.Message
			}
		}
	}
	return false, ""
}

// 获取某个服务的信息
// 2018-01-18 10:02
func GetDeploymentsService(namespace string, clientset kubernetes.Clientset, service string) []v1beta12.Deployment {
	opt := metav1.ListOptions{}
	opt.LabelSelector = "name=" + service
	deployments, err := clientset.AppsV1beta1().Deployments(namespace).List(opt)
	if err != nil {
		logs.Error("获取deployment 错误", err)
		return make([]v1beta12.Deployment, 0)
	}
	return deployments.Items
}

// 获取某个uuid标签的下面的pod服务
// 在创建完pod后自动创建service使用
func GetPodsFromUUid(namespace string, uuid string, clientSet kubernetes.Clientset) []ServicePod {
	pmap := make([]ServicePod, 0)
	deployments := GetDeployments(namespace, clientSet)
	logs.Info("deployments", deployments)
	name := strings.Split(namespace, "--")
	resouceName := name[1]
	appName := name[0]
	logs.Info("name", name)
	for _, deployment := range deployments {
		labels := deployment.Labels
		puuid := labels["uuid"]
		logs.Info("puuid", puuid, uuid)
		if puuid == uuid {
			temp := ServicePod{}
			temp.ResouceName = resouceName
			temp.AppName = appName
			temp.Selector = deployment.Spec.Selector.MatchLabels
			for _, c := range deployment.Spec.Template.Spec.Containers {
				temp.ClusterName = c.Name
				for _, port := range c.Ports {
					temp.ContainerPort = append(temp.ContainerPort, port.ContainerPort)
					pmap = append(pmap, temp)
				}
			}
			break
		}
	}
	return pmap
}

// 删除 deployment
func DeletelDeployment(namespace string, isService bool, name string, clusterName string) error {
	cl, err := GetYamlClient(clusterName, "apps", "v1beta1", "/apis")
	resource := &metav1.APIResource{Name: "Deployments", Namespaced: true}
	opt := metav1.DeleteOptions{}
	client, err := GetClient(clusterName)
	deploments := GetDeployments(namespace, client)
	if len(deploments) < 0 {
		logs.Error("删除yaml的deployment服务失败,没有找到对应的记录")
		return err
	}

	// 删除指定名称的服务
	if name != "" {
		logs.Info("删除deploy", namespace, name)
		err = cl.Resource(resource, namespace).Delete(name, &opt)
		if err != nil {
			logs.Error("删除yaml的deployment服务失败", err)
			//return err
		}

	} else {
		for _, deployment := range deploments {
			logs.Info("删除 deployment ", deployment.Namespace, deployment.Name)
			err = cl.Resource(resource, namespace).Delete(deployment.Name, &opt)
			if err != nil {
				logs.Error("删除yaml的deployment服务失败", err)
				//return err
			}
		}
	}

	obj := metav1.ListOptions{}
	rcl, _ := GetClient(clusterName)
	replications, err := rcl.ExtensionsV1beta1().ReplicaSets(namespace).List(obj)

	// 删除 Replicationcontrollers
	if name != "" {
		obj.LabelSelector = "name=" + name
		replications, err := rcl.ExtensionsV1beta1().ReplicaSets(namespace).List(obj)
		if err == nil {
			logs.Info("获取到 ReplicaSets 大小为", len(replications.Items))
			if len(replications.Items) > 0 {
				err := rcl.ExtensionsV1beta1().ReplicaSets(namespace).Delete(replications.Items[0].Name, &metav1.DeleteOptions{})
				if err != nil {
					logs.Error("删除 ReplicaSets 失败", err)
					return err
				} else {
					logs.Info("删除 ReplicaSets ", replications.Items[0].Name)
				}
			}
		} else {
			logs.Error("获取 ReplicaSets 失败", err)
		}

	} else {
		logs.Info("删除 Replicaset ")
		for _, v := range replications.Items {
			err := rcl.ExtensionsV1beta1().ReplicaSets(namespace).Delete(v.Name, &metav1.DeleteOptions{})
			if err != nil {
				logs.Error("获取 ReplicaSets 失败", err, v.Namespace, v.Name)
			}
		}
	}

	// 删除pod
	pods, err := rcl.CoreV1().Pods(namespace).List(obj)
	for _, v := range pods.Items {
		rcl.CoreV1().Pods(namespace).Delete(v.Name, &metav1.DeleteOptions{})
	}

	// 删除服务
	if isService {
		logs.Info("开始删除Service Service ... ", namespace)
		if name != "" {
			err = DeleteService(clusterName, namespace, name)
		} else {
			services, err := GetServices(rcl, namespace)
			if err == nil {
				for _, v := range services {
					err = DeleteService(clusterName, namespace, v.Name)
					if err == nil {
						logs.Info("删除服务成功", namespace, v.Name)
					} else {
						logs.Error("删除服务失败", namespace, v.Name)
					}
				}
			}
		}
	}

	return err
}

// svc {"metadata":{"name":"auto-nginx-3","namespace":"auto-nginx-3--dfsad","selfLink":"/api/v1/namespaces/auto-nginx-3--dfsad/services/auto-nginx-3","uid":"2c62631d-f773-11e7-8d1c-0894ef37b2d2","resourceVersion":"4030027","creationTimestamp":"2018-01-12T08:32:49Z","labels":{"app":"auto-nginx-3"}},"spec":{"ports":[{"name":"auto-nginx-3-0","protocol":"TCP","port":49873,"targetPort":80,"nodePort":49873}],"selector":{"name":"auto-nginx-3"},"clusterIP":"172.16.1.62","type":"NodePort","sessionAffinity":"None","externalTrafficPolicy":"Cluster"},"status":{"loadBalancer":{}}}
// deploy {"metadata":{"name":"auto-3","namespace":"auto-3--dfsad","selfLink":"/apis/apps/v1beta1/namespaces/auto-3--dfsad/deployments/auto-3","uid":"ee1a2658-f780-11e7-8d1c-0894ef37b2d2","resourceVersion":"4037873","generation":1,"creationTimestamp":"2018-01-12T10:11:18Z","labels":{"name":"auto-3"},"annotations":{"deployment.kubernetes.io/revision":"1"}},"spec":{"replicas":1,"selector":{"matchLabels":{"name":"auto-3"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"name":"auto-3"}},"spec":{"containers":[{"name":"auto-3","image":"nginx:1.10","ports":[{"containerPort":80,"protocol":"TCP"}],"resources":{"limits":{"cpu":"1","memory":"2Gi"},"requests":{"cpu":"1","memory":"2Gi"}},"terminationMessagePath":"/dev/termination-log","terminationMessagePolicy":"File","imagePullPolicy":"IfNotPresent"}],"restartPolicy":"Always","terminationGracePeriodSeconds":30,"dnsPolicy":"ClusterFirst","securityContext":{},"schedulerName":"default-scheduler"}},"strategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":"25%","maxSurge":"25%"}},"revisionHistoryLimit":2,"progressDeadlineSeconds":600},"status":{"observedGeneration":1,"replicas":1,"updatedReplicas":1,"readyReplicas":1,"availableReplicas":1,"conditions":[{"type":"Available","status":"True","lastUpdateTime":"2018-01-12T10:11:20Z","lastTransitionTime":"2018-01-12T10:11:20Z","reason":"MinimumReplicasAvailable","message":"Deployment has minimum availability."},{"type":"Progressing","status":"True","lastUpdateTime":"2018-01-12T10:11:20Z","lastTransitionTime":"2018-01-12T10:11:18Z","reason":"NewReplicaSetAvailable","message":"ReplicaSet \"auto-3-8548fd9d57\" has successfully progressed."}]}}
// 获取自己创建的namespace应用, 规则是app名加资源名区分
func GetDeploymentApp(clientSet kubernetes.Clientset, namespace string, service string) map[string]CloudApp {
	result := map[string]CloudApp{}

	deployments := make([]v1beta12.Deployment, 0)
	if service != "" {
		deployments = GetDeploymentsService(namespace, clientSet, service)
	} else {
		deployments = GetDeployments(namespace, clientSet)
	}
	//datas := []CloudApp{}
	for _, v := range deployments {
		data := CloudApp{}
		data.ServiceName = service
		data.ContainerNumber = int(v.Status.Replicas)
		data.AvailableReplicas = v.Status.AvailableReplicas
		data.CreateTime = util.GetMinTime(util.ReplaceTime(v.CreationTimestamp.String()))
		data.Image = v.Spec.Template.Spec.Containers[0].Image
		name := strings.Split(v.Namespace, "--")
		if len(name) < 2 {
			continue
		}

		access := make([]string, 0)
		if service != "" {
			svc := GetAppService(clientSet, v.Namespace, service)
			for _, svcport := range svc.Spec.Ports {
				a := svc.Name + "." + v.Namespace + ":" + strconv.Itoa(int(svcport.NodePort))
				access = append(access, a)
			}
			data.ServiceNumber = 1
		} else {
			svcs, _ := GetServices(clientSet, v.Namespace)
			if len(svcs) > 0 {
				for _, svc := range svcs[0].Spec.Ports {
					a := svcs[0].Name + "." + v.Namespace + ":" + strconv.Itoa(int(svc.NodePort))
					access = append(access, a)
				}
			}
			data.ServiceNumber = len(svcs)
		}
		data.Access = access

		//data.LastUpdateTime = v.Status.
		pods := GetPods(v.Namespace, clientSet)
		if service == "" {
			data.ContainerNumber = len(pods)
		}
		//data.ContainerNumber = len(pods)
		if len(pods) > 0 {
			if pods[0].Status.StartTime != nil {
				data.LastUpdateTime = util.ReplaceTime(pods[0].Status.StartTime.String())
			}
		} else {
			data.LastUpdateTime = "未知"
		}
		conditions := v.Status.Conditions
		if len(conditions) > 0 {
			data.Status = string(conditions[0].Status)
		} else {
			data.Status = "未知"
		}
		data.AppLabels = "无"
		data.AppName = name[0]
		data.ClusterName = name[1]
		result[v.Namespace+service] = data
	}

	return result
}

// 加工容器端口数据
// 获取容器端口数据
// 2018-01-11 13:47
func getPorts(port string, hostport string) []map[string]interface{} {
	ports := make([]map[string]interface{}, 0)
	hostsports := strings.Split(hostport, ",")
	for idx, p := range strings.Split(port, ",") {
		pv, err := strconv.Atoi(p)
		if err != nil {
			continue
		}
		data := map[string]interface{}{
			"containerPort": pv,
			"protocol":      "TCP",
		}
		if len(hostsports) > idx {
			hostport, err := strconv.Atoi(hostsports[idx])
			if err == nil {
				data["hostPort"] = hostport
			}
		}
		ports = append(ports, data)
	}

	return ports
}

// 获取服务和pod的映射端口数据
// 2018-01-12 15:36
func getServicePorts(param ServiceParam) []map[string]interface{} {
	var all []int
	if param.Namespace == util.Namespace("registryv2", "registryv2") {
		// 仓库单独使用的地址段,不分配给其他服务
		logs.Info("获取到要给仓库分配端口", param.Namespace, param.Name, param.ClusterName)
		all = GetServicePorts(param.Cl3, len(strings.Split(param.PortData, ",")), 49000, 50000)
	} else {
		all = GetServicePorts(param.Cl3, len(strings.Split(param.PortData, ",")), 0, 0)
	}
	ports := make([]map[string]interface{}, 0)
	for id, port := range strings.Split(param.PortData, ",") {
		if len(all) <= id {
			continue
		}
		free := all[id]
		oldData, ok := param.OldPort.Get(port)
		temp := make(map[string]interface{})
		//  保持不更新原有的
		if ok {
			oldp := oldData.(util.Lock)
			temp = map[string]interface{}{
				"name":       oldp.GetV("name"),
				"port":       oldp.GetV("port"),
				"targetPort": oldp.GetV("targetPort"),
				"nodePort":   oldp.GetV("nodePort"),
				"protocol":   oldp.GetV("protocol"),
			}
		} else {
			temp = map[string]interface{}{
				"name":       param.Name + "-" + strconv.Itoa(id),
				"port":       free,
				"targetPort": util.StringToInt(port),
				"nodePort":   free,
				"protocol":   "TCP",
			}
		}
		ports = append(ports, temp)
	}
	if param.IsRedeploy {
		data := make([]map[string]interface{}, 0)
		json.Unmarshal([]byte(param.PortYaml), &data)
		logs.Info("yaml 数据", data)
		if len(data) > 0 {
			for _, v := range data{
				d := v
				if d != nil {
					if d["kind"] == "Service" {
						spec := d["spec"].(map[string]interface{})
						if spec != nil {
							portData := spec["ports"].([]interface{})
							jsonData, err  := json.Marshal(portData)
							if err == nil {
								json.Unmarshal(jsonData, &ports)
							}
							logs.Info("重建服务获取到端口", util.ObjToString(ports))
						}
					}
				}
			}
		}
	}
	return ports
}

// 制作亲和性数据
// 2018-01-11 14:55
func getAffinity(affinityData string) []map[string]interface{} {
	data := make([]Affinity, 0)
	err := json.Unmarshal([]byte(affinityData), &data)
	if err != nil {
		logs.Error("创建Affinity数据失败", err)
		return make([]map[string]interface{}, 0)
	}
	affinitys := make([]map[string]interface{}, 0)
	for _, v := range data {
		d := map[string]interface{}{
			"key":      v.Type,
			"operator": "In",
			"values":   []string{v.Value},
		}
		affinitys = append(affinitys, d)
	}
	return affinitys
}

// 并且node上需要运行有app=store的label.
// 和关联的应用运行到一个机器上面去
func get() {
	v := map[string]interface{}{
		"affinity": map[string]interface{}{
			"podAffinity": map[string]interface{}{
				"requiredDuringSchedulingIgnoredDuringExecution": []map[string]interface{}{
					map[string]interface{}{
						"labelSelector": map[string]interface{}{
							"matchExpressions": []map[string]interface{}{
								map[string]interface{}{
									"key":      "app",
									"operator": "In",
									"values":   []string{"store"},
								},
							},
						},
						"topologyKey": "kubernetes.io/hostname",
					},
				},
			},
			"podAntiAffinity": map[string]interface{}{
				"requiredDuringSchedulingIgnoredDuringExecution": []map[string]interface{}{
					map[string]interface{}{
						"labelSelector": map[string]interface{}{
							"matchExpressions": []map[string]interface{}{
								map[string]interface{}{
									"key":      "app",
									"operator": "In",
									"values":   []string{"web-store"},
								},
							},
						},
						"topologyKey": "kubernetes.io/hostname",
					},
				},
			},
		},
	}
	fmt.Println(v)
}

// 选择指定的node节点,只IP地址选择
// 2018-01-11 16:10
// "kubernetes.io/hostname"
func getNodeSelectorNode(selector string) map[string]interface{} {
	d := NodeSelector{}
	err := json.Unmarshal([]byte(selector), &d)
	if err != nil {
		return make(map[string]interface{}, 0)
	}
	return map[string]interface{}{
		d.Lables: d.Value,
	}
}

// 获取健康检查数据
// cmd http tcp 3种类型的
// 2018-01-12 8:15
// periodSeconds 执行平率（秒）
// initialDelaySeconds 在容器开启之前的时间，就是等待容器启动完成后检查
func getHealthData(healthData string) (interface{}, bool) {
	check := make(map[string]interface{})
	d := HealthData{}
	err := json.Unmarshal([]byte(healthData), &d)

	if err != nil {
		logs.Warn("获取健康检查配置,配置为空或没有配在", err)
		return check, false
	}

	if d.HealthType == "" {
		logs.Warn("获取健康检查配置,配置为空或没有配在", err)
		return check, false
	}

	port, porterr := strconv.Atoi(strings.TrimSpace(d.HealthPort))
	switch strings.ToLower(d.HealthType) {
	case "cmd":
		check = map[string]interface{}{
			"livenessProbe": map[string]interface{}{
				"exec": map[string]interface{}{
					"command": strings.Split(d.HealthCmd, " "),
				},
			},
		}
	case "http":
		if porterr != nil {
			logs.Error("http检查端口异常")
			return check, false
		}
		check = map[string]interface{}{
			"livenessProbe": map[string]interface{}{
				"httpGet": map[string]interface{}{
					"path":   d.HealthPath,
					"port":   port,
					"scheme": "HTTP",
				},
			},
		}
	case "tcp":
		if porterr != nil {
			logs.Error("tcp检查端口异常")
			return check, false
		}
		check = map[string]interface{}{
			"livenessProbe": map[string]interface{}{
				"tcpSocket": map[string]interface{}{
					"port": port,
				},

			},
		}
	default:
		return check, false
		break
	}

	livenessProbe := check["livenessProbe"].(map[string]interface{})
	livenessProbe["initialDelaySeconds"] = util.StringToInt(d.HealthInitialDelay)
	livenessProbe["periodSeconds"] = util.StringToInt(d.HealthInterval)
	livenessProbe["timeoutSeconds"] = util.StringToInt(d.HealthTimeout)
	livenessProbe["failureThreshold"] = util.StringToInt(d.HealthFailureThreshold)
	return check["livenessProbe"].(map[string]interface{}), true
}

// 自动创建服务
// 2018-01-12 15:57
func createService(ports []map[string]interface{}, param ServiceParam) (interface{}, error) {

	conf := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[string]interface{}{
			"name": param.ServiceName,
			"labels": map[string]interface{}{
				"app": param.ServiceName,
			},
		},
		"spec": map[string]interface{}{
			"type": "NodePort",
			"selector": map[string]interface{}{
				"name": param.ServiceName,
			},
			"ports": ports,
		},
	}

	// 标签配置
	if len(param.Labels) > 0 {
		conf["metadata"].(map[string]interface{})["labels"] = param.Labels
	}

	// 更新的时候需要集群IP地址
	if _, ok := param.OldPort.Get("clusterIp"); ok {
		clusterIp := param.OldPort.GetV("clusterIp").(string)
		if clusterIp != "" {
			conf["spec"].(map[string]interface{})["clusterIP"] = clusterIp
			//conf["metadata"].(map[string]interface{})["resourceVersion"] = param.OldPort.GetV("resourceVersion")
		}
	}

	if param.SessionAffinity  != "" {
		conf["spec"].(map[string]interface{})["sessionAffinity"] = param.SessionAffinity
	}

	resource := &metav1.APIResource{Name: "Services", Namespaced: true}
	obj := unstructured.Unstructured{Object: conf}

	var d *unstructured.Unstructured
	var err error
	cl := param.Cl2

	delopt := metav1.DeleteOptions{}
	err = cl.Resource(resource, param.Namespace).Delete(param.ServiceName, &delopt)
	if err != nil {
		logs.Error("删除服务失败,原有服务不存在", err)
	}

	d, err = cl.Resource(resource, param.Namespace).Create(&obj)
	if err != nil {
		logs.Error("创建yaml的service失败", err)
		return "", err
	}
	logs.Info("创建yaml的service服务成功", d)
	return conf, err
}

// 2018-02-02 12:40
func setContainerCommand(param ServiceParam, v map[string]interface{}) map[string]interface{} {
	command := make([]string, 0)
	err := json.Unmarshal([]byte(param.Command), &command)
	if err == nil {
		v["spec"].(map[string]interface{})["template"].
		(map[string]interface{})["spec"].
		(map[string]interface{})["containers"].
		([]map[string]interface{})[0]["command"] = command
	}
	return v
}

// 2018-02-02 15:35
// 设置特权模式
func setPrivileged(param ServiceParam, v map[string]interface{}) map[string]interface{} {
	// 有特权的操作
	if param.Privileged {
		v["spec"].(map[string]interface{})["template"].
		(map[string]interface{})["spec"].
		(map[string]interface{})["containers"].
		([]map[string]interface{})[0]["securityContext"] = map[string]interface{}{
			"capabilities": map[string]interface{}{},
			"privileged":   true,
		}
	}
	return v
}

// 2018-02-09 21:43
// 镜像下载策略配置
func setImagePullPolice(param ServiceParam, v map[string]interface{}) map[string]interface{} {
	// 添加下载镜像secret
	if param.RegistryAuth != "" {
		secrets := []map[string]interface{}{
			map[string]interface{}{
				"Name": GetDockerImagePullName(param.Registry),
			},
		}
		v["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["ImagePullSecrets"] = secrets
	}
	return v
}

// 创建服务
//c1,_ := k8s.GetYamlClient("10.16.55.6","8080","apps","v1beta1","/apis")
//storageData := `[{"ContainerPath":"/tmp/mnt","HostPath":"/mnt","Volume":""},{"Volume":"","ContainerPath":"/tmp","HostPath":"/mnt"}]`
//lables := `{"Value":"10.16.55.102","Lables":"kubernetes.io/hostname"}`
//affinityData := `[{"Type":"zone","Value":"node103"}]`
//k8s.CreateServicePod(c1,"default","test1",0.1,"1024","80",storageData,"nginx:1.11",affinityData, lables)
// 2018-01-11 15:02
//c1,_ := k8s.GetYamlClient("10.16.55.6","8080","apps","v1beta1","/apis")
func CreateServicePod(param ServiceParam) (string, error) {
	status, result := GetDeploymentStatus(param.Namespace, param.Name, param.Cl3)
	CreateServiceAccount(param.Cl3, param.Namespace, "default")
	if ! status {
		return "服务已经在更新,不能更新", errors.NewTooManyRequestsError(result)
	}
	if len(param.ConfigureData) > 0 {
		logs.Info("需要创建或更新配置文件", param.ConfigureData)
		CreateConfigmap(param)
	}

	if param.TerminationSeconds == 0 {
		logs.Info("使用默认 TerminationSeconds 50秒")
		param.TerminationSeconds = 50
	}

	yamldata := make([]interface{}, 0)
	volumes, volumeMounts := getVolumes(param.StorageData, param.ConfigureData, param)
	name := param.ServiceName
	v := map[string]interface{}{
		"apiVersion": "apps/v1beta1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"labels": map[string]interface{}{
				"name": name,
			},
			"name": name,
		},
		"spec": map[string]interface{}{
			"replicas":        param.Replicas,
			"minReadySeconds": param.MinReady, // 滚动升级多少秒认为该pod就绪
			"strategy": map[string]interface{}{
				"rollingUpdate": map[string]interface{}{ // 假如replicas =3 ， 滚动升级pod数量到2-4个之间
					"maxSurge":       2, // 滚动升级时会先启动1个pod
					"maxUnavailable": 1, // 滚动升级时允许的最大Unavailable的pod个数
					//例如，该值设置成30%，启动rolling update后旧的ReplicatSet将会立即缩容到期望的Pod数量的70%。
					// 新的Pod ready后，随着新的ReplicaSet的扩容，旧的ReplicaSet会进一步缩容，确保在升级的所有时刻可以用的Pod数量至少是期望Pod数量的70%。
				},
			},
			"selector": map[string]interface{}{
				"matchLabels": map[string]interface{}{
					"name": name,
				},
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"name": name,
					},
				},
				"spec": map[string]interface{}{
					"nodeSelector": getNodeSelectorNode(param.Selector),
					"terminationGracePeriodSeconds": param.TerminationSeconds, // 优雅的关闭进程,默认30秒
					"containers": []map[string]interface{}{
						map[string]interface{}{
							"image": param.Image,
							"name":  name,
							"ports": getPorts(param.Port, param.HostPort),
							"resources": map[string]interface{}{
								"limits": map[string]interface{}{
									"memory": param.Memory + "Mi",
									"cpu":    param.Cpu,
								},
								"requests": map[string]interface{}{
									"memory": param.Memory + "Mi",
									"cpu":    param.Cpu,
								},
							},
							"env":          getEnv(param.Envs),
							"volumeMounts": volumeMounts,
						},
					},
					"volumes": volumes,
				},
			},
		},
	}

	// 主机网络模式
	if len(param.NetworkMode) > 0 {
		v["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["hostNetwork"] = true
	}

	// 绑定filebeat容器
	if len(param.Kafka) > 0 && len(param.LogPath) > 0 {
		spec := v["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})
		spec["containers"] = append(spec["containers"].([]map[string]interface{}), map[string]interface{}{})
		filebeatContainer := CreateFilebeatConfig(param)
	    spec["containers"].([]map[string]interface{})[1] = filebeatContainer
		filebeatVolumes, filebeatMount := getFilebeatStorage(param)
		logs.Info("filebeatMount", filebeatMount)
		if len(filebeatVolumes) > 0 {
			oldvolumes :=  spec["volumes"].([]map[string]interface{})
			oldvolumes = append(oldvolumes, filebeatVolumes...)
			oldvolumes = append(oldvolumes, getFilebeatConfig(param))
			spec["volumes"] = oldvolumes
			filebeatOldVolumeMounts := spec["containers"].([]map[string]interface{})[1]["volumeMounts"].([]map[string]interface{})
			filebeatOldVolumeMounts = append(filebeatOldVolumeMounts, filebeatMount...)
			spec["containers"].([]map[string]interface{})[1]["volumeMounts"] = filebeatOldVolumeMounts

			serviceOldVolumeMounts := spec["containers"].([]map[string]interface{})[0]["volumeMounts"].([]map[string]interface{})
			serviceOldVolumeMounts = append(serviceOldVolumeMounts, filebeatMount...)
			spec["containers"].([]map[string]interface{})[0]["volumeMounts"] = serviceOldVolumeMounts
		}
	}

	healthdata, ok := getHealthData(param.HealthData)
	if ok {
		// readinessProbe
		v["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]map[string]interface{})[0]["readinessProbe"] = make(map[string]interface{})
		v["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]map[string]interface{})[0]["readinessProbe"] = healthdata
	}
	v = setImagePullPolice(param, v)
	if param.Command != "" {
		v = setContainerCommand(param, v)
	}

	if param.Privileged {
		v = setPrivileged(param, v)
	}

	logs.Info(util.ObjToString(v))

	yamldata = append(yamldata, v)
	t, _ := json.Marshal(v)
	yaml := string(t)
	resource := &metav1.APIResource{Name: "Deployments", Namespaced: true}
	obj := unstructured.Unstructured{Object: v}

	var d *unstructured.Unstructured
	var err error
	if param.Update {
		d, err = param.C1.Resource(resource, param.Namespace).Update(&obj)
	} else {
		d, err = param.C1.Resource(resource, param.Namespace).Create(&obj)
	}

	if err != nil {
		logs.Error("创建yaml的Deployment失败", err)
		return yaml, err
	}

	logs.Info("创建yaml的 Deployment 服务成功", d)
	if param.Update && param.UpdateType != "port" {
		logs.Info("更新服务完成，退出..")
		return yaml, err
	}

	param.OldPort = GetCurrentPort(param.Cl3, param.Namespace, param.Name)
	logs.Info("创建和更新服务", "已经存在的端口", param.OldPort.String())
	ports := getServicePorts(param)

	logs.Info("新端口", "获取到新端口", ports)

	d1, err := createService(ports, param)

	if err != nil {
		logs.Error("创建服务失败", err)
		return "", err
	}
	yamldata = append(yamldata, d1)
	conf, err := json.Marshal(yamldata)
	return string(conf), err
}

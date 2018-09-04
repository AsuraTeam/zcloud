package k8s

import (
	"strings"
	"k8s.io/client-go/kubernetes"
	"github.com/astaxie/beego/logs"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
	"cloud/util"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strconv"
	"sort"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/sets"
	"fmt"
)

type NodeIp struct {
	Ip string
}

// 获取集群的IP地址
func GetNodesIp(clientset kubernetes.Clientset) []NodeIp {
	result := make([]NodeIp, 0)
	data := GetNodes(clientset, "")
	for _, d := range data {
		temp := NodeIp{}
		temp.Ip = d.Name
		result = append(result, temp)
	}
	return result
}

// 获取nodes
func GetNodes(clientset kubernetes.Clientset, labels string) []v1.Node {
	opt := metav1.ListOptions{}
	if labels != "" {
		opt.LabelSelector = labels
	}
	nodes, err := clientset.CoreV1().Nodes().List(opt)
	if err != nil {
		logs.Error("获取Nodes错误", err.Error())
		return make([]v1.Node, 0)
	}
	return nodes.Items
}

// 设置节点是否可调度
func UpdateNodeStatus(client kubernetes.Clientset, ip string, unschdulable bool) error {
	node := getNodes(client, ip)
	node.Spec.Unschedulable = unschdulable
	d, err := client.CoreV1().Nodes().Update(node)
	fmt.Println(d)
	if err != nil {
		logs.Error("更新节点失败", d, err)
	}
	return err
}

// 2018-02-12 16:26
// 获取node
func getNodes(clientset kubernetes.Clientset, ip string) *v1.Node {
	opt := metav1.GetOptions{}
	nodes, err := clientset.CoreV1().Nodes().Get(ip, opt)
	if err == nil {
		return nodes
	}
	return &v1.Node{}
}

// 2018-02-13 09:54
// 获取某个节点的镜像
func GetNodeImage(clustername string, ip string) []HostImages {
	cl, err := GetClient(clustername)
	result := make([]HostImages, 0)
	if err == nil {
		node := getNodes(cl, ip)
		count := 1
		for _, v := range node.Status.Images {
			if len(v.Names) > 1 {
				if !strings.Contains(v.Names[1], "none") {
					temp := HostImages{}
					temp.Id = count
					names := strings.Split(v.Names[1], ":")
					tname := names[0 : len(names)-1]
					temp.Name = strings.Join(tname, ":")
					temp.Tag = names[len(names)-1]
					temp.Size = strconv.FormatInt(v.SizeBytes/1024/1024, 10) + "MB"
					result = append(result, temp)
					count += 1
				}
			}
		}
	}
	return result
}

// 更新某个node的标签
// 2018-01-11 18:00
// k8s.UpdateNodeLabels("10.16.55.6","8080","10.16.55.102","sshd","sshd","")
func UpdateNodeLabels(clustername string, ip string, labelsData string) error {
	clientset, err := GetClient(clustername)
	opt := metav1.ListOptions{}
	opt.LabelSelector = "kubernetes.io/hostname=" + ip
	nodes, err := clientset.CoreV1().Nodes().List(opt)
	if err != nil {
		logs.Error("获取Node数据失败", err)
		return err
	}
	dclient, err := GetYamlClient(clustername, "", "v1", "api")
	if err != nil {
		return err
	}
	for _, item := range nodes.Items {
		lables := item.GetLabels()
		for _, v := range strings.Split(labelsData, "\n") {
			fmt.Println("获取到标签", v, ip)
			temp := strings.Split(v, "=")
			if len(temp) > 1 {
				lables[temp[0]] = temp[1]
				fmt.Println("设置标签", temp[0], temp[1])
			}
		}
		item.SetLabels(lables)
		resource := &metav1.APIResource{Name: "Nodes", Namespaced: false}

		conf := map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Node",
			"metadata": map[string]interface{}{
				"labels": item.GetLabels(),
				"name":   item.GetName(),
			},
			"spec": item.Spec,
		}
		fmt.Println("更新标签数据", util.ObjToString(conf))
		obj := unstructured.Unstructured{Object: conf}
		d, err := dclient.Resource(resource, "").Update(&obj)
		fmt.Println("更新Node标签:", d, err)
		return err
	}
	return err
}

func GetNodeFromCluster(clientSet kubernetes.Clientset) ClusterStatus {
	nodes := GetNodes(clientSet, "")
	clusterStatus := ClusterStatus{}
	for _, item := range nodes {
		if clusterStatus.MemSize == 0 && clusterStatus.CpuNum == 0 {
			clusterStatus.CpuNum = item.Status.Capacity.Cpu().Value()
			clusterStatus.MemSize = item.Status.Capacity.Memory().Value()
			clusterStatus.Nodes = 1
		} else {
			clusterStatus.CpuNum = clusterStatus.CpuNum + item.Status.Capacity.Cpu().Value()
			clusterStatus.MemSize = clusterStatus.MemSize + item.Status.Capacity.Memory().Value()
			clusterStatus.Nodes = clusterStatus.Nodes + 1
		}
		clusterStatus.OsVersion = item.Status.NodeInfo.OSImage
	}
	clusterStatus.PodNum = GetPodsNumber("", clientSet)
	clusterStatus.Services = GetServiceNumber(clientSet, "")
	clusterStatus.MemSize = clusterStatus.MemSize / 1024 / 1024 / 1024

	return clusterStatus
}

// 获取节点ready状态
func getNodeReady(d []v1.NodeCondition) (string, string) {
	for _, k := range d {
		if k.Type == "Ready" {
			if k.Status == "True" {
				return "运行中", k.Message
			} else {
				return "错误", k.Message
			}
		}
	}
	return "失败", "未知"
}

// 获取nodes状态数据
func GetNodesFromIp(ip string, clientset kubernetes.Clientset, nodes []v1.Node) NodeStatus {
	pods := GetPods("", clientset)
	//nodes := GetNodes(clientset)
	nodeStatus := NodeStatus{}
	for _, item := range nodes {
		if item.Name == ip {
			nodeStatus.CpuNum = item.Status.Capacity.Cpu().Value()
			nodeStatus.MemSize = item.Status.Capacity.Memory().Value() / 1024 / 1024 / 1024
			nodeStatus.PodNum = GetIpPodNumber(pods, ip)
			nodeStatus.HostIp = ip
			nodeStatus.CreateTime = util.ReplaceTime(item.CreationTimestamp.String())
			nodeStatus.Status, nodeStatus.ErrorMsg = getNodeReady(item.Status.Conditions)
			if item.Spec.Unschedulable == true {
				nodeStatus.Status = "不可调度"
			}
			nodeStatus.ImageNum = len(item.Status.Images)
			nodeStatus.K8sVersion = item.Status.NodeInfo.KubeletVersion
			nodeStatus.OsVersion = item.Status.NodeInfo.OSImage
			lables := make([]string, 0)
			for _, v := range item.Labels {
				if ! strings.Contains(v, "kubernetes") {
					lables = append(lables, v)
				}
			}
			nodeStatus.Lables = lables
			break
		}
	}
	return nodeStatus
}

var LEVEL_0 = 0
var LEVEL_1 = 1

func DescribeNodeResource(clientSet kubernetes.Clientset, ip string) []NodeReport {
	node := getNodes(clientSet, ip)
	pods := GetPodsFromNode(ip, clientSet)
	return describeNodeResource(pods, *node, ip)
}

// 2018-09-04 10:00
// 节点资源分配报表
func describeNodeResource(nodeNonTerminatedPodsList v1.PodList, node v1.Node, ip string) []NodeReport {
	result := make([]NodeReport, 0)
	info := make([]string, 0)
	info = append(info, fmt.Sprintf("Non-terminated Pods:\t(%d in total)\n", len(nodeNonTerminatedPodsList.Items)))
	info = append(info, fmt.Sprintf("Namespace\tName\t\tCPU Requests\tCPU Limits\tMemory Requests\tMemory Limits\n"))
	info = append(info, fmt.Sprintf("---------\t----\t\t------------\t----------\t---------------\t-------------\n"))
	allocatable := node.Status.Capacity
	if len(node.Status.Allocatable) > 0 {
		allocatable = node.Status.Allocatable
	}

	for _, pod := range nodeNonTerminatedPodsList.Items {
		req, limit := PodRequestsAndLimits(&pod)
		cpuReq, cpuLimit, memoryReq, memoryLimit := req[v1.ResourceCPU], limit[v1.ResourceCPU], req[v1.ResourceMemory], limit[v1.ResourceMemory]
		fractionCpuReq := float64(cpuReq.MilliValue()) / float64(allocatable.Cpu().MilliValue()) * 100
		fractionCpuLimit := float64(cpuLimit.MilliValue()) / float64(allocatable.Cpu().MilliValue()) * 100
		fractionMemoryReq := float64(memoryReq.Value()) / float64(allocatable.Memory().Value()) * 100
		fractionMemoryLimit := float64(memoryLimit.Value()) / float64(allocatable.Memory().Value()) * 100
		report := NodeReport{
			Ip: ip,
			Name:           pod.Name,
			Namespace:      pod.Namespace,
			CpuLimits:      cpuLimit.String() + " (" + strconv.FormatInt(int64(fractionCpuLimit), 10) + "%)",
			CpuRequests:    cpuReq.String() + " (" + strconv.FormatInt(int64(fractionCpuReq), 10) + "%)",
			MemoryLimits:   memoryLimit.String() + " (" + strconv.FormatInt(int64(fractionMemoryLimit), 10) + "%)",
			MemoryRequests: memoryReq.String() + "（" + strconv.FormatInt(int64(fractionMemoryReq), 10) + "%)",
		}
		result = append(result, report)
		r := fmt.Sprintf("%s\t%s\t\t%s (%d%%)\t%s (%d%%)\t%s (%d%%)\t%s (%d%%)\n", pod.Namespace, pod.Name,
			cpuReq.String(), int64(fractionCpuReq), cpuLimit.String(), int64(fractionCpuLimit),
			memoryReq.String(), int64(fractionMemoryReq), memoryLimit.String(), int64(fractionMemoryLimit))
		info = append(info, r)
	}

	reqs, limits := getPodsTotalRequestsAndLimits(nodeNonTerminatedPodsList)
	cpuReqs, cpuLimits, memoryReqs, memoryLimits := reqs[v1.ResourceCPU], limits[v1.ResourceCPU], reqs[v1.ResourceMemory], limits[v1.ResourceMemory]
	fractionCpuReqs := float64(0)
	fractionCpuLimits := float64(0)
	if allocatable.Cpu().MilliValue() != 0 {
		fractionCpuReqs = float64(cpuReqs.MilliValue()) / float64(allocatable.Cpu().MilliValue()) * 100
		fractionCpuLimits = float64(cpuLimits.MilliValue()) / float64(allocatable.Cpu().MilliValue()) * 100
	}
	fractionMemoryReqs := float64(0)
	fractionMemoryLimits := float64(0)
	if allocatable.Memory().Value() != 0 {
		fractionMemoryReqs = float64(memoryReqs.Value()) / float64(allocatable.Memory().Value()) * 100
		fractionMemoryLimits = float64(memoryLimits.Value()) / float64(allocatable.Memory().Value()) * 100
	}
	report := NodeReport{
		Ip:ip,
		Name:           "所有",
		Namespace:      "所有",
		CpuLimits:      cpuLimits.String() + " (" + strconv.FormatInt(int64(fractionCpuLimits), 10) + "%)",
		CpuRequests:    cpuReqs.String() + " (" + strconv.FormatInt(int64(fractionCpuReqs), 10) + "%)",
		MemoryLimits:   memoryLimits.String() + " (" + strconv.FormatInt(int64(fractionMemoryLimits), 10) + "%)",
		MemoryRequests: memoryReqs.String() + " (" + strconv.FormatInt(int64(fractionMemoryReqs), 10) + "%)",
	}
	result = append(result, report)
	info = append(info, fmt.Sprintf("%s\t%s (%d%%)\t%s (%d%%)\n",
		v1.ResourceCPU, cpuReqs.String(), int64(fractionCpuReqs), cpuLimits.String(), int64(fractionCpuLimits)))
	info = append(info, fmt.Sprintf("%s\t%s (%d%%)\t%s (%d%%)\n",
		v1.ResourceMemory, memoryReqs.String(), int64(fractionMemoryReqs), memoryLimits.String(), int64(fractionMemoryLimits)))
	extResources := make([]string, 0, len(allocatable))
	for resource := range allocatable {
		if !IsStandardContainerResourceName(string(resource)) && resource != v1.ResourcePods {
			extResources = append(extResources, string(resource))
		}
	}
	sort.Strings(extResources)
	for _, ext := range extResources {
		extRequests, extLimits := reqs[v1.ResourceName(ext)], limits[v1.ResourceName(ext)]
		info = append(info, fmt.Sprintf("%s\t%s\t%s\n", ext, extRequests.String(), extLimits.String()))
	}
	logs.Info(strings.Join(info, ""))
	return result
}

func getPodsTotalRequestsAndLimits(podList v1.PodList) (reqs map[v1.ResourceName]resource.Quantity, limits map[v1.ResourceName]resource.Quantity) {
	reqs, limits = map[v1.ResourceName]resource.Quantity{}, map[v1.ResourceName]resource.Quantity{}
	for _, pod := range podList.Items {
		podReqs, podLimits := PodRequestsAndLimits(&pod)
		for podReqName, podReqValue := range podReqs {
			if value, ok := reqs[podReqName]; !ok {
				reqs[podReqName] = *podReqValue.Copy()
			} else {
				value.Add(podReqValue)
				reqs[podReqName] = value
			}
		}
		for podLimitName, podLimitValue := range podLimits {
			if value, ok := limits[podLimitName]; !ok {
				limits[podLimitName] = *podLimitValue.Copy()
			} else {
				value.Add(podLimitValue)
				limits[podLimitName] = value
			}
		}
	}
	return
}

// PodRequestsAndLimits returns a dictionary of all defined resources summed up for all
// containers of the pod.
func PodRequestsAndLimits(pod *v1.Pod) (reqs v1.ResourceList, limits v1.ResourceList) {
	reqs, limits = v1.ResourceList{}, v1.ResourceList{}
	for _, container := range pod.Spec.Containers {
		addResourceList(reqs, container.Resources.Requests)
		addResourceList(limits, container.Resources.Limits)
	}
	// init containers define the minimum of any resource
	for _, container := range pod.Spec.InitContainers {
		maxResourceList(reqs, container.Resources.Requests)
		maxResourceList(limits, container.Resources.Limits)
	}
	return
}

// addResourceList adds the resources in newList to list
func addResourceList(list, new v1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = *quantity.Copy()
		} else {
			value.Add(quantity)
			list[name] = value
		}
	}
}

// maxResourceList sets list to the greater of list/newList for every resource
// either list
func maxResourceList(list, new v1.ResourceList) {
	for name, quantity := range new {
		if value, ok := list[name]; !ok {
			list[name] = *quantity.Copy()
			continue
		} else {
			if quantity.Cmp(value) > 0 {
				list[name] = *quantity.Copy()
			}
		}
	}
}

var standardContainerResources = sets.NewString(
	string(v1.ResourceCPU),
	string(v1.ResourceMemory),
	string(v1.ResourceEphemeralStorage),
)

// IsStandardContainerResourceName returns true if the container can make a resource request
// for the specified resource
func IsStandardContainerResourceName(str string) bool {
	return standardContainerResources.Has(str) || IsHugePageResourceName(v1.ResourceName(str))
}

// IsHugePageResourceName returns true if the resource name has the huge page
// resource prefix.
func IsHugePageResourceName(name v1.ResourceName) bool {
	return strings.HasPrefix(string(name), v1.ResourceHugePagesPrefix)
}

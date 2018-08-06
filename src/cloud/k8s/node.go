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
	logs.Info(d)
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
	result := []HostImages{}
	if err == nil {
		node := getNodes(cl, ip)
		count := 1
		for _, v := range node.Status.Images {
			if len(v.Names) > 1 {
				if !strings.Contains(v.Names[1], "none") {
					temp := HostImages{}
					temp.Id = count
					names := strings.Split(v.Names[1], ":")
					tname := names[0:len(names)-1]
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
			logs.Info("获取到标签", v, ip)
			temp := strings.Split(v, "=")
			if len(temp) > 1 {
				lables[temp[0]] = temp[1]
				logs.Info("设置标签", temp[0], temp[1])
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
			"spec":item.Spec,
		}
		logs.Info("更新标签数据" , util.ObjToString(conf))
		obj := unstructured.Unstructured{Object: conf}
		d, err := dclient.Resource(resource, "").Update(&obj)
		logs.Info("更新Node标签:", d, err)
		return err
	}
	return err
}

func GetNodeFromCluster(clientset kubernetes.Clientset) ClusterStatus {
	nodes := GetNodes(clientset, "")
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
	}
	clusterStatus.PodNum = GetPodsNumber("", clientset)
	clusterStatus.Services = GetServiceNumber(clientset, "")
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

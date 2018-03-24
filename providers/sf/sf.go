package sf

import (
  "github.com/virtual-kubelet/virtual-kubelet/manager"
  "k8s.io/api/core/v1"
  "fmt"
  "log"
  "time"
  "os"
  "k8s.io/apimachinery/pkg/api/resource"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SFProvider struct {
  nodeName string
  operatingSystem string
  internalIP string
  daemonEndpointPort int32
  pods map[string]*v1.Pod
}

func NewSFProvider(config string, rm *manager.ResourceManager, nodeName, operatingSystem string, internalIP string, daemonEndpointPort int32) (*SFProvider, error) {
  provider := SFProvider{
	  nodeName: nodeName,
	  operatingSystem: operatingSystem,
	  internalIP: internalIP,
	  daemonEndpointPort: daemonEndpointPort,
	  pods: make(map[string]*v1.Pod)}

  return &provider, nil
}

func (p *SFProvider) CreatePod(pod *v1.Pod) error {
   log.Printf("receive CreatePod %v \n", pod)
   log.Printf("MY_POD_NAME=%s\n",os.Getenv("MY_POD_NAME"))
   log.Printf("MY_POD_NAMESPACE=%s\n", os.Getenv("MY_POD_NAMESPACE"))
   log.Printf("VKUBELET_POD_ID=%s\n", os.Getenv("VKUBELET_POD_ID"))
   key, err := buildKey(pod)
   if err != nil {
     return err
   }
   p.pods[key] = pod
   return nil
}

func (p *SFProvider) UpdatePod(pod *v1.Pod) error {
  log.Printf("receive UpdatePod %q \n", pod.Name)

  key, err := buildKey(pod)
  if err != nil {
    return err
  }
  p.pods[key] = pod
  return nil
}

func (p *SFProvider) DeletePod(pod *v1.Pod) error {
  log.Printf("receive DeletedPod %v \n", pod)

  key, err := buildKey(pod)
  if err != nil {
    return err
  }
  delete(p.pods, key)
  return nil
}

func (p *SFProvider) GetPod(namespace, name string) (*v1.Pod, error) {
  log.Printf("receive GetPod %q \n", name)

  key, err := buildKeyFromNames(namespace, name)
  if err != nil {
    return nil, err
  }

  if pod, ok := p.pods[key]; ok {
    log.Printf(" returning: %v \n", p.pods[key])
    return pod, nil
  }
  return nil, nil
}

func (p *SFProvider) GetPodStatus(namespace, name string) (*v1.PodStatus, error) {
  log.Printf("receive GetPodStatus %q \n", name)
  now := metav1.NewTime(time.Now())

  status := &v1.PodStatus {
	  Phase: v1.PodRunning,
	  HostIP: "1.2.3.4",
	  PodIP: "5.6.7.8",
	  StartTime: &now,
	  Conditions: []v1.PodCondition {
		  {
			  Type: v1.PodInitialized,
			  Status: v1.ConditionTrue},
		  {
			  Type: v1.PodReady,
			  Status: v1.ConditionTrue},
		  {
			  Type: v1.PodScheduled,
			  Status: v1.ConditionTrue}}}
  pod, err := p.GetPod(namespace, name)
  if err != nil {
	  return status, err
  }
  for _, container := range pod.Spec.Containers {
	  status.ContainerStatuses = append(status.ContainerStatuses, v1.ContainerStatus{
		  Name: container.Name,
		  Image: container.Image,
		  Ready: true,
		  RestartCount: 0,
		  State: v1.ContainerState{
			  Running: &v1.ContainerStateRunning{
				  StartedAt: now}}})
  }
  return status, nil
}

func (p *SFProvider) GetPods() ([]*v1.Pod, error) {
  log.Printf("receive GetPods \n")

  var pods []*v1.Pod

  for _, pod := range p.pods {
	  log.Printf("adding %v \n", pod)
	  pods = append(pods, pod)
  }
  return pods,nil
}

func (p *SFProvider) Capacity() v1.ResourceList {
  return v1.ResourceList {
    "cpu": resource.MustParse("20"),
    "memory": resource.MustParse("100Gi"),
    "pods": resource.MustParse("20")}
}

func (p *SFProvider) NodeConditions() []v1.NodeCondition {
  return []v1.NodeCondition{
    {
      Type: "Ready",
      Status: v1.ConditionTrue,
      LastHeartbeatTime: metav1.Now(),
      LastTransitionTime: metav1.Now(),
      Reason: "KubeletReady",
      Message: "kubelet is ready."},
   {
      Type: "OutofDisk",
      Status: v1.ConditionFalse,
      LastHeartbeatTime: metav1.Now(),
      LastTransitionTime: metav1.Now(),
      Reason: "KubeletHasSufficientDisk",
      Message: "kubelet has sufficient disk space available"},
   {
      Type: "MemoryPressure",
      Status: v1.ConditionFalse,
      LastHeartbeatTime: metav1.Now(),
      LastTransitionTime: metav1.Now(),
      Reason: "KubeletHasSufficientMemory",
      Message: "kubelet has sufficient memory available"},
   {
      Type: "DiskPressure",
      Status: v1.ConditionFalse,
      LastHeartbeatTime: metav1.Now(),
      LastTransitionTime: metav1.Now(),
      Reason: "KubeletHasNoDiskPressure",
      Message: "kubelet has no disk pressure"},
   {
      Type: "NetworkUnavailable",
      Status: v1.ConditionFalse,
      LastHeartbeatTime: metav1.Now(),
      LastTransitionTime: metav1.Now(),
      Reason: "RouteCreated",
      Message: "RouteController created a route"}}
}

func (p *SFProvider) OperatingSystem() string {
  return p.operatingSystem
}

func (p *SFProvider) GetContainerLogs(namespace, podName, containerName string, tail int) (string, error) {
  log.Printf("receive GetContainerLogs %q \n", podName) 
  return "", nil
}

func (p *SFProvider) NodeAddresses() []v1.NodeAddress {
  return []v1.NodeAddress {
	  {
		  Type: "InternalIP",
		  Address: p.internalIP}}
}

func (p *SFProvider) NodeDaemonEndpoints() *v1.NodeDaemonEndpoints {
  return &v1.NodeDaemonEndpoints {
	  KubeletEndpoint: v1.DaemonEndpoint {
		  Port: p.daemonEndpointPort}}
}

func buildKey(pod *v1.Pod) (string, error) {
  if pod.ObjectMeta.Namespace == "" {
    return "", fmt.Errorf("pod namespace not found")
 }
 if pod.ObjectMeta.Name == "" {
    return "", fmt.Errorf("pod name not found")
 }

 return buildKeyFromNames(pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
}

func buildKeyFromNames(namespace string, name string) (string, error) {
  return fmt.Sprintf("%s-%s", namespace, name), nil
}

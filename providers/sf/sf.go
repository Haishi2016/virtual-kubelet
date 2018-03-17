package sf

import (
  "k8s.io/api/core/v1"
)

type SFProvider struct {
  operatingSystem string
}

func (p *SFProvider) CreatePod(pod *v1.Pod) error {
}

func (p *SFProvider) UpdatePod(pod *v1.Pod) error {
}

func (p *SFProvider) DeletePod(pod *v1.Pod) error {
}

func (p *SFProvider) GetPod(namespace, name string) (*v1.Pod, error) {
}

func (p *SFProvider) GetPodStatus(namespace, name string) (*v1.PodStatus, error) {
}

func (p *SFProvider) GetPods() ([]*v1.Pod, error) {
}

func (p *SFProvider) Capacity() v1.ResourceList {
}

func (p *SFProvider) NodeConditions() []v1.NodeCondition {
}

func (p *SFProvider) OperatingSystem() string {
  return p.operatingSystem
}


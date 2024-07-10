package k8s

import (
	"fmt"
	"os/exec"
)

type KubernetesDeployer struct{}

func NewKubernetesDeployer() *KubernetesDeployer {
	return &KubernetesDeployer{}
}

func (k *KubernetesDeployer) CreatePod(name string) error {
	cmd := exec.Command("kubectl", "run", name, "--image=nginx")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create pod: %w", err)
	}
	return nil
}

func (k *KubernetesDeployer) DeletePod(name string) error {
	cmd := exec.Command("kubectl", "delete", "pod", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete pod: %w", err)
	}

	return nil
}

func (k *KubernetesDeployer) GetAllPodList() ([]string, error) {
	return []string{"pod1", "pod2"}, nil
}

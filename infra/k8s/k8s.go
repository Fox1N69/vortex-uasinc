package k8s

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type KubernetesDeployer struct{}

func NewKubernetesDeployer() *KubernetesDeployer {
	return &KubernetesDeployer{}
}

func (k *KubernetesDeployer) CreatePod(name string, image string) error {
	cmd := exec.Command("kubectl", "run", name, image)
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
	cmd := exec.Command("kubectl", "get", "pods", "-o", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	var result struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
		} `json:"items"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	podNames := make([]string, len(result.Items))
	for i, item := range result.Items {
		podNames[i] = item.Metadata.Name
	}

	return podNames, nil
}

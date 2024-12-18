package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HelmDeploySpec defines the desired state of HelmDeploy
type HelmDeploySpec struct {
	RepositoryURL string   `json:"repositoryURL"` // GitHub repository to monitor
	TriggerBranch string   `json:"triggerBranch"` // Default: "release/*"
	HelmCharts    []string `json:"helmCharts"`    // List of Helm charts to deploy
}

// HelmDeployStatus defines the observed state of HelmDeploy
type HelmDeployStatus struct {
	DeployedCharts []string `json:"deployedCharts,omitempty"` // Track deployed charts
	LastSynced     string   `json:"lastSynced,omitempty"`     // Last time sync was successful
}

// HelmDeploy is the Schema for the helmdeploys API
type HelmDeploy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelmDeploySpec   `json:"spec,omitempty"`
	Status HelmDeployStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HelmDeployList contains a list of HelmDeploy
type HelmDeployList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelmDeploy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HelmDeploy{}, &HelmDeployList{})
}

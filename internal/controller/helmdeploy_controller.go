package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	fluxv1alpha1 "github..com/ADugal1/release-operator/api/v1alpha1"
)

// HelmDeployReconciler reconciles a HelmDeploy object
type HelmDeployReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type ReleasePayload struct {
	Action  string `json:"action"`
	Release struct {
		TagName string `json:"tag_name"`
	} `json:"release"`
}

func (r *HelmDeployReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Fetch the HelmDeploy resource
	var helmDeploy fluxv1alpha1.HelmDeploy
	if err := r.Get(ctx, req.NamespacedName, &helmDeploy); err != nil {
		log.Error(err, "unable to fetch HelmDeploy")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Example: Check if GitHub release exists (pseudo-code, replace with actual API call)
	if isReleaseAvailable(helmDeploy.Spec.RepositoryURL, helmDeploy.Spec.TriggerBranch) {
		for _, chart := range helmDeploy.Spec.HelmCharts {
			err := deployHelmChart(chart)
			if err != nil {
				log.Error(err, "failed to deploy Helm chart", "chart", chart)
				return ctrl.Result{}, err
			}
		}
		// Update the status
		helmDeploy.Status.DeployedCharts = helmDeploy.Spec.HelmCharts
		helmDeploy.Status.LastSynced = time.Now().Format(time.RFC3339)
		if err := r.Status().Update(ctx, &helmDeploy); err != nil {
			log.Error(err, "failed to update HelmDeploy status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{RequeueAfter: time.Hour}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HelmDeployReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&fluxv1alpha1.HelmDeploy{}).
		Complete(r)
}

func deployHelmChart(chartName string) error {
	// Use Helm SDK to deploy the chart
	// Placeholder for now
	log.Info("Deploying Helm chart", "chart", chartName)
	return nil
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var payload ReleasePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if payload.Action == "published" {
		log.Printf("New release detected: %s", payload.Release.TagName)
		// Trigger deployment logic here
	}
}

func fetchHelmChart(chartName, releaseTag string) (string, error) {
	nexusURL := fmt.Sprintf("https://nexus.example.com/repository/helm/%s-%s.tgz", chartName, releaseTag)
	resp, err := http.Get(nexusURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	filePath := fmt.Sprintf("/tmp/%s-%s.tgz", chartName, releaseTag)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return filePath, err
}

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	configclient "github.com/openshift/client-go/config/clientset/versioned"
	adminrelease "github.com/openshift/oc/pkg/cli/admin/release"
)

func main() {
	ctx := context.Background()

	// Load Kubernetes configuration from KUBECONFIG environment variable or default location
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		kubeconfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("Failed to load kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create OpenShift config client
	configClient, err := configclient.NewForConfig(config)
	if err != nil {
		fmt.Printf("Failed to create OpenShift config client: %v\n", err)
		os.Exit(1)
	}

	// Retrieve the ClusterVersion object
	clusterVersion, err := configClient.ConfigV1().ClusterVersions().Get(ctx, "version", metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Failed to get ClusterVersion: %v\n", err)
		os.Exit(1)
	}

	releaseImage := clusterVersion.Status.Desired.Image
	fmt.Printf("Release Image: %s\n", releaseImage)

	info := adminrelease.InfoOptions{}

	release, err := info.LoadReleaseInfo(releaseImage, false)
	if err != nil {
		fmt.Printf("Failed to load release info: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Release: %+v\n", release)
}

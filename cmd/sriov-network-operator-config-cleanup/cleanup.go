package main

import (
	"context"
	"os"

	snolog "github.com/k8snetworkplumbingwg/sriov-network-operator/pkg/log"
	"github.com/spf13/cobra"
	"k8s.io/client-go/dynamic"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/log"

	sriovnetworkv1 "github.com/k8snetworkplumbingwg/sriov-network-operator/pkg/client/clientset/versioned/typed/sriovnetwork/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	namespace string
)

func init() {
	rootCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "designated SriovOperatorConfig namespace")
}

func DynamicClientFor(g schema.GroupVersionKind, obj *unstructured.Unstructured, namespace string) (dynamic.ResourceInterface, error) {
	return nil, nil
}

func runCleanupCmd(cmd *cobra.Command, args []string) error {
	var (
		config *rest.Config
		err    error
	)
	// init logger
	snolog.InitLog()
	setupLog := log.Log.WithName("sriov-network-operator-config-cleanup")

	setupLog.Info("Run sriov-network-operator-config-cleanup")

	// creates the in-cluster config
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		setupLog.Error(err, "failed initialization k8s rest config")
	}

	sriovcs, err := sriovnetworkv1.NewForConfig(config)
	if err != nil {
		setupLog.Error(err, "failed to create 'sriovnetworkv1' clientset")
	}

	sriovcs.SriovOperatorConfigs(namespace).Delete(context.Background(), "default", metav1.DeleteOptions{})
	if err != nil {
		setupLog.Error(err, "failed to delete SriovOperatorConfig")
		return err
	}

	return nil

}

package main

import (
	"context"

	sriovnetworkv1 "github.com/k8snetworkplumbingwg/sriov-network-operator/api/v1"
	"github.com/k8snetworkplumbingwg/sriov-network-operator/test/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var (
	testNamespace            string = "sriov-network-operator"
	defaultSriovOperatorSpec        = sriovnetworkv1.SriovOperatorConfigSpec{
		EnableInjector:        true,
		EnableOperatorWebhook: true,
		LogLevel:              2,
		FeatureGates:          nil,
	}
)

var _ = Describe("cleanup", Ordered, func() {

	defaultSriovOperatorConfig := &sriovnetworkv1.SriovOperatorConfig{
		ObjectMeta: metav1.ObjectMeta{Name: "default", Namespace: testNamespace},
		Spec:       defaultSriovOperatorSpec,
	}

	BeforeEach(func() {
		Expect(k8sClient.Create(context.Background(), defaultSriovOperatorConfig)).NotTo(HaveOccurred())
		err := util.WaitForNamespacedObject(defaultSriovOperatorConfig, k8sClient, testNamespace, "default", util.RetryInterval, util.APITimeout*3)
		Expect(err).NotTo(HaveOccurred())
	})

	It("test webhook cleanup flow", func() {
		cmd := &cobra.Command{}
		namespace = testNamespace
		Expect(runCleanupCmd(cmd, []string{})).Should(Succeed())

		config := &sriovnetworkv1.SriovOperatorConfig{}
		Expect(k8sClient.Get(context.Background(), types.NamespacedName{Name: "default", Namespace: testNamespace}, config)).To(MatchError(
			ContainSubstring("sriovoperatorconfigs.sriovnetwork.openshift.io \"default\" not found")))

	})
})

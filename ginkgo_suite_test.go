package e2e_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Define timeout and interval for waiting checks
const waitTimeout = time.Second * 30
const waitInterval = time.Second

func TestGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ginkgo Suite")
}

var _ = Describe("CRD Presence", func() {
	var (
		apiExtClient *apiextensionsclient.Clientset
		config       *rest.Config
	)

	BeforeSuite(func() {
		// Load kubeconfig from the default location or use an in-cluster config
		var err error
		config, err = loadKubeConfig()
		Expect(err).NotTo(HaveOccurred())

		// Create API extensions client
		apiExtClient, err = apiextensionsclient.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should find the Custom Resource Definition", func() {
		// Specify the name of your CRD
		crdName := "your-crd-name"

		// Check if the CRD is present
		Eventually(func() bool {
			_, err := apiExtClient.ApiextensionsV1().CustomResourceDefinitions().Get(context.TODO(), crdName, metav1.GetOptions{})
			return err == nil
		}, waitTimeout, waitInterval).Should(BeTrue(), fmt.Sprintf("CRD %s not found", crdName))
	})

	It("should list Custom Resource Definitions", func() {
		// List all CRDs
		crdList, err := apiExtClient.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())

		// Print CRD names
		for _, crd := range crdList.Items {
			fmt.Printf("Found CRD: %s\n", crd.GetName())
		}
	})
})

// loadKubeConfig loads the kubeconfig from the default location or uses in-cluster config if running in a pod.
func loadKubeConfig() (*rest.Config, error) {
	home := homedir.HomeDir()
	kubeconfigPath := os.Getenv("KUBECONFIG")
	log.Println("kubeconfigPath : " + kubeconfigPath)
	if home != "" && kubeconfigPath == "" {
		kubeconfigPath = home + "/.kube/config"
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil && kubeconfigPath == "" {
		config, err = rest.InClusterConfig()
	}
	return config, err
}

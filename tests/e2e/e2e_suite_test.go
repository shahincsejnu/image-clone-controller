package e2e_test

import (
	"flag"
	"github.com/shahincsejnu/image-clone-controller/tests/e2e/framework"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	clientSetScheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/scale/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	kubeconfigPath = func() string {
		kubecfg := os.Getenv("KUBECONFIG")
		if kubecfg != "" {
			return kubecfg
		}
		return filepath.Join(homedir.HomeDir(), ".kube", "config")
	}()

	registry = "shahincsejnu"

	whichController = func() string {
		whichCon := os.Getenv("WHICH_CONTROLLER")
		if whichCon != "" {
			return whichCon
		}
		return framework.DEPLOYMENT
	}()
)

func init() {
	utilruntime.Must(scheme.AddToScheme(clientSetScheme.Scheme))

	//flag.StringVar(&kubeconfigPath, "kubeconfig", kubeconfigPath, "Path to kubeconfig file")
	flag.StringVar(&whichController, "which-controller", whichController, "Define which controller you want to check")
	flag.StringVar(&registry, "registry", registry, "your dockerhub username")
}

const (
	TIMEOUT = 20 * time.Minute
)

var (
	root *framework.Framework
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TIMEOUT)

	RunSpecs(t, "e2e Suite test")
}

var _ = BeforeSuite(func() {
	By("Using kubeconfig from " + kubeconfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	Expect(err).NotTo(HaveOccurred())
	config.Burst = 100
	config.QPS = 100

	// Clients
	kubeClient := kubernetes.NewForConfigOrDie(config)

	// Framework
	root = framework.New(config, kubeClient)

	// Create namespace
	By("Using namespace " + root.Namespace())
	err = root.CreateNamespace()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	By("Deleting Namespace")
	err := root.DeleteNamespace()
	Expect(err).NotTo(HaveOccurred())
})

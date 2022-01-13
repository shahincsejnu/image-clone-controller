package framework

import (
	"gomodules.xyz/x/crypto/rand"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	ALL        = "all"
	DEPLOYMENT = "deployment"
	DAEMONSET  = "daemonset"
)

type Framework struct {
	restConfig *rest.Config
	kubeClient kubernetes.Interface
	namespace  string
	name       string
}

func New(
	restConfig *rest.Config,
	kubeClient kubernetes.Interface,
) *Framework {
	return &Framework{
		restConfig: restConfig,
		kubeClient: kubeClient,
		name:       "sample-obj",
		namespace:  "image-clone-controller-system",
	}
}

func (f *Framework) Invoke() *Invocation {
	return &Invocation{
		Framework: f,
		app:       "apiserver",
	}
}

func (fi *Invocation) GetRandomName(extraSuffix string) string {
	return rand.WithUniqSuffix(fi.name + extraSuffix)
}

func RunTest(controller, whichController string) bool {
	if whichController == ALL || controller == whichController {
		return true
	} else {
		return false
	}
}

type Invocation struct {
	*Framework
	app string
}

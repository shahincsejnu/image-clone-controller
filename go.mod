module github.com/shahincsejnu/image-clone-controller

go 1.16

require (
	github.com/go-logr/logr v0.4.0
	github.com/google/go-containerregistry v0.5.2-0.20210609162550-f0ce2270b3b4
	github.com/google/go-containerregistry/pkg/authn/k8schain v0.0.0-20220107180046-bf65fd766231
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	k8s.io/api v0.22.1
	k8s.io/apimachinery v0.22.1
	k8s.io/client-go v0.22.1
	k8s.io/klog/v2 v2.9.0
	kmodules.xyz/client-go v0.0.0-20210822203828-5e9cebbf1dfa
	sigs.k8s.io/controller-runtime v0.10.0
)

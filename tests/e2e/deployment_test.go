package e2e_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/shahincsejnu/image-clone-controller/tests/e2e/framework"
	appsv1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

var _ = Describe("Deployment", func() {
	var (
		err error
		f   *framework.Invocation
	)
	BeforeEach(func() {
		f = root.Invoke()
		if !framework.RunTest(framework.DEPLOYMENT, whichController) {
			Skip(fmt.Sprintf("`%s` test is applied only when whichController flag is either `all` or `%s` but got `%s`", framework.DEPLOYMENT, framework.DEPLOYMENT, whichController))
		}
	})

	Describe("ImageCloneController", func() {
		Context("DeploymentController", func() {
			var (
				dockerCred     *core.Secret
				secretName     string
				deploymentName string
				deployment     *appsv1.Deployment
			)

			BeforeEach(func() {
				secretName = "image-clone-controller-cred"
				dockerCred = f.DockerCredSecret(secretName)
				deploymentName = f.GetRandomName("")
				deployment = f.Deployment(deploymentName)
			})

			AfterEach(func() {
				By("Deleting Deployment")
				err = f.DeleteDeployment(deployment.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())

				By("Wait for Deleting Deployment")
				f.EventuallyDeploymentDeleted(deployment.ObjectMeta).Should(BeTrue())

				By("Deleting secret")
				err = f.DeleteSecret(dockerCred.ObjectMeta)
			})

			It("should create, image clone and delete Deployment successfully", func() {
				By("Creating Secret")
				err = f.CreateSecret(dockerCred)
				Expect(err).NotTo(HaveOccurred())

				By("Creating Deployment")
				err = f.CreateDeployment(deployment)
				Expect(err).NotTo(HaveOccurred())

				By("Wait for Image Get cloned")
				f.EventuallyImageCloned(deployment.ObjectMeta, registry+"/").Should(BeTrue())
			})
		})
	})
})

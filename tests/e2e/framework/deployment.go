package framework

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	"strings"
	"time"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	meta_util "kmodules.xyz/client-go/meta"
)

func (i *Invocation) Deployment(name string) *appsv1.Deployment {
	image := "sakibalamin/apiserver:1.0.1"
	replicas := int32(2)

	return &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: i.Namespace(),
			Labels: map[string]string{
				"app": i.app,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "apiserver",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"app": "apiserver",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "apiserver",
							Image: image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (f *Framework) CreateDeployment(obj *appsv1.Deployment) error {
	_, err := f.kubeClient.AppsV1().Deployments(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	return err
}

func (f *Framework) DeleteDeployment(meta metav1.ObjectMeta) error {
	return f.kubeClient.AppsV1().Deployments(meta.Namespace).Delete(context.TODO(), meta.Name, meta_util.DeleteInBackground())
}

func (f *Framework) EventuallyImageCloned(meta metav1.ObjectMeta, registry string) GomegaAsyncAssertion {
	return Eventually(
		func() bool {
			deployment, err := f.kubeClient.AppsV1().Deployments(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			tmp := true

			for _, container := range deployment.Spec.Template.Spec.Containers {
				img := container.Image
				if !strings.HasPrefix(img, registry) {
					tmp = false
					break
				}
			}
			return tmp
		},
		time.Minute*15,
		time.Second*10,
	)
}

func (f *Framework) EventuallyDeploymentDeleted(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(
		func() bool {
			_, err := f.kubeClient.AppsV1().Deployments(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{})
			return errors.IsNotFound(err)
		},
		time.Minute*15,
		time.Second*10,
	)
}

/*
Copyright @shahincsejnu 2022.
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	meta_util "kmodules.xyz/client-go/meta"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("deployment", req.NamespacedName)

	// We must ignore the Deployments of "kube-system" namespace
	// could have ignore "kube-system" namespace by checking req.Namespace here
	// but Ignored from SetupWithManager function by Event filter

	// Getting the Deployment Object
	obj := &appsv1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		log.Error(err, "unable to fetch Deployment")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	for index, container := range obj.Spec.Template.Spec.Containers {
		img := container.Image
		// Ignore containers who are using cloned image
		if strings.HasPrefix(img, "shahincsejnu/") {
			continue
		}
		// Add "clone/" as prefix of the image for marking that it's cloned image
		modifiedImage := "shahincsejnu/" + strings.ReplaceAll(img, "/", "-")

		err := imagePull(img)
		if err != nil {
			return ctrl.Result{}, err
		}
		err = imageTag(img, modifiedImage)
		if err != nil {
			return ctrl.Result{}, err
		}
		err = imagePush(modifiedImage)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Use cloned image
		obj.Spec.Template.Spec.Containers[index].Image = modifiedImage
	}

	// Update the Deployment Object
	err := r.Update(ctx, obj)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithEventFilter(predicate.Funcs{
			CreateFunc: func(e event.CreateEvent) bool {
				return !meta_util.MustAlreadyReconciled(e.Object)
			},
			UpdateFunc: func(e event.UpdateEvent) bool {
				return (e.ObjectNew.(metav1.Object)).GetDeletionTimestamp() != nil || !meta_util.MustAlreadyReconciled(e.ObjectNew)
			},
		}).
		WithEventFilter(predicate.NewPredicateFuncs(func(e client.Object) bool {
			if e.GetNamespace() == "kube-system" {
				klog.Infof("Ignoring kube-system namespace's events")
				return false
			}
			return true
		})).
		Complete(r)
}

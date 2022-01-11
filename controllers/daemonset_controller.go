/*
Copyright @shahincsejnu 2022.
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/crane"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	meta_util "kmodules.xyz/client-go/meta"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
)

// DaemonSetReconciler reconciles a DaemonSet object
type DaemonSetReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=daemonsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=daemonsets/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=*,resources=*,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *DaemonSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("daemonset", req.NamespacedName)

	// We must ignore the Deployments of "kube-system" namespace
	// could have ignore "kube-system" namespace by checking req.Namespace here
	// but Ignored from SetupWithManager function by Event filter

	// Getting the DaemonSet Object
	obj := &appsv1.DaemonSet{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		log.Error(err, "unable to fetch DaemonSet")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// traverse all of the containers of deployment to check & clone the images
	for index, container := range obj.Spec.Template.Spec.Containers {
		img := container.Image
		// Ignore containers who are using cloned image
		if strings.HasPrefix(img, "shahincsejnu/") {
			continue
		}
		// Add "shahincsejnu/" as prefix of the image for marking that it's cloned image
		modifiedImage := "shahincsejnu/" + strings.ReplaceAll(img, "/", "-")

		// copy the modified image to own repository using crane.Copy
		err := crane.Copy(container.Image, modifiedImage)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Use cloned image
		obj.Spec.Template.Spec.Containers[index].Image = modifiedImage
	}

	// Update the DaemonSet Object
	err := r.Update(ctx, obj)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DaemonSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.DaemonSet{}).
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

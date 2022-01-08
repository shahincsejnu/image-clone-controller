# image-clone-controller
Kubernetes controller which watches applications (Deployment and DaemonSet) and "caches" the images (public container images) by re-uploading to our own registry repository and reconfiguring the applications to use these copies.

## Project's Motivation
- We’d like to be safe against the risk of public container images disappearing from the registry while we use them, breaking our deployments. 
- Suppose, we have a Kubernetes cluster on which we can run applications. These applications will often use publicly available container images, like official images of popular programs, e.g. Jenkins, PostgreSQL, and so on. Since the images reside in repositories over which we have no control, it is possible that the owner of the repo deletes the image while our pods are configured to use it. In the case of a subsequent node rotation, the locally cached copies of the images would be deleted and Kubernetes would be unable to re-download them in order to re-provision the applications.
- So, we want to have a controller which watches the applications and “caches” the images by re-uploading to our own registry repository and reconfiguring the applications to use these copies.
 

## Resources:
- https://book.kubebuilder.io/quick-start.html
- https://github.com/kubernetes-sigs/controller-runtime  
- https://github.com/kubernetes-sigs/controller-runtime/tree/master/examples/builtins  
- https://github.com/google/go-containerregistry/blob/master/pkg/v1/remote/README.md  
- https://godoc.org/github.com/google/go-containerregistry/pkg/v1/remote  




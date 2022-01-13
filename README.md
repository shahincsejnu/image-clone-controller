# image-clone-controller
Kubernetes controller which watches applications (Deployment and DaemonSet) and "caches" the images (public container images) by re-uploading to our own registry repository and reconfiguring the applications to use these copies.

## Project's Motivation
- We’d like to be safe against the risk of public container images disappearing from the registry while we use them, breaking our deployments. 
- Suppose, we have a Kubernetes cluster on which we can run applications. These applications will often use publicly available container images, like official images of popular programs, e.g. Jenkins, PostgreSQL, and so on. Since the images reside in repositories over which we have no control, it is possible that the owner of the repo deletes the image while our pods are configured to use it. In the case of a subsequent node rotation, the locally cached copies of the images would be deleted and Kubernetes would be unable to re-download them in order to re-provision the applications.
- So, we want to have a controller which watches the applications and “caches” the images by re-uploading to our own registry repository and reconfiguring the applications to use these copies.

## Demo
- Watch the working demo: https://asciinema.org/a/hxq0A6vjNBJKq940t63cXUgkU

## Use

### Locally running the manager
- clone this repo
- open the repo locally
- run `make`
- run `./bin/manager`
- open another terminal and go to samples: `cd config/samples`
- apply docker cred secret & sample deployment: 
    - give <base64 encoded username:password of your docker registry> in the `auth:` field of the docker-cred-k8s-secret 
    - run `kubectl apply -f docker-cred-secret.yaml`
    - run `kubectl apply -f sample-deployment.yaml`
- check in the sample deployment image, it will get cloned & pushed to your given docker registry and re-use in the deployment

### InCluster manager running
- `export IMG="<your_registry>/<controller_image_name>:<tag>"`
- `make docker-build`
- `make docker-push` (Note: for docker push you need to login in your dockerhub from the current terminal by `docker login`)
- `make deploy`
- verify the deployment by: `kubectl get all -n image-clone-controller-system`  
- open another terminal and go to samples: `cd config/samples`
- apply docker cred secret & sample deployment:
    - give <base64 encoded username:password of your docker registry> in the `auth:` field of the docker-cred-k8s-secret
    - run `kubectl apply -f docker-cred-secret.yaml`
    - run `kubectl apply -f sample-deployment.yaml`
- check in the sample deployment image, it will get cloned & pushed to your given docker registry and re-use in the deployment
- undeploy by: `make undeploy`

## e2e test
- Added e2e test for deployment controller, similarly will add for DaemonSet controller
- For using Deployment controller test follow below steps:
  - run the controller (either locally or incluster running the manager)
  - in another terminal go to project's : `cd tests/e2e`
  - in the `tests/e2e/framework/docker-cred-secret.go` file provide your dockerhub "username:password" in the "auth" field
  - run `ginkgo run --which-controller=<controller_name> --registry=<your_dockerhub_username>`
  - ex: `ginkgo run -- --which-controller=deployment --registry=shahincsejnu`
  - Note: make sure you sync the namespace, registry name among test files & controllers
  
## Disclaimer
- It's a hobby project, not a production grade

## What's Next?
- make this controller code more generic 
- make helm chart of this operator

## Resources:
- https://book.kubebuilder.io/quick-start.html
- https://github.com/kubernetes-sigs/controller-runtime  
- https://github.com/kubernetes-sigs/controller-runtime/tree/master/examples/builtins  
- https://github.com/google/go-containerregistry/blob/master/pkg/v1/remote/README.md  
- https://godoc.org/github.com/google/go-containerregistry/pkg/v1/remote  
- https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
- https://github.com/google/go-containerregistry/blob/main/pkg/authn/k8schain/README.md
- https://github.com/google/go-containerregistry/blob/main/cmd/crane/recipes.md
- https://github.com/google/go-containerregistry/blob/main/cmd/crane/doc/crane_copy.md




## Problem Statement:
We use namespace isolation to allow engineering teams to share Kubernetes resources safely. But DNS is global so we need to put guardrails in place so that no two services can use the same ingress hostname. How can we enforce that within the cluster?

## Expected Output:
Your code should be stored in a publicly accessible location.

### Dev Environment
- minikube
- minikube ingress-dns plugin

### Outcome
- some ingress controller will, by default, prevent duplicate hostnames in different namespaces with their admission controllers
- we will use a the ValidatingAdmissionWebhook admission controller to enforce an fqdn policy 

### References & Credit
- https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#validatingadmissionwebhook 
- https://docs.giantswarm.io/advanced/custom-admission-controller/
- https://minikube.sigs.k8s.io/docs/handbook/addons/ingress-dns
- https://github.com/open-policy-agent/gatekeeper
- https://github.com/joelspeed/webhook-certificate-generator
- https://github.com/chipmk/docker-mac-net-connect
- https://github.com/slok/kubewebhook/

### Helpful commands
```
# load docker image to minikube
minikube image load hostnamevalidator

# approve cert generating by wcg
kubectl certificate approve validator.default.svc

# generate certs for the server and ca bundle for the VAW
go run cmd/webhook-certificate-generator/main.go --service-name=validator --namespace=default --secret-name=validator-certs --in-cluster=false -k ~/.kube/config --patch-validating=validator
```
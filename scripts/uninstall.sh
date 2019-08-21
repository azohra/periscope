kubectl delete -f $(pwd)/deployment/k8s/periscope/deployment.yaml
kubectl delete -f $(pwd)/deployment/k8s/periscope/destination-rule.yaml
kubectl delete -f $(pwd)/deployment/k8s/periscope/gateway.yaml
kubectl delete -f $(pwd)/deployment/k8s/periscope/role.yaml
kubectl delete -f $(pwd)/deployment/k8s/periscope/rolebinding.yaml 
kubectl delete -f $(pwd)/deployment/k8s/periscope/service-account.yaml
kubectl delete -f $(pwd)/deployment/k8s/periscope/service.yaml
kubectl delete -f $(pwd)/deployment/k8s/periscope/virtualservice.yaml

#kubectl delete -f $(pwd)/deployment/k8s/namespace.yaml
kubectl apply -f $(pwd)/deployment/k8s/periscope/deployment.yaml
kubectl apply -f $(pwd)/deployment/k8s/periscope/destination-rule.yaml
kubectl apply -f $(pwd)/deployment/k8s/periscope/gateway.yaml
kubectl apply -f $(pwd)/deployment/k8s/periscope/role.yaml
kubectl apply -f $(pwd)/deployment/k8s/periscope/rolebinding.yaml 
kubectl apply -f $(pwd)/deployment/k8s/periscope/service-account.yaml
kubectl apply -f $(pwd)/deployment/k8s/periscope/service.yaml
kubectl apply -f $(pwd)/deployment/k8s/periscope/virtualservice.yaml

#kubectl apply -f $(pwd)/deployment/k8s/namespace.yaml
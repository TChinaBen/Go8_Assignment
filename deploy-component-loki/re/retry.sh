#!/bin/bash

NAMESPACE=loki

kubectl delete -n $NAMESPACE -f loki2.yaml
#kubectl delete -n $NAMESPACE -f ruler.yaml
kubectl delete -n $NAMESPACE -f promtail.yaml


sleep 10

kubectl apply -n $NAMESPACE -f promtail.yaml
kubectl apply -n $NAMESPACE -f loki2.yaml
# kubectl apply -n $NAMESPACE -f ruler.yaml


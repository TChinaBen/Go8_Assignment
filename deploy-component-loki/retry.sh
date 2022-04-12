#!/bin/bash

NAMESPACE=loki

kubectl delete -n $NAMESPACE -f loki.yaml
kubectl delete -n $NAMESPACE -f ruler.yaml
kubectl delete -n $NAMESPACE -f promtail.yaml

sleep 5

kubectl apply -n $NAMESPACE -f loki.yaml
kubectl apply -n $NAMESPACE -f ruler.yaml
kubectl apply -n $NAMESPACE -f promtail.yaml


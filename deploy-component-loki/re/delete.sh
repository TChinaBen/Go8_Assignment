#!/bin/bash

NAMESPACE=loki

kubectl delete -n $NAMESPACE -f object-bucket-claim-delete.yaml
kubectl delete -n $NAMESPACE -f loki.yaml
kubectl delete -n $NAMESPACE -f ruler.yaml
kubectl delete -n $NAMESPACE -f promtail.yaml
kubectl delete ns $1

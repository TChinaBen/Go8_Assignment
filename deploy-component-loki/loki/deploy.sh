#!/bin/bash

NAMESPACE=loki

# LOKI_IMAGE=docker.io/grafana/loki:2.2.0
# PROMTAIL_IMAGE=docker.io/grafana/promtail:2.2.0
# NGINX_IMAGE=docker.io/nginxinc/nginx-unprivileged:1.19-alpine
# SYNCER_IMAGE=docker.io/makoto126/rule-syncer:v1

if [[ $HARBOR_URL == "" ]]; then
    echo "harbor url required"
    exit 1
fi

if [[ $MODE != "fs" && $MODE != "s3" ]]; then
    echo "mode should be fs or s3"
    exit 1
fi

LOKI_IMAGE=$HARBOR_URL/loki:2.2.0
PROMTAIL_IMAGE=$HARBOR_URL/promtail:2.2.0
NGINX_IMAGE=$HARBOR_URL/nginx-unprivileged:1.19-alpine
SYNCER_IMAGE=$HARBOR_URL/rule-syncer:v1

# if [[ $LOKI_IMAGE == "" ]]; then
#     echo "loki image required"
#     exit 1
# fi
# 
# if [[ $PROMTAIL_IMAGE == "" ]]; then
#     echo "promtail image required"
#     exit 1
# fi
# 
# if [[ $NGINX_IMAGE == "" ]]; then
#     echo "nginx image required"
#     exit 1
# fi
# 
# if [[ $SYNCER_IMAGE == "" ]]; then
#     echo "rule-syncer image required"
#     exit 1
# fi
#

## create ns
kubectl create ns $NAMESPACE

## create secret docker-registry
kubectl create secret docker-registry -n $NAMESPACE registrysecret --docker-server=$REPO --docker-username=$USERNAME --docker-password=$PASSWORD

if [[ $MODE == "s3" ]]; then
    ## request object
    kubectl apply -n $NAMESPACE -f object-bucket-claim-delete.yaml
    ## wait rook to sync
    sleep 3

    ## get s3 info
    S3_BUCKET_HOST=$(kubectl -n rook-ceph get svc rook-ceph-rgw-my-store -o jsonpath='{.spec.clusterIP}')
    S3_BUCKET_NAME=$(kubectl -n $NAMESPACE get cm loki-distributed-object-storage -o jsonpath='{.data.BUCKET_NAME}')
    S3_ACCESS_KEY_ID=$(kubectl -n $NAMESPACE get secret loki-distributed-object-storage -o jsonpath='{.data.AWS_ACCESS_KEY_ID}' | base64 --decode)
    S3_SECRET_ACCESS_KEY=$(kubectl -n $NAMESPACE get secret loki-distributed-object-storage -o jsonpath='{.data.AWS_SECRET_ACCESS_KEY}' | base64 --decode)

    S3="s3://$S3_ACCESS_KEY_ID:$S3_SECRET_ACCESS_KEY@$S3_BUCKET_HOST:80/$S3_BUCKET_NAME"

    ## deploy loki
    sed "s!{{S3}}!$S3!g; s!{{NAMESPACE}}!$NAMESPACE!g; s!{{LOKI_IMAGE}}!$LOKI_IMAGE!g; s!{{NGINX_IMAGE}}!$NGINX_IMAGE!g" loki.yaml.template > loki.yaml
    kubectl apply -n $NAMESPACE -f loki.yaml
else
    sed "s!{{NAMESPACE}}!$NAMESPACE!g; s!{{LOKI_IMAGE}}!$LOKI_IMAGE!g; s!{{NGINX_IMAGE}}!$NGINX_IMAGE!g" loki-fs.yaml.template > loki.yaml
    kubectl apply -n $NAMESPACE -f loki.yaml
fi

## deploy ruler
sed "s!{{NAMESPACE}}!$NAMESPACE!g; s!{{LOKI_IMAGE}}!$LOKI_IMAGE!g; s!{{SYNCER_IMAGE}}!$SYNCER_IMAGE!g" ruler.yaml.template > ruler.yaml
kubectl apply -n $NAMESPACE -f ruler.yaml

## deploy promtail
sed "s!{{NAMESPACE}}!$NAMESPACE!g; s!{{PROMTAIL_IMAGE}}!$PROMTAIL_IMAGE!g" promtail.yaml.template > promtail.yaml
kubectl apply -n $NAMESPACE -f promtail.yaml


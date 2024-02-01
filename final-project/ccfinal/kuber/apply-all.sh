#!/bin/bash

kubectl apply -f ./database-secret.yaml
kubectl apply -f ./database-sts.yaml
kubectl apply -f ./http-monitor-configmap.yml
kubectl apply -f ./http-monitor-deployment.yml
kubectl apply -f ./http-monitor-service.yml

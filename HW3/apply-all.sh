#!/bin/bash

kubectl apply -f redis-config.yaml
kubectl apply -f redis-pv.yaml
kubectl apply -f redis-pvc.yaml
kubectl apply -f redis-secret.yaml
kubectl apply -f weather-app-config.yaml
kubectl apply -f redis.yaml
kubectl apply -f weather-app.yaml

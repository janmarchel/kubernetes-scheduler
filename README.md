# Custom Kubernetes Scheduler with Predictive Modeling

## Project Overview
This project introduces a custom Kubernetes scheduler that utilizes load history and predictive models to optimize resource scheduling in advance. The aim is to enhance resource utilization efficiency and overall cluster performance by dynamically adjusting scheduling based on predicted loads.

## Features
- **Custom Kubernetes Scheduler:** Implementation of a custom scheduler that integrates with the Kubernetes API.
- **Scheduling Webhooks:** Use of scheduling webhooks to modify scheduling decisions made by the default Kubernetes scheduler.
- **Predictive Modeling:** Application of advanced predictive algorithms to analyze load history and forecast future resource requirements.

## Prerequisites
- A Kubernetes cluster installed and configured.
- `kubectl` installed on your local machine for interacting with the Kubernetes cluster.

## Deploying the Custom Scheduler
Build the custom scheduler image (Docker example):

docker build -t custom-scheduler:latest .

Deploy the custom scheduler to the Kubernetes cluster:
kubectl apply -f deployment.yaml

## Setting Up Scheduling Webhooks
Deploy the webhook server:

kubectl apply -f webhook-deployment.yaml
Configure the webhooks in the Kubernetes cluster:

kubectl apply -f webhook-configuration.yaml
## Usage
After setting up the custom scheduler and scheduling webhooks, the Kubernetes cluster will use the new scheduling logic for new jobs. Monitor the cluster behavior and job scheduling using kubectl:

kubectl get pods -o wide

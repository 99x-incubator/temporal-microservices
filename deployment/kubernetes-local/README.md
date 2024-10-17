
# Temporal kubernetes deployment

**Work in process.**   

Deploying temporal runtime side by side with native workloads using kind as the example
Guide.

 1. (install kubetl and something like openlens)

 2. setup a local cluster eg kind (https://kind.sigs.k8s.io/docs/user/quick-start/#installation) openlens to check the cluster
		 - kind create cluster --name local-01  (once the installation is successful there), you kube context should have been set to the new 
  3. building the microservice images 
  -  docker files 
		  - microservices/admin_gateway/Dockerfile
		  - microservices/admin_notifications/Dockerfile
		  - temporal/workers/Dockerfile
		  - clients/av1-admin/Dockerfile
  - load the local images to the kind cluster 
		  - `kind load docker-image admin_web admin_gateway admin_worker admin_notifications --name local-01`

 3. Deploying services to the cluster 
	 - deployment/kubernetes-local/00_base.yaml
	 - deployment/kubernetes-local/01_temporal.txt
	 - deployment/kubernetes-local/02_temporal_admin_task_queue_worker.yaml
	 - deployment/kubernetes-local/03_admin_gateway.yaml
	 - deployment/kubernetes-local/04_admin_notification.yaml
	 - deployment/kubernetes-local/05_admin_web.yaml



housekeeping

kind delete cluster --name local-01
docker system prune -a



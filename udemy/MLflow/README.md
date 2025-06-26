# [Build scalable MLOps with MLflow, KServe, Docker, and Kubernetes. Automate deployments, monitor models, and workflows.](https://www.udemy.com/course/mlflow-kubernetes-mlops)

```bash
python3 -m venv venv
source venv/bin/activate
pip install 'mlflow[server]'
```


## Kind cluster install

```bash
./kind.sh
```

## [KServe install](https://kserve.github.io/website/latest/get_started/#install-helm)

```bash
curl -s "https://raw.githubusercontent.com/kserve/kserve/release-0.15/hack/quick_install.sh" | bash\n
```

## [Kubernetes dashboard](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/)

```bash
helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard/
helm upgrade --install kubernetes-dashboard kubernetes-dashboard/kubernetes-dashboard --create-namespace --namespace kubernetes-dashboard
 ```


 ## [Deploying Models to Kubernetes with AIStor, MLflow and KServe: ](https://blog.min.io/deploying-models-to-kubernetes-with-aistor-mlflow-and-kserve/)

 ```bash
kind create cluster
kubectl config use-context kind-kind
kubectl create namespace mlflow-kserve-test
```


# Course notes

## [Udemy course 1: https://www.udemy.com/course/mastering-ollama-a-comprehensive-guide](https://www.udemy.com/course/mastering-ollama-a-comprehensive-guide)

## [Udemy course 2: https://www.udemy.com/course/mastering-ollama-a-comprehensive-guide](https://www.udemy.com/course/ollama-starttech)

## ChromaDB Helm chart

```bash
helm repo add chroma https://amikos-tech.github.io/chromadb-chart/
helm repo update
helm install chroma chroma/chromadb --set chromadb.apiVersion="0.5.11"

kubectl get secret chromadb-auth -o jsonpath="{.data.token}" | base64 --decode
# or use this to directly export variable
export CHROMA_TOKEN=$(kubectl get secret chromadb-auth -o jsonpath="{.data.token}" | base64 --decode)

kubectl port-forward svc/chroma-chromadb 8000:8000

## [chromadb-admin repo](https://github.com/flanker/chromadb-admin)

docker run -p 3001:3000 fengzhichao/chromadb-admin
```

- NOTE: Use http://host.docker.internal:8000 for the connection string if you want to connect to a ChromaDB instance running locally.

## Golang client

- https://cookbook.chromadb.dev/core/clients/#http-client


## Ollama commands

- Exiting

```bash
/bye
```

- Verboes mode

```bash
/set verbose
```

- Set/Save system parameters

```bash
/set system "You are Jarvis from Iron-man and the user is the iron man. respond in 2 lines only."
/save jarvisai
```


source $HOME/adk-hello-world-go/set_env.sh

if ss -ltn | grep -q :8081; then
  echo "Proxy is already running at http://127.0.0.1:8081/ui/?app=hello_time_agent"
else
  /usr/bin/gcloud run services proxy adk-hello-world-go --project comglitn --port 8081 --region us-central1
fi

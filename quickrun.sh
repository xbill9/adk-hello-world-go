source $HOME/adk-hello-world-go/set_env.sh


echo `pwd`
echo go run agent.go quick deploy

#--set-env-vars GOOGLE_CLOUD_PROJECT=comglitn,GOOGLE_CLOUD_LOCATION=us-central1,GOOGLE_GENAI_USE_VERTEXAI=true    

/usr/bin/gcloud run deploy adk-hello-world-go --source . --region us-central1 --project comglitn --ingress all --clear-secrets --no-allow-unauthenticated

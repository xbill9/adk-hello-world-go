source $HOME/adk-hello-world-go/set_env.sh


echo `pwd`
echo gcloud run deploy $SERVICE_NAME --source . --region $REGION --project $PROJECT_ID

gcloud run deploy $SERVICE_NAME --source . --region $REGION --project $PROJECT_ID --ingress all --clear-secrets --no-allow-unauthenticated \
--set-env-vars GOOGLE_CLOUD_PROJECT=$PROJECT_ID,GOOGLE_CLOUD_LOCATION=$REGION,GOOGLE_GENAI_USE_VERTEXAI=true

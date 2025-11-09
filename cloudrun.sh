source $HOME/adk-hello-world-go/set_env.sh


echo `pwd`
echo go run agent.go deploy

#--set-env-vars GOOGLE_CLOUD_PROJECT=comglitn,GOOGLE_CLOUD_LOCATION=us-central1,GOOGLE_GENAI_USE_VERTEXAI=true    

$HOME/adk-go/adkgo deploy cloudrun \
    -p $GOOGLE_CLOUD_PROJECT \
    -r $GOOGLE_CLOUD_LOCATION \
    -s $SERVICE_NAME \
    -e $AGENT_PATH 

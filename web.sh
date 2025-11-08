source $HOME/adk-hello-world-go/set_env.sh

cd hello-agent

echo `pwd`
echo go run agent.go web api webui
go run agent.go web api webui

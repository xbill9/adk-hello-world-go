source $HOME/adk-hello-world-go/set_env.sh

cd a2a-client-go

echo `pwd`
echo go run agent.go web api webui
go run main.go web api webui

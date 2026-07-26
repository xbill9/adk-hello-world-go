source $HOME/adk-hello-world-go/set_env.sh

cd a2a-client-go

echo `pwd`
# Any arguments are joined into the user prompt; default prompt is used when none are given.
echo go run main.go "$@"
go run main.go "$@"

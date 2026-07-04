package client

//go:generate sh -c "curl -s -u ${USERNAME:-opencode}:${PASSWORD:-pass1} http://localhost:${HOST_PORT:-4002}/doc > schema.json && sh $PWD/../../.scripts/fix-schema.sh schema.json && go run github.com/ogen-go/ogen/cmd/ogen@latest --config $PWD/../../.build/ogen.yml --target . --clean --package client schema.json && rm schema.json"

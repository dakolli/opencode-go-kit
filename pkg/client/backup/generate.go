package client

//go:generate sh -c "curl -s -u opencode:pass1 http://localhost:4002/doc > schema.json && sh $PWD/../../config/fix-schema.sh schema.json && go run github.com/ogen-go/ogen/cmd/ogen@latest --config $PWD/../../config/ogen.yml --target . --clean --package client schema.json && rm schema.json"

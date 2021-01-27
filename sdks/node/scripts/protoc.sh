protoc \
  --proto_path=../../proto \
  --plugin="protoc-gen-ts=node_modules/.bin/protoc-gen-ts" \
  --plugin="protoc-gen-grpc=node_modules/.bin/grpc_tools_node_protoc_plugin" \
  --js_out="import_style=commonjs,binary:./dist/proto" \
  --grpc_out="grpc_js:./dist/proto" \
  --ts_out="service=grpc-node,mode=grpc-js:./dist/proto" \
  ../../proto/ingress.proto

protoc \
  --proto_path=../../proto \
  --plugin="protoc-gen-ts=node_modules/.bin/protoc-gen-ts" \
  --js_out="import_style=commonjs,binary:./dist/proto" \
  --ts_out="./dist/proto" \
  ../../proto/manifest.proto \

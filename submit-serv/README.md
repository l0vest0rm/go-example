
mkdir submit
thrift -o submit --gen js submit.thrift
thrift -o submit --gen go submit.thrift
protoc submit.proto --go_out=plugins=grpc:submit
protoc submit.proto --js_out=import_style=commonjs:submit
protoc submit.proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:submit

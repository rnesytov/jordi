.PHONY: compile_spec
compile_spec:
	@protoc --go_out=. --go-grpc_out=. test_spec.proto

.PHONY: run
run:
	@go run . $(ARGS)

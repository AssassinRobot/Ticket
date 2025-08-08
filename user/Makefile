PROTO_DIR=../proto
OUT_DIR=./gen

generate:
	@protoc -I $(PROTO_DIR) \
		--go_out=$(OUT_DIR) \
		--go_opt=paths=source_relative \
		$(PROTO_DIR)/user.proto \
		$(PROTO_DIR)/user_req_rep.proto \

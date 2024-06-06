GREEN="\033[00;32m"
RESTORE="\033[0m"

# make the output of the message appear green
define style_calls
	$(eval $@_msg = $(1))
	echo ${GREEN}${$@_msg}
	echo ${RESTORE}
endef

.PHONY: test format

test:
	@$(call style_calls, "Testing...")
	@go test -v ./...

format:
	@$(call style_calls, "Formatting...")
	@gofmt -s -w .

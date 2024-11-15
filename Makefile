dev:
	@./run_dev_server.sh
	

test:
	go test ./store ./handlers

int_test:
	hurl --test ./integrationtests/test.hurl

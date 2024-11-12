start_pg:
	dc start

dev:
	air

test:
	go test ./store ./handlers

int_test:
	hurl --test ./integrationtests/test.hurl

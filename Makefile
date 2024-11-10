run_templ:
	templ generate --watch --proxy="http://localhost:3000" --open-browser=false -proxyport=7000

run_server:
	air

live:
	make -j 2 run_templ run_server

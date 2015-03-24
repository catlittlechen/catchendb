all: compile

compile:
	./githash.sh
	cd server; go build . && install ./server ../../bin/ccdb-server && rm ./server && cd ../../bin && tar zcvf ./server.tar.gz ./ccdb-server;
	cd client; go build . && install ./client ../../bin/ccdb-cli && rm ./client && cd ../../bin && tar zcvf ./client.tar.gz ./ccdb-cli;
			        

clean:
	cd server; go clean; cd ../logic; go clean; cd ../client; go clean;

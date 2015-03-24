all: compile

compile:
	./githash.sh
	cd server; go build . && install ./server ../../bin/ccdb && rm ./server && cd ../../bin && tar zcvf ./ccdb.tar.gz ./ccdb;
			        

clean:
	cd server; go clean; cd ../logic; go clean;

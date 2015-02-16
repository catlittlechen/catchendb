all: compile

compile:
	./githash.sh
	cd main; go build . && install ./main ../../bin/ccdb && rm ./main && cd ../../bin && tar zcvf ./ccdb.tar.gz ./ccdb;
			        

clean:
	cd main; go clean; cd ../logic; go clean;

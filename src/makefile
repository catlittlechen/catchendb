all: compile

compile:
	./githash.sh
	cd main; go build . && install ./main ../../bin && rm ./main && cd ../../bin && tar zcvf ./main.tar.gz ./main;
			        

clean:
	cd main; go clean; cd ../logic; go clean;

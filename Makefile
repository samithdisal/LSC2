all: clean
	go build .

clean:
	rm -rf lew

dep: all
	scp lew l1r:/home/rout/st

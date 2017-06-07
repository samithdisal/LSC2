all: clean
	go build .

clean:
	rm -rf lew

dep: all
	scp lew mle:/home/ubuntu/st

test: all
	./lew -s
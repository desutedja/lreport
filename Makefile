build:
	CGO_ENABLED=0 go build -o lreport .

build-container: build
	docker build -t lreport/lreport:dev .

run-container:
	docker run -d -p 1111:1111 lreport/lreport:dev
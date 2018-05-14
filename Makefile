export GOPATH :=/Users/rock/Documents/gitlab_work/ant
export PATH := ${PATH}:${GOPATH}/bin
export GOBIN := ${GOPATH}/bin
main:
	go run main.go
build:
	go install main.go
images:
	docker build -t ant .
run:
	docker run -it -v /Users/rock/Documents/gitlab_work/ant:/gopath --rm ant
docker: images run


go-get:
	# go get github.com/henrylee2cn/pholcus_lib
	# go get github.com/henrylee2cn/pholcus
	# go get github.com/henrylee2cn/goutil
	# go get github.com/gocolly/colly
	# go get github.com/golang/protobuf
	go get github.com/PuerkitoBio/goquery
	go get github.com/astaxie/beego


spider-build:
	cd src/github.com/henrylee2cn/pholcus && go install

run-censor:
	go run src/censor/main.go
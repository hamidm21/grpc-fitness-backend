#!make
#TODO must handle trailing whitespace in include .env

version := $(shell git describe --abbrev=0 --tags)


e2e:
	export $$(cat .env | grep -v ^\# | xargs) && \
		go run cmd/e2e/main.go
grpcurl-remote:
	grpcurl -proto vendor/gitlab.com/mefit/mefit-api/proto/mefit.proto grpc.fitex.app.yottab.io:443 $(cmd)
grpcurl-local:
	grpcurl -H "$(h)" -d @ -insecure -proto vendor/gitlab.com/mefit/mefit-api/proto/mefit.proto 127.0.0.1:8443 $(cmd)
# Samples:
# make grpcurl cmd=Mefit/PhoneSignIn <<EOF       
# {
# "phoneNo": "9376866378"
# }
# EOF
# 
# Login
# make grpcurl cmd=Mefit/Login <<EOF                     master
# {
# "phoneNo": "9376866378",
# "activationToken": "123456"
# }
# EOF
#
# ProfileUpdate
# make grpcurl h="x-token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjYzMDMyOTEsImlhdCI6MTU0MDM4MzI5MSwic3ViIjoxfQ.u1xzx1diCXAxPMrgEeP53s24MbFZikQFoEXLMnHfKxc" cmd=Mefit/ProfileUpdate <<EOF
# {
# "name": "homayoun",
#  "age": 25
# }
# EOF
#
# Profile info
# grpcurl -H "x-token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjYzMDMyOTEsImlhdCI6MTU0MDM4MzI5MSwic3ViIjoxfQ.u1xzx1diCXAxPMrgEeP53s24MbFZikQFoEXLMnHfKxc" -insecure -proto vendor/gitlab.com/mefit/mefit-api/proto/mefit.proto 127.0.0.1:8443 Mefit/ProfileInfo

all: build

build-payment:
	go build -o build/payment cmd/payment/main.go

build-admin:
	go build -o build/admin cmd/admin/main.go

build-server:
	go build -o build/server cmd/grpc/main.go

git-config:
	git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"

clean:
	rm -rf build/*

docker-build-payment:
	docker build -t hub.yottab.io/fitex/fitex-payment:$(version) -f payment.Dockerfile .

docker-build-admin:
	docker build -t hub.yottab.io/fitex/fitex-admin:$(version) -f admin.Dockerfile .

docker-build-server:
	docker build -t hub.yottab.io/fitex/fitex-server:$(version) .

docker-push-payment:
	docker push hub.yottab.io/fitex/fitex-payment:$(version)

docker-push-admin:
	docker push hub.yottab.io/fitex/fitex-admin:$(version)

docker-push-server:
	docker push hub.yottab.io/fitex/fitex-server:$(version)

run-payment: build-payment
	export $$(cat .env | grep -v ^\# | xargs) && \
		./build/payment

run-admin: build-admin
	export $$(cat .env | grep -v ^\# | xargs) && \
		./build/admin

run-server: build-server
	export $$(cat .env | grep -v ^\# | xargs) && \
		./build/server
		
# To update the server content, we must change the .env respectfuly
manager: 
	export $$(cat .env | grep -v ^\# | xargs) && \
		go run cmd/manager/main.go

dep-ensure-tor:
	https_proxy=socks5://127.0.0.1:9150 http_proxy=socks5://127.0.0.1:9150 dep ensure -v

test:
	go test -v controller/...

watch:
	while inotifywait -e modify -r .; do killall server && make run &; done
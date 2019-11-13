FROM golang:1.11-stretch AS build

ADD ./ /go/src/2019_2_Next_Level

WORKDIR /go/src/2019_2_Next_Level


RUN go get -v ./...
RUN go build -o post ./cmd/post/main.go

ENTRYPOINT ["./post"]

# CMD go run post
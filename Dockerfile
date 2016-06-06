FROM golang:1.6

COPY . /go/src/github.com/kyleterry/tweetsave
WORKDIR /go/src/github.com/kyleterry/tweetsave
RUN go build -v
RUN go install -v

RUN tweetsave

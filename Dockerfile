FROM golang:1.6

ADD . /go/src/github.com/kyleterry/tweetsave
RUN go install /go/src/github.com/kyleterry/tweetsave
ENTRYPOINT /go/bin/tweetsave

EXPOSE 8080

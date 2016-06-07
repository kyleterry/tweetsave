FROM golang:1.6

COPY . /go/src/github.com/kyleterry/tweetsave
WORKDIR /go/src/github.com/kyleterry/tweetsave
RUN go install -v

EXPOSE 8080

CMD ["sh", "-c", "tweetsave -db-bind=${DATABASE_URL} -api-bind=:8080 -consumer-key=${TWEETSAVE_CONSUMER_KEY} -consumer-secret=${TWEETSAVE_CONSUMER_SECRET} -access-token=${TWEETSAVE_ACCESS_TOKEN} -access-secret=${TWEETSAVE_ACCESS_SECRET}"]

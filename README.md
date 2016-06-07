# Tweet Save

## Building

`go build`

## Running Locally

Ensure you have postgres running and a tweetsave database created with the
correct permissions. Look up instructions for your OS on how to do this.

```bash
./tweetsave \
 -consumer-key <consumerkey> \
 -consumer-secret <consumersecret> \
 -access-token <accesstoken> \
 -access-secret <accesssecret>
```

## Installing locally

`go install`

## Deploying to Heroku

```bash
# using heroku toolbelt
heroku create
heroku addons:create heroku-postgresql:hobby-dev
heroku config:set TWEETSAVE_CONSUMER_KEY=<YOUR KEY>
heroku config:set TWEETSAVE_CONSUMER_SECRET=<YOUR KEY>
heroku config:set TWEETSAVE_ACCESS_TOKEN=<YOUR KEY>
heroku config:set TWEETSAVE_ACCESS_SECRET=<YOUR KEY>
git push heroku master
```

## Running in Docker

```bash
docker build -t tweetsave
docker run -it --rm --name tweetsave \
-e TWEETSAVE_CONSUMER_KEY=<your twitter consumer key> \
-e TWEETSAVE_CONSUMER_SECRET=<your twitter consumer secret> \
-e TWEETSAVE_ACCESS_TOKEN=<your access token> \
-e TWEETSAVE_ACCESS_SECRET=<your access secret> \
-e DATABASE_URL=postgres://<your remote db>/tweetsave \
tweetsave
```

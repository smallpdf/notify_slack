# notify\_slack

Turns your failing `go test` into slack messages

## installation
```sh
wget -nv https://github.com/smallpdf/notify_slack/releases/latest/download/notify_slack && chmod +x notify_slack
```

## setup
environment variables
```sh
NOTIFY_LOGLEVEL=DEBUG|INFO*|ERROR|...
NOTIFY_LOGFORMAT=console|json*
NOTIFY_DRYRUN=true|false*
NOTIFY_USERS=user@smallpdf.com,...
NOTIFY_GROUPS=backend,...
NOTIFY_SLACKTOKEN=....
```

`NOTIFY_USERS` and `NOTIFY_GROUPS` accept lists with a `,` as a separator.

## example
```
go test -json ./testdata | tee notify_slack
go test -json ./testdata | tee >(NOTIFY_USERS=user@smallpdf.com notify_slack)
```
Don't forget to add `Tester` to the Account or Group.

## releases
```
git tag v0.0.1
git push --tags
#see github actions
```

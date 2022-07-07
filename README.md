# notify\_slack

Turns your failing `go test` into slack messages

## setup
environment variables
```sh
NOTIFY_LOGLEVEL=DEBUG|INFO*|ERROR|...
NOTIFY_LOGFORMAT=console|json*
NOTIFY_DRYRUN=true|false*
NOTIFY_USERS=user@smallpdf.com,...
NOTIFY_GROUPS=backend,...
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

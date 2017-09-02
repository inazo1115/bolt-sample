# bolt-sample

## Usage example

Server
```
$ go run main.go
```

Client
```
$ curl 'http://localhost:8080/put?key=foo&value=bar'
ok
$ curl 'http://localhost:8080/put?key=foofoo&value=barbar'
ok
$ curl 'http://localhost:8080/put?key=bar&value=baz'
ok
$ curl 'http://localhost:8080/get?key=foo'
key=foo, value=bar
$ curl 'http://localhost:8080/list'
key=bar, value=baz
key=foo, value=bar
key=foofoo, value=barbar
$ curl 'http://localhost:8080/list?prefix=foo'
key=foo, value=bar
key=foofoo, value=barbar
$ curl 'http://localhost:8080/backup'

$N¾-

    UM
      _RMyBucket foobarfoofoobarbahMyBucket0&barbazfoobarfoofoobarbar
```

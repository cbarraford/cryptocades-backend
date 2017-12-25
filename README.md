[![CircleCI](https://circleci.com/gh/CBarraford/lotto/tree/master.svg?style=svg&circle-token=e41380660a1b6ecd373ffc742a8f6df7cd821bcb)]
[![codecov](https://codecov.io/gh/CBarraford/lotto/branch/master/graph/badge.svg?token=e1O9Ww2XUC)](https://codecov.io/gh/CBarraford/lotto)

Lotto
=====

## Development
This project adheres to the [golang
standard](https://golang.org/doc/code.html#Organization) of GOPATH=$HOME. So make sure you
clone this project in ~/src/github.com/CBarraford/lotto

### Create schema migration
```
make get # downloads the migrate binary needed
make create-migration name="another_one"
```

### Run app
This starts the achiever service with postgres. Be aware, when stopped,
postgres continues to run. If you want to stop it `docker-compose down` will
do it.
```
make run
```

### Run tests
With integration tests (ie db)
```
make test
```

Without integration tests
```
make test-short
```

### Linting
```
make lint
```

### Masquerade API
You can make API requests masquerading as any user (enabled by the config,
false by default).

```
curl localhost:8080/users/<id>
{"error":"unauthorized"}

curl -H "Masquerade: <id>" localhost:8080/users/<id>
Success!
```

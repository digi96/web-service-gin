# Sqlc

# Swaggo

## Install

```console
go install github.com/swaggo/swag/cmd/swag@latest
export PATH=$PATH:/Users/{username}/bin
source ~/.zshrc
echo 'cd to go project folder'
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

## Annotation

###### example: https://santoshk.dev/posts/2022/how-to-integrate-swagger-ui-in-go-backend-gin-edition/

# Testify

## Install

```console
go get github.com/stretchr/testify
```

## Design Test Function

```code
package yours

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {

  assert.True(t, true, "True is true!")

}
```

## Run Test

```console
go test -v -cover ./...
```

# Rabbitmq

## Reference

###### https://medium.com/@celalsahinaltinisik/rabbitmq-publish-consume-with-golang-26412784b6e5

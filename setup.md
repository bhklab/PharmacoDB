#### Some setup steps followed:

**Install go**

**Export GOPATH.**

```
export GOPATH=~/goland
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$GOPATH/bin
```
(Don't forget to create the folder `~/goland`, which will be the go workspace)

**Install dependencies:**

```
go get "github.com/go-sql-driver/mysql"
go get "gopkg.in/gin-gonic/gin.v1"
```

**export the following environment variables** (replace `some_name` and `some_password` with appropriate values)

```
export pharmacodb_api_dbname="some_name"
export local_mysql_passwd="some_password"
```

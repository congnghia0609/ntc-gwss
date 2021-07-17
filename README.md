# ntc-gwss
ntc-gwss is golang code example websocket server using library gorilla.  

## Install dependencies
```bash
# Install dependencies
#make deps
go mod download

# update go.mod file
go mod tidy
```

## Build
```bash
export GO111MODULE=on
make build
```

## Clean file build
```bash
make clean
```

## Run with environment: development | test | staging | production
### development
```bash
make run
```
### test
```bash
make run-test
```
### staging
```bash
make run-stag
```
### production
```bash
make run-prod
```


## License
This code is under the [Apache License v2](https://www.apache.org/licenses/LICENSE-2.0).  

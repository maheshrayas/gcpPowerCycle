# Cost saving utility for google compute engine and google kubernetes engine (stateless containers)
![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)

## Initialize common
``` go mod init github.com/maheshrayas/powerCycle/common ```

## Initialize functions
``` go mod init  github.com/maheshrayas/powerCycle ```
Update go.mod file in function folder
``` require github.com/maheshrayas/powerCycle/common v0.0.0 ```
``` replace github.com/maheshrayas/powerCycle/common => ../common ```

## Initilize vendor in function folder
``` go mod vendor ```

## Deploy in cloud functions
``` cd functions ```

 ``` gcloud functions deploy PowerCycle --runtime go111 --trigger-http ```


## Tagging standard:

start_08-00_mon-fri_stop_16-30 : Keep the instance up and running from Monday to Friday between 08:00 to 16:30

## Known issue with go mod vendor
go mod vendor
go: modules disabled inside GOPATH/src by GO111MODULE=auto; see 'go help modules'
export GO111MODULE=on and then run
go mod vendor

## TODO:
1. Run the cloudfunctions at the timely interval (every 30 mins)
2. Documentation
3. Unit test

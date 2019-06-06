# Initialize common
``` go mod init github.com/maheshrayas/powerCycle/common ```

# Initialize functions
``` go mod init  github.com/maheshrayas/powerCycle ```
Update go.mod file in function folder
``` require github.com/maheshrayas/powerCycle/common v0.0.0 ```
``` replace github.com/maheshrayas/powerCycle/common => ../common ```

# Initilize vendor in function folder
``` go mod vendor ```

# Deploy in cloud functions
``` cd functions ```
 ``` gcloud functions deploy PowerCycle --runtime go111 --trigger-http ```

 #TODO
 * Stop the instance
 * Take the snapshot of the disk
 * Edit the instance and delete the dist associated with it
 * Save the instance
 * Delete the disk

* Create a disk using the snapshot
* Edit the instance and attach this disk
* start the instance

* Issue with go mod vendor
go mod vendor
go: modules disabled inside GOPATH/src by GO111MODULE=auto; see 'go help modules'
export GO111MODULE=on and then run
go mod vendor




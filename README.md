# freedrive
free storage for all

### setup
- clone the repo: `git clone https://github.com/brandoneprice31/freedrive.git`
- install dependencies: `dep ensure`
- compile the code: `make build`
- make a directory for storing freedrive keys: `mkdir /path/to/freedrive/dir`

### usage
- backup: `FD=/path/to/freedrive/dir BP=/path/to/backup/dir make backup`
- download: `FD=/path/to/freedrive/dir DL=/path/to/download/dir make download`

### creating a new service
- fork this repo
- create a struct that fulfills the `Service` interface defined in `service/service.go`
- look at `service/braintree.go` for an example
- create a pull request

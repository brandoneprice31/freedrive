# freedrive
free storage for all

### setup
- clone the repo: `git clone https://github.com/brandoneprice31/freedrive.git`
- install dependencies: `dep ensure`
- compile the code: `make build`
- make a directory for storing freedrive keys: `mkdir /path/to/freedrive/dir`

### usage
- setup the `config.yaml` file with your credentials
- backup: `make backup`
- download: `make download`

### creating a new service
- fork this repo
- create a struct that fulfills the `Service` interface defined in `service/service.go`
- look at `service/braintree.go` for an example
- create a pull request

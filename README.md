# freedrive
free storage for all

### what is it
- a free online storage hack inspired by Jamie Davies' [Medium article](https://medium.com/@viralpickaxe/how-we-hacked-the-braintree-api-to-store-an-unlimited-number-of-files-302860736c25) on how he hacked the Braintree API to store unlimited data
- freedrive uploads / downloads entire directories to multiple different "services" like the braintree sandbox, twitter, dropbox, and hopefully more to come
- see my freedrive [twitter account](https://twitter.com/freedrivetest) to see my encrypted files

### setup
- clone the repo: `git clone https://github.com/brandoneprice31/freedrive.git`
- install dependencies: `dep ensure`
- compile the code: `make build`
- make a directory for storing freedrive keys: `mkdir /path/to/freedrive/dir`

### usage
- setup the `config.yaml` file with your credentials (feel free to use my credentials although its easy to set up your own twitter api keys / braintree sandbox, etc.)
- backup: `make backup`
- download: `make download`

### creating a new service
- fork this repo
- create a struct that fulfills the `Service` interface defined in `service/service.go`
- look at `service/braintree.go` for an example
- create a pull request

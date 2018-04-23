BACKUP=bin/backup
DOWNLOAD=bin/download

ifeq ($(TRAVIS), true)
	CGO_ENABLED := 0
else
	CGO_ENABLED := 1
endif

build:
	CGO_ENABLED=${CGO_ENABLED} go build -o ${BACKUP} cmd/backup/*.go
	CGO_ENABLED=${CGO_ENABLED} go build -o ${DOWNLOAD} cmd/download/*.go

backup:
	${BACKUP}

download:
	${DOWNLOAD}

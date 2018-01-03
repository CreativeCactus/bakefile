default:
	go build bake.go

test:
	echo ${1:-test}
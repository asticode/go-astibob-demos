i:
	go run index/main.go

w1:
	go run workers/1/main.go

w2:
	go run workers/2/main.go

w3:
	CGO_CXXFLAGS="-I${CURDIR}/tmp/include" \
	LIBRARY_PATH=${CURDIR}/tmp/lib:${LIBRARY_PATH} \
	LD_LIBRARY_PATH=${CURDIR}/tmp/lib:${LD_LIBRARY_PATH} \
	go run workers/3/main.go
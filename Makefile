linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -buildmode=c-shared -o gophpfetch.so main.go \
	 && gcc -E -P gophpfetch.h -o gophpfetchx.h \
	  && echo '#define FFI_LIB "./gophpfetch.so"'>gophpfetch.h \
	  && sed -i -e 's/extern size_t _GoStringLen(_GoString_ s);//g' \
	   -e 's/extern const char \*_GoStringPtr(_GoString_ s);//g' \
	    -e 's/typedef float _Complex GoComplex64;//g' \
	     -e 's/typedef double _Complex GoComplex128;//g' gophpfetchx.h  \
	   &&  cat gophpfetchx.h >> gophpfetch.h && rm gophpfetchx.h
macos:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -buildmode=c-shared -o gophpfetch.so main.go
windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -buildmode=c-shared main.go
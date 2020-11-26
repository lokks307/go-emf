module github.com/lokks307/go-emf

go 1.15

require (
	github.com/mattn/go-colorable v0.1.8
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68
)

replace (
	github.com/lokks307/go-emf/emf => ./emf
	github.com/lokks307/go-emf/w32 => ./w32
)

module github.com/lokks307-dev/go-emf

go 1.15

require (
	github.com/disintegration/imaging v1.6.2
	github.com/mattn/go-colorable v0.1.8
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/sys v0.0.0-20201119102817-f84b799fce68
)

replace (
	github.com/lokks307-dev/go-emf/emf => ./emf
	github.com/lokks307-dev/go-emf/w32 => ./w32
)

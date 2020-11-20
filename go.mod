module github.com/lokks307/go-emf

go 1.15

require (
	github.com/disintegration/imaging v1.6.2
	github.com/mattn/go-colorable v0.1.8
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/image v0.0.0-20200927104501-e162460cd6b5 // indirect
	golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae
	golang.org/x/text v0.3.4
)

replace (
	github.com/lokks307/go-emf/emf => ./emf
	github.com/lokks307/go-emf/fontname => ./fontname
	github.com/lokks307/go-emf/w32 => ./w32
)

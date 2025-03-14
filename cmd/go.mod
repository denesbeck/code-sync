module github.com/denesbeck/code-sync

go 1.23.0

toolchain go1.24.1

require cli v0.0.0-00010101000000-000000000000

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	golang.org/x/sys v0.31.0 // indirect
)

replace cli => ./cli

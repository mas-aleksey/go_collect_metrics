module main

go 1.19

replace utils => ../../internal/utils
replace clients => ../../internal/clients

require (
	clients v0.0.0-00010101000000-000000000000 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	utils v0.0.0-00010101000000-000000000000 // indirect
)

module main

go 1.19

replace utils => ../../internal/utils
replace handlers => ../../internal/handlers
replace storage => ../../internal/storage

require (
	handlers v0.0.0-00010101000000-000000000000 // indirect
	storage v0.0.0-00010101000000-000000000000 // indirect
	utils v0.0.0-00010101000000-000000000000 // indirect
)

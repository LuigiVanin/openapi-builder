module github.com/LuigiVanin/openapi-builder/cmd

go 1.25.3

require (
	github.com/LuigiVanin/openapi-builder v0.0.0
	github.com/flowchartsman/swaggerui v0.0.0-20221017034628-909ed4f3701b
)

require github.com/goccy/go-yaml v1.19.2 // indirect

// Build the demos against the local library source instead of a published version.
replace github.com/LuigiVanin/openapi-builder => ../

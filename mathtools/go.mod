module github.com/clinaresl/mathprob/mathtools

go 1.16

replace github.com/clinaresl/mathprob/components => ../components

replace github.com/clinaresl/mathprob/fstools => ../fstools

require (
	github.com/clinaresl/mathprob/components v0.0.0-00010101000000-000000000000
	github.com/clinaresl/mathprob/fstools v0.0.0-00010101000000-000000000000
	github.com/clinaresl/mathprob/helpers v0.0.0-00010101000000-000000000000
)

replace github.com/clinaresl/mathprob/helpers => ../helpers

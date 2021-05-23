module github.com/clinaresl/mathprob/mathtools

go 1.16

replace github.com/clinaresl/mathprob/components => ../components

replace github.com/clinaresl/mathprob/fstools => ../fstools

require (
	github.com/clinaresl/mathprob/components v0.0.0-20210523185513-4af87bf01910
	github.com/clinaresl/mathprob/fstools v0.0.0-20210523185513-4af87bf01910
	github.com/clinaresl/mathprob/helpers v0.0.0-20210523200000-5d4315b08dbd
)

replace github.com/clinaresl/mathprob/helpers => ../helpers

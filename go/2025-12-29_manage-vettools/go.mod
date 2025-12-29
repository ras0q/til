module manage-vettools

go 1.25.5

require (
	golang.org/x/exp/typeparams v0.0.0-20231108232855-2478ac86f678 // indirect
	golang.org/x/mod v0.31.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/tools v0.40.0 // indirect
	honnef.co/go/tools v0.6.1 // indirect
	vettools v0.0.0-00010101000000-000000000000 // indirect
)

replace vettools => ./vettools/

tool vettools

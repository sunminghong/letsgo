go test -bench=.* -cpuprofile=cpu.prof
go test -bench=.* -cpuprofile=cpu.prof -c
go tool pprof net.test cpu.prof

echo "--> running go benchmarks..."
go test . -bench . -run none | sed 's/^/    /'
echo "--> running gopherjs benchmarks..."
gopherjs test github.com/albrow/vdom --bench=. --run=none | sed 's/^/    /'
echo "    compiling karma benchmarks to js..."
gopherjs build karma/go/bench.go -o karma/js/bench.js | sed 's/^/    /'
echo "    running benchmarks with karma..."
karma run | sed 's/^/    /'
echo "DONE."

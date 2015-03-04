echo "--> running go tests..."
go test . | sed 's/^/    /'
echo "--> running gopherjs tests..."
gopherjs test github.com/albrow/vdom | sed 's/^/    /'
echo "--> running karma tests..."
echo "    compiling karma tests to js..."
gopherjs build karma/go/test.go -o karma/js/test.js | sed 's/^/    /'
echo "    running tests with karma..."
karma run | sed 's/^/    /'
echo "DONE."

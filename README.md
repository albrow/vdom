vdom
====

[![GoDoc](https://godoc.org/github.com/albrow/vdom?status.svg)](https://godoc.org/github.com/albrow/vdom)

vdom is a virtual dom implementation written in go which is compatible with
[gopherjs](http://www.gopherjs.org/) and inspired by
[react.js](http://facebook.github.io/react/). The primary purpose of
vdom is to improve the performance of view rendering in
[humble](https://github.com/soroushjp/humble), a framework that lets you write
frontend web apps in pure go and compile them to js to be run in the browser.
However, vdom is framework agnostic, and generally will work whenever you can
render html for your views as a slice of bytes.


Development Status
------------------

vdom is no longer actively maintained. Additionally, ad hoc testing suggests that it might currently
be slower than `setInnerHTML` in at least some cases.

Most users today will likely want to use WebAssembly instead of GopherJS.

Still, this repo might be a decent starting point for anyone wishing to create a virtual DOM implementation
in Go.

Browser Compatibility
---------------------

vdom has been tested and works with IE9+ and the latest versions of Chrome, Safari, and Firefox.

Javascript code generated with gopherjs uses typed arrays, so in order to work with IE9, you will
need a polyfill. There is one in karma/js/support/polyfill/typedarray.js which is used for the
karma tests.


Installing
----------

Assuming you have already installed go and set up your go workspace, you can install
vdom like you would any other package:

`go get github.com/albrow/vdom`

Import in your go code:

`import "github.com/albrow/vdom"`

Install the latest version of [gopherjs](https://github.com/gopherjs/gopherjs), which
compiles your go code to javascript:

`go get -u github.com/gopherjs/gopherjs`

When you are ready, compile your go code to javascript using the `gopherjs` command line
tool. Then include the resulting js file in your application.


Quickstart Guide
----------------

I'll update this section when all functionality is completed. For now, here's a preview
of what usage will probably look like.

Assuming you have a go html template called todo.tmpl:

```html
<li class="todo-list-item {{ if .Completed }} completed {{ end }}">
	<input class="toggle" type="checkbox" {{ if .Completed }} checked {{ end }}>
	<label class="todo-label">{{ .Title }}</label>
	<button class="destroy"></button>
	<input class="edit" onfocus="this.value = this.value;" value="{{ .Title }}">
</li>
```

And a Todo view/model type that looks like this:

```go
type Todo struct {
	Title     string
	Completed bool
	Root    dom.Element
	tree      *vdom.Tree
}
```

You could do the following:

```go
var todoTmpl = template.Must(template.ParseFiles("todo.tmpl"))

func (todo *Todo) Render() error {
	// Execute the template with the given todo and write to a buffer
	buf := bytes.NewBuffer([]byte{})
	if err := tmpl.Execute(buf, todo); err != nil {
		return err
	}
	// Parse the resulting html into a virtual tree
	newTree, err := vdom.Parse(buf.Bytes())
	if err != nil {
		return err
	}
	// Calculate the diff between this render and the last render
	patches, err := vdom.Diff(todo.tree, newTree)
	if err != nil {
		return err
	}
	// Effeciently apply changes to the actual DOM
	if err := patches.Patch(todo.Root); err != nil {
		return err
	}
	// Remember the virtual DOM state for the next render to diff against
	todo.tree = newTree
}
```

Testing
-------

vdom uses three sets of tests. If you're on a unix system, you can run all the tests
in one go with `scripts/test.sh`. The script also compiles the go files to javsacript each
time it runs. You will still need to install the dependencies for the script to work correctly.

### Go Tests

Traditional go tests can be run with `go test .`. These tests are for code which does not
interact with the DOM or depend on js-specific features.

### Gopherjs Tests

You can run `gopherjs test github.com/albrow/vdom` to compile the same tests from above
to javascript and tests them with node.js. This will also test some code which might depend
on js-specific features (but not the DOM) and can't be tested with pure go. You will need
to install [node.js](http://nodejs.org/) to run these tests.

### Karma Tests

vdom uses karma and the jasmine test framework to test code that interacts with the DOM in
real browsers. You will need to install these dependencies:

- [node.js](http://nodejs.org/)
- [karma](http://karma-runner.github.io/0.12/index.html)
- [karma-jasmine](https://github.com/karma-runner/karma-jasmine)

Don't forget to also install the karma command line tools with `npm install -g karma-cli`.

You will also need to install a launcher for each browser you want to test with, as well as the
browsers themselves. Typically you install a karma launcher with `npm install -g karma-chrome-launcher`.
You can edit the config files `karma/test-mac.conf.js` and `karma/test-windows.conf.js` if you want
to change the browsers that are tested on. The Mac OS config specifies Chrome, Firefox, and Safari, and
the Windows config specifies IE9-11, Chrome, and Firefox. You only need to install IE11, since the 
older versions can be tested via emulation.

Once you have installed all the dependencies, start karma with `karma start karma/test-mac.conf.js` or 
`karma start karma/test-windows.conf.js` depending on your operating system. If you are using a unix
machine, simply copy one of the config files and edit the browsers section as needed. Once karma is
running, you can keep it running in between tests.

Next you need to compile the test.go file to javascript so it can run in the browsers:

```
gopherjs build karma/go/test.go -o karma/js/test.js
```

Finally run the tests:

```
karma run
```

Benchmarking
------------

vdom uses three sets of benchmarks. If you're on a unix system, you can run all the benchmarks
in one go with `scripts/bench.sh`. The script also compiles the go files to javsacript each
time it runs. You will still need to install the dependencies for the script to work correctly.

**NOTE:** There are some additional dependencies for benchmarking that are not needed for testing.

### Go Benchmarks

Traditional go benchmarks can be run with `go test -bench . -run none`. I don't expect you
to be using vdom in a pure go context (but there's nothing stopping you from doing so!), so
these tests mainly serve as a comparison to the gopherjs benchmarks. It also helps with
catching obvious performance problems early.

### Gopherjs Benchmarks

To compile the library to javascript and benchmark it with node.js, you can run
`gopherjs test github.com/albrow/vdom --bench=. --run=none`. These benchmarks are only
for code that doesn't interact directly with the DOM. You will need to install
[node.js](http://nodejs.org/) to run these benchmarks.

### Karma Benchmarks

vdom uses karma and benchmark.js to test code that interacts with the DOM in real browsers.
You will need to install these dependencies:

- [node.js](http://nodejs.org/)
- [karma](http://karma-runner.github.io/0.12/index.html)
- [karma-benchmark](https://github.com/JamieMason/karma-benchmark)

Don't forget to also install the karma command line tools with `npm install -g karma-cli`.

Just like with the tests, you will need to install a launcher for each browser you want to test with.

Once you have installed all the dependencies, start karma with `karma start karma/bench-mac.conf.js` or 
`karma start karma/bench-windows.conf.js` depending on your operating system. We have to use
different config files because of a [limitation of karma-benchmark](https://github.com/JamieMason/karma-benchmark/issues/7).
You will probably want to kill karma and restart it if you were running it with the test configuration.
If you are using a unix machine, simply copy one of the config files and edit the browsers section as
needed. Once karma is running, you can keep it running in between benchmarks.

Next you need to compile the bench.go file to javascript so it can run in the browsers:

```
gopherjs build karma/go/bench.go -o karma/js/bench.js
```

Finally run the benchmarks:

```
karma run
```

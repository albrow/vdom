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


Progress
--------

vdom is now very close to being feature-complete, and it is pretty rigourously tested. Version
0.0.1 will be released soon. All that's left to do is benchmark and, if necessary, improve
performance. Ad hoc testing suggests that it might currently be slower than `setInnerHTML` in
at least some cases. I plan to fix this if possible.


Browser Compatibility
---------------------

vdom is tested with and works with the latest versions of Chrome, Safari, Firefox, and
Internet Explorer. For windows systems, vdom works only with IE10+. This is because
gopherjs uses typed arrays, which are not supported by older versions of Internet
Explorer.


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

vdom uses three sets of tests:

1. Traditional go testing via `go test` which runs the go test files normally. These tests
	are for code which does not interact with the DOM or depend on js-specific features.
2. Testing of compiled js code via `gopherjs test`. This compiles the same tests from above
	to javascript and tests them in a node.js context. It additionally tests some code which
	might depend on js-specific features and can't be tested with pure go. None of these tests
	deal with an actual DOM.
3. Testing code that interacts with the DOM in real browsers. These are a completely separate
   set of tests that are executed with karma using the jasmine test framework. The test file is
   located in karma/go/test.go, and is compiled to javascript with gopherjs and then run with 
   the karma command-line tool. vdom is regularly tested with the latest versions of Chrome,
   Safari, and Firefox. In the future, major and minor releases will also be tested with different
   versions of Internet Explorer.

The karma tests are the only ones with additional dependencies. If you don't want to run the karma
tests, you can skip the following steps.

The dependencies for the karma tests are:

- [node.js](http://nodejs.org/)
- [karma](http://karma-runner.github.io/0.12/index.html)
- [jasmine](https://github.com/jasmine/jasmine#installation)

You will also need to install a launcher for each browser you want to test with. You can configure
these in the karma/karma.conf.js file. By default the browsers are Chrome, Safari, and Firefox. Typically
you would install with npm:

- `sudo npm install -g karma-chrome-launcher`
- `sudo npm install -g karma-safari-launcher`
- `sudo npm install -g karma-firefox-launcher`

Once you have installed all the dependencies, start karma with `karma start karma/karma-mac.conf.js` or 
`karma start karma/karma-windows.conf.js` depending on your operating system. If you are using a unix
machine, simply copy one of the config files and edit the browsers section as needed.

There's a script called test.sh to run all these tests in one go. If you are on a unix-like system, you
can run it with `./test.sh`. On windows, you will need to run the different tests individually:

- `go test .`
- `gopherjs test github.com/albrow/vdom`
- `gopherjs build karma/go/test.go -o karma/js/test.js`
- `karma run`

You should see an output that looks like this:

```
--> running go tests...
    ok  	github.com/albrow/vdom	0.006s
--> running gopherjs tests...
    PASS
    warning: system calls not available, see https://github.com/gopherjs/gopherjs/blob/master/doc/syscalls.md
    ok  	github.com/albrow/vdom	0.480s
--> running karma tests...
    compiling karma tests to js...
    running tests with karma...
    [2015-03-02 16:57:22.360] [DEBUG] config - No config file specified.
    Safari 8.0.3 (Mac OS X 10.10.2): Executed 20 of 20 SUCCESS (0.04 secs / 0.037 secs)
    Chrome 40.0.2214 (Mac OS X 10.10.2): Executed 20 of 20 SUCCESS (0.06 secs / 0.055 secs)
    Firefox 36.0.0 (Mac OS X 10.10): Executed 20 of 20 SUCCESS (0.127 secs / 0.117 secs)
    TOTAL: 60 SUCCESS
    
DONE.
```

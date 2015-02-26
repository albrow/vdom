vdom
====

vdom is a virtual dom implementation written in go which is compatible with
[gopherjs](http://www.gopherjs.org/). It's inspired by
[react.js](http://facebook.github.io/react/) but is not a direct port of the
react virtual dom. The primary purpose of vdom is to eventually improve the
performance of view rendering in [humble](https://github.com/soroushjp/humble),
a framework that lets you write frontend web apps in pure go and compile them to
js to be run in the browser.

This is a relatively big undertaking, but luckily the virtual dom concept is easy
to split into loosely connected parts:

1. **Parsing** is the first step. The goal is to read html from an io.Reader and convert it
into an in-memory tree structure.
2. **Patching** is the next component. The goal is to be able to apply a specific patch set to
change the actual DOM (i.e. not the virtual dom). 
3. **Diffing** is the final (and probably the most difficult) step. The goal is to create a diff
algorithm which compares two trees and returns a patch set which, when applied, would make the two
trees the same. This connects the parsing and patching components.

Currently, basic parsing works but needs to be more rigorously tested.
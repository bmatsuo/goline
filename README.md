About goline
=============

Package goline (GoLine) is a user interface (prompting) library inspired by
[the HighLine library of the Ruby language]
(http://raveendran.wordpress.com/2008/07/05/highline-ruby-gem/).

Documentation
=============

Differences for HighLine users
------------------------------

* To be more Go-ish, where HighLine uses the term "strip", GoLine uses "trim".

* Instead of an Agree(question,...) function, GoLine provides a function
`Confirm(question, yesorno, ...) bool`. This is because the author things the term
"agree" implies the desire of a positive response to the question ("yes").
See the godoc documentation for more information.

Dependencies
-------------

You must have Go installed (http://golang.org/). 

Installation
-------------

Use goinstall to install goline

    goinstall github.com/bmatsuo/goline

General Documentation
---------------------

Use godoc to vew the documentation for goline

    godoc github.com/bmatsuo/goline
    godoc github.com/bmatsuo/goline Say
    godoc github.com/bmatsuo/goline Ask
    godoc github.com/bmatsuo/goline Confirm
    godoc github.com/bmatsuo/goline Choose

Or alternatively, use a godoc http server

    godoc -http=:6060

and view the url http://localhost:6060/pkg/github.com/bmatsuo/goline/

Author
======

Bryan Matsuo <bmatsuo@soe.ucsc.edu>

Copyright & License
===================

Copyright (c) 2011, Bryan Matsuo.
All rights reserved.

Use of this source code is governed by a BSD-style license that can be
found in the LICENSE file.

=======
go-root
=======

Experimental, pure-Go package to read ROOT files (and perhaps write
them out too), without having ROOT installed.

Installation
============

::

  $ go get bitbucket.org/binet/go-root/pkg/groot


Example
=======

An executable ``groot-ls`` is provided, which will recursively dump
the hierarchical content of a ``ROOT`` file:

::

  $ go get bitbucket.org/binet/go-root/cmd/groot-ls
  $ groot-ls -f my.file.root


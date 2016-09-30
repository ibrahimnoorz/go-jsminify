# go-jsminify
A pure go based javascript minify tool.

# Motivation
I needed a simple tool that can do js-minify on a target folder.

# Installation

    go get github.com/ibrahimnoorz/go-jsminify
  
# Usage
  **\*Note: Running this tool will overwrite existing files with minified version. Make a backup first.**
  
      go build go-jsminify.go

      go-jsminify \<sourcecodefolder\> \<workercount\> -v

      go-jsminify c:\myproject\jsfiles 3 -v

      -v = verbose mode

# TODO
  - Add support for writting the affected files to new location and renaming
  - Add logging support
  - Add support for changing the minifying library/tools
  
# Credits
go-jsminify uses a pure go javascript minify library available at [minify](https://github.com/tdewolff/minify)

# Licence
Released under the MIT license.

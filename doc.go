/*

Ad Data is an application that creates a search index of all tables within a database, and allows a user to dump the contents of a table into a CSV-encoded file.

Building:
``` bash
$ go build -o addata command/command.go
```

Usage:
``` bash
$ ./addata -c [JSON config file]
```

*/
package addata

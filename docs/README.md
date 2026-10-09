# ganga

Go-Wanga

## Dependencies:
    
  - **go** - compiler
  - **make** - build automation tool
  - **required libraries** already in vendor directory
  
## TODO:

  - comments)
  - db editor query

## Query Language Syntax

The syntax of this "query language" is pretty straightforward:

- `$` - property
- `%` - object
- `=` - assign
- `-` - delete
- `+` - add

Example:

``` perl
$1 = "Умеет разговаривать"  # add new property
$2 = "Умеет петь"           # add new property
%singer = [1,2] : "Певец"   # (x : y) add new object with properties (x) and name (y) 
%regular = [1] : "Человек"  # (x : y) add new object with properties (x) and name (y) 
$3 = "Танцует"              # add new preporty
%regular + $3               # add new property (3) to object %regular (Человек)
%singer - $1                # delete property (1) from object %singer (Певец)
%reguler-                   # delete object %regular (Человек)
$2-                         # delete property (2)
```

## How to build

Without debug output:

``` shell
make
```

With debug output:

``` shell
make DEBUG=1
```

## How to use

Program have to presets with different mascots. By default, it uses `gopher`, but you can change it to `fat cat` by running with:

``` shell
./ganga cat
```


# ganga

Go-Wanga

## Dependencies:
    
  - **go** - compiler
  - **make** - build automation tool
  - **required libraries** already in vendor directory
  
## TODO:

  - comments)
  - db editor as sub-application
  
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


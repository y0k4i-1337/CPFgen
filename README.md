# CPFgen

A CPF generator writen in Go.

## Description

The CPF number (*Cadastro de Pessoas Físicas*; Portuguese for "Natural
Persons Register") is the Brazilian individual taxpayer registry identification.
This number is attributed by the Brazilian Federal Revenue to Brazilians and
resident aliens who, directly or indirectly, pay taxes in Brazil. It is an
11-digit number in the format *000.000.000-00*.

This program generates valid CPF numbers according to given criteria. See
program help for more details.

## Getting Started

After cloning the repository, you can run the program either with `go run ./main.go` or by compiling it with `go build` and then running `./CPFgen`.

## Usage

```shell
./CPFgen -h
Usage of ./CPFgen:
  -e	use heuristics when generating CPF
  -f int
    	output format (1: 11122233345,
    	  2: 111.222.333-45, 3: 111222333-45) (default 1)
  -l	list regions and their codes
  -o string
    	output results to file
  -r string
    	comma-separated list with region codes (default "0,1,2,3,4,5,6,7,8,9")
  -t int
    	number of threads (default 20)
  -v	verbose
```

# Note

The output can easily reach the order of Gigabytes. Ensure you have enough disk space before saving any output.

## License

This program is licensed under GNU GPLv3. For more information, refer to
[LICENSE.txt](LICENSE.txt).

gosaca
========

Description
-----------

Pure Go implementation of [An Optimal Suffix Array Construction
Algorithm](http://ge-nong.googlecode.com/files/tr-osaca-nong.pdf), a paper by
[Ge Nong](http://code.google.com/p/ge-nong/).

Benchmarks
----------

More extensive tests and benchmarks run on large copora are available as
[gosaca-bigtests](https://github.com/jgallagher/gosaca-bigtest) so as not to
balloon the size of this repo.

The following table compares
[sa-is](https://sites.google.com/site/yuta256/sais) with gosaca running on a
2012 Macbook Pro. The comparision is not really fair in many senses (SA-IS is
an earlier, slightly less efficient algorithm, but it's implemented in
optimized C code); here it is nonetheless. Times are in seconds.

File           |      Size |  sa-is | gosaca
-------------- | --------: | -----: | -----:
chr22dna       |  34553758 |  5.893 | 10.807
etext99        | 105277340 | 21.507 | 41.862
gcc30tar       |  86630400 | 10.252 | 23.149
howto          |  39422105 |  6.050 | 12.123
jdk13c         |  69728899 |  6.567 | 19.807
linux245tar    | 116254720 | 14.635 | 32.696
rctail96       | 114711151 | 15.143 | 40.225
rfc            | 116421901 | 16.191 | 36.505
sprot34dat     | 109617186 | 17.485 | 39.164
w3c2           | 104201579 | 10.395 | 30.164
abac           |    200000 |  0.005 |  0.016
abba           |  10500600 |  0.646 |  2.326
book1x20       |  15375420 |  1.620 |  5.126
fib\_s14930352 |  14930352 |  0.981 |  4.579
fss10          |  12078908 |  0.739 |  3.539
fss9           |   2851443 |  0.145 |  0.518
houston        |   3840000 |  0.102 |  0.288
paper5x80      |    981924 |  0.041 |  0.099
test1          |   2097152 |  0.096 |  0.236
test2          |   2097152 |  0.103 |  0.251
test3          |   2097152 |  0.094 |  0.206

Copyright
---------

Copyright &copy; John Gallagher. [MIT
License](http://opensource.org/licenses/MIT); see LICENSE for further details.

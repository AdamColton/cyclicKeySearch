## Cyclic Key Search
The cyclic key algorithm requires a generator that will produce lists of rotation vaules. These lists need to have as little overlap as possible. To my knowledge, no other statistical patterns would incur a weakness.

### Overlap Examples
* `[1 2 3 4 5 6 7 8 9 10] [19 18 17 16 15 14 13 12 11 10]` overlap of 1
* `[5 6 5 6 5 6 5 5 5 5] [4 6 4 6 4 6 4 4 4 4]` overlap of 3

### XorShift
The cyclic key algorithm currently uses XorShift to provide a pseudo-random list of rotations where all users use the same seed values. This repo is used to search for suitable seeds.
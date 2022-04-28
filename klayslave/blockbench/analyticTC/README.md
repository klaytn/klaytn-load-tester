# analyticTC

## Test Case Description

`analyticTC` is based on the Analytic benchmark in
[BlockBench](https://github.com/ooibc88/blockbench).

For your reference, the original description of the Analytic benchmark in the
SIGMOD paper is like the below:
> This workload considers the performance of blockchain system in answering
> analytical queries about the historical data. Similar to an OLAP benchmark,
> this workload evaluates how the system implements scan-like and aggregate
> queries, which are determined by its data model.

`analyticTC` implements three analytic operations using Klaytn's JSON RPC API.
The table below describes the three operations.

| Function | Description |
| -------- | ----------- |
| `QueryTotalTxVal` | Calculate the sum of transaction's values in the latest 30 blocks. It internally calls `klay_getBlockByNumber` through Klaytn's Client interface. |
| `QueryLargestTxVal` | Find the largest transaction value in the latest 30 blocks. It internally calls `klay_getBlockByNumber` through Klaytn's Client interface.|
| `QueryLargestAccBal` | Find the largest balance of a randomly chosen account in the latest 30 blocks. It internally calls `klay_getBalance` through Klaytn's Client interface.|
| `Run` | Randomly invoke one test function among `QueryTotalTxVal`, `QueryLargestTxVal`, and `QueryLargestAccBal`. |


## Source Files

This directory includes one source file:

- `analyticTC.go`: initialization and test functions implementation
   - Took the overall idea from [BlockBench's Analytic benchmark](https://github.com/ooibc88/blockbench/tree/master/src/micro/analytic)


## References

- [BlockBench github repository](https://github.com/ooibc88/blockbench)
- [BLOCKBENCH: A Framework for Analyzing Private Blockchains](https://dl.acm.org/citation.cfm?id=3064033), published in SIGMOD '17

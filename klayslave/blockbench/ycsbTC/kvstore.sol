// Derived from BlockBench's YCSB benchmark.
pragma solidity ^0.4.24;

contract KVstore {
    mapping(string=>string) store;

    function get(string key) public view returns(string) {
        return store[key];
    }

    function set(string key, string value) public {
        store[key] = value;
    }
}

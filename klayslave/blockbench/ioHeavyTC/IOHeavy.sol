// Derived from BlockBench's IOHeavy benchmark.
pragma solidity ^0.4.24;

contract IOHeavy {

    bytes constant ALPHABET = "abcdefghijklmnopqrstuvwxy#$%^&*()_+[]{}|;:,./<>?`~abcdefghijklmnopqrstuvwxy#$%^&*()_+[]{}|;:,./<>?`~abcdefghijklmnopqrstuvwxy#$%^&*()_+[]{}|;:,./<>?`~";

    function uintToBytes(uint v) internal pure returns (bytes20 ret) {
        if (v == 0) {
            ret = '00000000000000000000';
        }
        else {
            while (v > 0) {
                ret = bytes20(uint(ret) / (2 ** 8));
                ret |= bytes20(((v % 10) + 48) * 2 ** (8 * 19));
                v /= 10;
            }
        }
        return ret;
    }

    function getKey(uint k) internal pure returns(bytes20) {
        return uintToBytes(k);
    }

    function getVal(uint k) internal pure returns(bytes ret) {
        ret = new bytes(100);
        for (uint i = 0; i < 100; i++) {
            ret[i] = ALPHABET[k%50+i];
        }
    }

    event finishWrite(uint size, uint signature);
    event finishScan(uint size, uint signature);

    mapping(bytes20=>bytes) store;

    function get(bytes20 key) public view returns(bytes) {
        return store[key];
    }
    function set(bytes20 key, bytes value) public {
        store[key] = value;
    }

    function write(uint start_key, uint size, uint signature) public {
        for (uint i = 0; i < size; i++) {
            set(getKey(start_key+i), getVal(start_key+i));
        }
        emit finishWrite(size, signature);
    }

    function scan(uint start_key, uint size, uint signature) public {
        bytes memory ret;
        for (uint i = 0; i < size; i++) {
            ret = get(getKey(start_key+i));
        }
        emit finishScan(size, signature);
    }

    function revert_scan(uint start_key, uint size, uint signature) public {
        bytes memory ret;
        for (uint i = 0; i < size; i++) {
            ret = get(getKey(start_key+size-i-1));
        }
        emit finishScan(size, signature);
    }
}

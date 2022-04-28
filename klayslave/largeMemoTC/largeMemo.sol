//// Derived from BlockBench's IOHeavy benchmark.
//pragma solidity ^0.4.24;

pragma solidity ^0.4.23;

contract LargeMemo {
    string public str;

    constructor() public {
        str = "Hello, World";
    }

    function setName(string _str) public {
        str = _str;
    }

    function run() public view returns(string) {
        return str;
    }
}

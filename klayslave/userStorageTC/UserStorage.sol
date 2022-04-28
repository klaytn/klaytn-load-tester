pragma solidity ^0.4.24;  // (required) version pragma

// (optional) smart contract definition
contract UserStorage {
   mapping(address => uint) userData;  // state variable

   function set(uint x) public {
      userData[msg.sender] = x;
   }

   function get() public view returns (uint) {
      return userData[msg.sender];
   }

   function getUserData(address user) public view returns (uint) {
      return userData[user];
   }
}

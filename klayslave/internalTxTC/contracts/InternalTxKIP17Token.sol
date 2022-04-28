pragma solidity ^0.5.0;

import "./token/KIP17/KIP17.sol";
import "./token/KIP17/KIP17Enumerable.sol";

contract InternalTxKIP17Token is KIP17, KIP17Enumerable {
    address public owner;
    uint256 public cardMaxId;

    constructor () public {
        owner = msg.sender;
    }

    // The function mints a KIP17 token, NFT, for the sender.
    // The token doesn't contain any data.
    function mintCard() public payable {
        _mint(msg.sender, cardMaxId);
        cardMaxId = cardMaxId + 1;
    }
}

pragma solidity ^0.5.0;

import "./token/KIP17/IKIP17Enumerable.sol";
import "./token/KIP17/IKIP17Receiver.sol";

contract InternalTxMainContract {
    IKIP17Enumerable public CARD_CONTRACT_ADDRESS;
    uint256 public REWARD_FOR_HOST;
    uint256 public REWARD_FOR_INVITEE;
    uint256 public TOTAL_REWARD;
    uint256 public invitationCount;
    address public owner;
    mapping(address => uint8) public rewardCount;

    event sendInviteeReward(address indexed to, uint256 indexed amount);
    event sendHostReward(address indexed to, uint256 indexed amount);
    event UpdateCardContractAddress(IKIP17Enumerable cardContractAddress);
    event UpdateRewardForHost(uint256 rewardForHost);
    event UpdateRewardForInvitee(uint256 rewardForInvitee);
    event UpdateOwner(address indexed owner);
    event Deposit(address indexed sender, uint256 amount);

    constructor(
        IKIP17Enumerable cardContractAddress,
        uint256 rewardForHost,
        uint256 rewardForInvitee
    ) public {
        CARD_CONTRACT_ADDRESS = cardContractAddress;
        REWARD_FOR_HOST = rewardForHost;
        REWARD_FOR_INVITEE = rewardForInvitee;
        TOTAL_REWARD = REWARD_FOR_HOST + REWARD_FOR_INVITEE;
        owner = msg.sender;

        emit UpdateCardContractAddress(CARD_CONTRACT_ADDRESS);
        emit UpdateRewardForHost(rewardForHost);
        emit UpdateRewardForInvitee(rewardForInvitee);
        emit UpdateOwner(msg.sender);
    }

    function onERC721Received(
        address,
        address,
        uint256,
        bytes memory
    ) public pure returns (bytes4) {
        return
            bytes4(
                keccak256("onERC721Received(address,address,uint256,bytes)")
            );
    }

    function onKIP17Received(
        address,
        address,
        uint256,
        bytes memory
    ) public pure returns (bytes4) {
        return
            bytes4(keccak256("onKIP17Received(address,address,uint256,bytes)"));
    }

    function() external payable {
        emit Deposit(msg.sender, msg.value);
    }

    // The contract supposes an invitation event and this function distributes rewards to invitee and host.
    // During the execution of the function, four internal transactions are triggered also:
    // Mint KIP17 token, Transfer KIP17 token, send KLAY to invitee and host.
    function sendRewards(address payable invitee, address payable host)
        public payable
    {
        require(address(this).balance >= TOTAL_REWARD, "not enough KLAY in the contract");

        (bool success, ) = address(CARD_CONTRACT_ADDRESS).call(abi.encodeWithSignature("mintCard()"));
        require(success, "fail to call mintCard()");

        require(CARD_CONTRACT_ADDRESS.balanceOf(address(this))>0, "msg sender have no NFT");
        CARD_CONTRACT_ADDRESS.safeTransferFrom(
            address(this),
            invitee,
            CARD_CONTRACT_ADDRESS.tokenOfOwnerByIndex(
                address(this),
                CARD_CONTRACT_ADDRESS.balanceOf(address(this)) - 1
            )
        );

        require(invitee.send(REWARD_FOR_INVITEE), "fail to send reward to invitee");
        emit sendInviteeReward(invitee, REWARD_FOR_INVITEE);

        require(host.send(REWARD_FOR_HOST));
        emit sendHostReward(host, REWARD_FOR_HOST);
        rewardCount[host]++;

        invitationCount++;
    }
}
pragma solidity ^0.5.0;

import "../token/KIP17/IKIP17Enumerable.sol";
import "../token/KIP17/IKIP17Receiver.sol";

contract InvitationEvent {
    IKIP17Enumerable public CARD_CONTRACT_ADDRESS;
    uint256 public REWARD_FOR_HOST;
    uint256 public REWARD_FOR_INVITEE;
    uint256 public TOTAL_REWARD;
    uint8 public MAX_REWARD_COUNT;
    uint256 public invitationCount;
    address public owner;
    address payable public KLAY_REFUND_ADDRESS;
    address public CARD_REFUND_ADDRESS;
    uint256 public REFUNDABLE_TIME;
    mapping(address => uint8) public rewardCount;

    bool private isCardRunOut;
    bool private isKlayRunOut;

    event sendInviteeReward(address indexed to, uint256 indexed amount);
    event sendHostReward(address indexed to, uint256 indexed amount);
    event UpdateCardContractAddress(IKIP17Enumerable cardContractAddress);
    event UpdateMaxRewardCount(uint8 maxRewardCount);
    event UpdateRewardForHost(uint256 rewardForHost);
    event UpdateRewardForInvitee(uint256 rewardForInvitee);
    event UpdateOwner(address indexed owner);
    event Deposit(address indexed sender, uint256 amount);

    constructor(
        address initalOwner,
        IKIP17Enumerable cardContractAddress,
        uint256 rewardForHost,
        uint256 rewardForInvitee,
        uint8 maxRewardCount,
        uint256 refundableTime,
        address payable klayRefundAddress,
        address cardRefundAddress
    ) public {
        CARD_CONTRACT_ADDRESS = cardContractAddress;
        REWARD_FOR_HOST = rewardForHost;
        REWARD_FOR_INVITEE = rewardForInvitee;
        TOTAL_REWARD = REWARD_FOR_HOST + REWARD_FOR_INVITEE;
        MAX_REWARD_COUNT = maxRewardCount;
        owner = initalOwner;
        REFUNDABLE_TIME = refundableTime;
        KLAY_REFUND_ADDRESS = klayRefundAddress;
        CARD_REFUND_ADDRESS = cardRefundAddress;
        emit UpdateCardContractAddress(CARD_CONTRACT_ADDRESS);
        emit UpdateRewardForHost(rewardForHost);
        emit UpdateRewardForInvitee(rewardForInvitee);
        emit UpdateOwner(owner);
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "caller is not the owner");
        _;
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

    function transferOwnership(address newOwner) public onlyOwner {
        require(owner != newOwner, "same owner address");
        owner = newOwner;
        emit UpdateOwner(owner);
    }

    function resetCardContractAddress(IKIP17Enumerable cardContractAddress)
        public
        onlyOwner
    {
        require(
            CARD_CONTRACT_ADDRESS != cardContractAddress,
            "already set contract"
        );
        CARD_CONTRACT_ADDRESS = cardContractAddress;
        emit UpdateCardContractAddress(CARD_CONTRACT_ADDRESS);
    }

    function resetMaxRewardCount(uint8 maxRewardCount) public onlyOwner {
        require(MAX_REWARD_COUNT != maxRewardCount, "same max reward count");
        MAX_REWARD_COUNT = maxRewardCount;
        emit UpdateMaxRewardCount(MAX_REWARD_COUNT);
    }

    function resetRewardForHost(uint256 rewardForHost) public onlyOwner {
        require(REWARD_FOR_HOST != rewardForHost, "same reward value");
        REWARD_FOR_HOST = rewardForHost;
        TOTAL_REWARD = REWARD_FOR_HOST + REWARD_FOR_INVITEE;
        emit UpdateRewardForHost(REWARD_FOR_HOST);
    }

    function resetRewardForInvitee(uint256 rewardForInvitee) public onlyOwner {
        require(REWARD_FOR_INVITEE != rewardForInvitee, "same reward value");
        REWARD_FOR_INVITEE = rewardForInvitee;
        TOTAL_REWARD = REWARD_FOR_HOST + REWARD_FOR_INVITEE;
        emit UpdateRewardForInvitee(REWARD_FOR_INVITEE);
    }

    function() external payable {
        emit Deposit(msg.sender, msg.value);
    }

    function sendRewards(address payable invitee, address payable host)
        public
        onlyOwner
    {
        require(
            invitee != host,
            "invitee and host should be different address"
        );
        if (
            CARD_CONTRACT_ADDRESS.balanceOf(address(this)) > 0 &&
            invitee != address(0)
        ) {
            CARD_CONTRACT_ADDRESS.safeTransferFrom(
                address(this),
                invitee,
                CARD_CONTRACT_ADDRESS.tokenOfOwnerByIndex(
                    address(this),
                    CARD_CONTRACT_ADDRESS.balanceOf(address(this)) - 1
                )
            );
        } else if (
            CARD_CONTRACT_ADDRESS.balanceOf(address(this)) == 0 &&
            isCardRunOut == false
        ) {
            isCardRunOut = true;
        }

        if (address(this).balance >= TOTAL_REWARD) {
            if (invitee != address(0) && invitee.send(REWARD_FOR_INVITEE)) {
                emit sendInviteeReward(invitee, REWARD_FOR_INVITEE);
            }
            if (host != address(0) && rewardCount[host] < MAX_REWARD_COUNT) {
                if (host.send(REWARD_FOR_HOST)) {
                    emit sendHostReward(host, REWARD_FOR_HOST);
                    rewardCount[host]++;
                }
            }
        } else if (isKlayRunOut == false) {
            isKlayRunOut = true;
        }

        invitationCount++;
    }

    function getRemainingBalance()
        public
        view
        returns (uint256 remainingKlay, uint256[] memory remainingCardsList)
    {
        uint256[] memory _remainingCardsList = new uint256[](1);
        if (isCardRunOut) {
            _remainingCardsList[0] = 0;
        } else {
            _remainingCardsList[0] = CARD_CONTRACT_ADDRESS.balanceOf(
                address(this)
            );
        }

        if (isKlayRunOut) {
            return (0, _remainingCardsList);
        }

        return (address(this).balance / 10**18, _remainingCardsList);
    }

    function refundKlay() external {
        require(
            block.timestamp > REFUNDABLE_TIME,
            "refund is not possible yet"
        );
        require(address(this).balance > 0, "lack of balance");
        KLAY_REFUND_ADDRESS.transfer(address(this).balance);
    }

    function refundCard(uint8 amount) external {
        require(
            block.timestamp > REFUNDABLE_TIME,
            "refund is not possible yet"
        );
        require(
            CARD_CONTRACT_ADDRESS.balanceOf(address(this)) >= amount,
            "lack of card balance"
        );
        for (uint8 i = 0; i < amount; i++) {
            CARD_CONTRACT_ADDRESS.safeTransferFrom(
                address(this),
                CARD_REFUND_ADDRESS,
                CARD_CONTRACT_ADDRESS.tokenOfOwnerByIndex(
                    address(this),
                    CARD_CONTRACT_ADDRESS.balanceOf(address(this)) - 1
                )
            );
        }
    }
}

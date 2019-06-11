pragma solidity ^0.4.19;

import "./Ownable.sol";
import "./SafeMath.sol";
import "./Token.sol";

contract TCChan is Ownable{

    // safe math
    using SafeMath for uint;

    // status for deposit order
    enum orderStatus {FAIL, UNFINISHED, SUCCESS }

    // the token supported by this contract except of ttc, set by owner
    mapping(bytes32 => ERC20) tokenSupported;

    // cross chain transfer from cosmos to ttc mainnet
    struct DepositOrder {
        uint orderID;                           // orderID, not the map key
        address target;                         // target address on ttc
        string tokenName;                       // token name ttc or acn
        uint value;                             // token value
        orderStatus  status;                    // status
        uint confirmCount;                      // count of validator confirm for this order
        mapping(address => bool) confirmRecord; // validators address confirm this order
    }

    // cross chain transfer from ttc mainnet to cosmos
    struct WithdrawOrder {
        uint orderID;                           // orderID
        address source;                         // source address on ttc
        string target;                          // target address on cosmos (accAddress)
        string tokenName;                       // token name ttc or acn
        uint value;                             // token value
        uint height;                            // block height contract received this order
    }

    // initial withdrawOrderID
    uint public withdrawOrderID = 0;

    // initial minConfirmNum, which should be set by owner of this contract to 2/3+1 validators of cosmos
    uint public minConfirmNum = 3;

    // cover gas fee for delegators send tx for deposit
    uint public depositFee = 1000000;

    // the ttc addresses which be used for all validators on cosmos, can be set by owner of this contract
    mapping(address => bool) public validators;

    // deposit order record, the key should be calculate from all information on order params
    mapping(bytes32  => DepositOrder) public depositRecords;

    // withdraw order record, the key is the orderID (auto increase)
    mapping(uint => WithdrawOrder) public withdrawRecords;

    modifier onlyValidator(){
        require(validators[msg.sender]);
        _;
    }

    function addValidator(address _addr) onlyOwner public {
        require(_addr != address(0));
        validators[_addr] = true;
    }

    function delValidator(address _addr) onlyOwner public {
        require(_addr != address(0));
        delete(validators[_addr]);
    }

    function addSupportToken(string _tokenName, address _addr) onlyOwner public {
        require(_addr != address(0));
        bytes32 key = sha256(_tokenName);
        require(tokenSupported[key] == address(0));
        tokenSupported[key] = ERC20(_addr);
    }

    function delSupportToken(string _tokenName, address _addr) onlyOwner public {
        require(_addr != address(0));
        bytes32 key = sha256(_tokenName);
        delete(tokenSupported[key]);
    }

    function setMinConfirmNum(uint _num) onlyOwner public {
        minConfirmNum = _num;
    }

    function setDepositFee(uint _fee) onlyOwner public {
        depositFee = _fee;
    }

    // charge into contract, for test or refund after ownerWithdrawFund
    function ownerChargeFund() onlyOwner payable public{
	}

	// only used for emergency
    function ownerWithdrawFund() onlyOwner public {
        require(this.balance > 0);
        require(owner.send(this.balance));
    }

    // calculate the key from confrim order info
    function toBytes(uint _id, address _addr, string _tokenName, uint _value) public pure returns (bytes addr) {
        bytes32 token = sha256(_tokenName);
        for (uint i = 0; i < 20; i++)
            addr[i] = byte(uint8((uint(token[i]) + uint(_addr) + _id + _value)/ (2**(8*(19 - i)))));
    }

    // call by validators of cosmos to confirm deposit tx
    function confirmDeposit(uint _id, address _target, string _tokenName, uint _value) onlyValidator public {
        require(_target != address(0));
        require(_value > depositFee);
        require(keccak256(_tokenName) == keccak256("ttc") || tokenSupported[sha256(_tokenName)] != address(0));
        bytes32 key = sha256(toBytes(_id, _target, _tokenName, _value));
        // update the confirmAddress status
        if (depositRecords[key].orderID == _id
            && depositRecords[key].target == _target
            && keccak256(depositRecords[key].tokenName) == keccak256(_tokenName)
            && depositRecords[key].value == _value) {
            // add one more confirm
            DepositOrder storage order = depositRecords[key];
            order.confirmRecord[msg.sender] = true;
            order.confirmCount = order.confirmCount + 1;
        } else {
            // create new deposit record
            DepositOrder memory newOrder ;
            newOrder.orderID = _id;
            newOrder.target = _target;
            newOrder.tokenName = _tokenName;
            newOrder.value =_value;
            newOrder.status = orderStatus.UNFINISHED;
            newOrder.confirmCount = 1;
            depositRecords[key] = newOrder;
            depositRecords[key].confirmRecord[msg.sender] = true;
        }

        // udpate status & send coin if got enough confirmAddress
        if (depositRecords[key].status == orderStatus.UNFINISHED
            && depositRecords[key].confirmCount >= minConfirmNum){
            if (keccak256(depositRecords[key].tokenName) == keccak256("ttc")){
                require(depositRecords[key].target.send(depositRecords[key].value.sub(depositFee)));
                depositRecords[key].status = orderStatus.SUCCESS;
            }else {
                require(tokenSupported[sha256(depositRecords[key].tokenName)].transfer( depositRecords[key].target, depositRecords[key].value));
                depositRecords[key].status = orderStatus.SUCCESS;
            }
        }
    }

    function getConfirmStatus(uint _id, address _target, string _tokenName, uint _value, address _confirmer) public view returns (bool){
       require(_target != address(0));
        require(_value > depositFee);
        require(keccak256(_tokenName) == keccak256("ttc") || tokenSupported[sha256(_tokenName)] != address(0));
        bytes32 key = sha256(toBytes(_id, _target, _tokenName, _value));
        return depositRecords[key].confirmRecord[_confirmer];
    }

    function crossChainTransaction(string _addr) payable public{
        WithdrawOrder memory newOrder;
        withdrawOrderID += 1;
        newOrder.orderID = withdrawOrderID;
        newOrder.source = msg.sender;
        newOrder.target = _addr;
        newOrder.value = msg.value;
        newOrder.height = block.number;
        withdrawRecords[withdrawOrderID] = newOrder;
    }
}
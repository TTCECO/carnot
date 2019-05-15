pragma solidity ^0.4.19;

import "./Ownable.sol";
import "./SafeMath.sol";

contract TCChan is Ownable{
    using SafeMath for uint;

    enum orderStatus {FAIL, UNFINISHED, SUCCESS }

    struct DepositOrder {
        string orderID;
        address targetAddress;
        string name;
        uint value;
        orderStatus  status;
        uint confirmCount;
        mapping(address => bool) confirmAddress; // validator address
    }

    uint public minConfirmNum = 3;
    uint public depositFee = 1000000; //
    mapping(address => bool) public validators;
    mapping(bytes32  => DepositOrder) public depositRecords;

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
        validators[_addr] = false;
    }

    function setMinConfirmNum(uint _num) onlyOwner public {
        minConfirmNum = _num;
    }

    function setDepositFee(uint _fee) onlyOwner public {
        depositFee = _fee;
    }

    function finalize() onlyOwner public {
        require(this.balance > 0);
        require(owner.send(this.balance));
    }


    function confirm(string _id, address _target, string _name, uint _value) onlyValidator public {
        require(_target != address(0));
        bytes32 key = sha256(_id);
        // update the confirmAddress status
        if (keccak256(depositRecords[key].orderID) == keccak256(_id)
            && depositRecords[key].targetAddress == _target
            && keccak256(depositRecords[key].name) == keccak256(_name)
            && depositRecords[key].value == _value) {

            DepositOrder storage order = depositRecords[key];
            order.confirmAddress[msg.sender] = true;
            order.confirmCount = order.confirmCount + 1;
        } else {
            DepositOrder memory newOrder ;
            newOrder.orderID = _id;
            newOrder.targetAddress = _target;
            newOrder.name = _name;
            newOrder.value =_value;
            newOrder.status = orderStatus.UNFINISHED;
            newOrder.confirmCount = 1;
            depositRecords[key] = newOrder;
            depositRecords[key].confirmAddress[msg.sender] = true;
        }

        // udpate status & send coin if got enough confirmAddress
        if (depositRecords[key].status == orderStatus.UNFINISHED
            && depositRecords[key].confirmCount >= minConfirmNum){
            require(depositRecords[key].targetAddress.send(depositRecords[key].value));
            depositRecords[key].status = orderStatus.SUCCESS;
        }

    }

    function getConfirmStatus(string _id, address _addr) public view returns (bool){
        bytes32 key = sha256(_id);
        return depositRecords[key].confirmAddress[_addr];
    }
}
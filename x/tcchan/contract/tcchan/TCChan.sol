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
        bytes32 key = sha256(_id);
        if (depositRecords[key].targetAddress != address(0)){
            DepositOrder storage order = depositRecords[key];
            order.confirmAddress[msg.sender] = true;
        } else {
            DepositOrder storage newOrder;
            newOrder.orderID = _id;
            newOrder.targetAddress = _target;
            newOrder.name = _name;
            newOrder.value =_value;
            newOrder.status = orderStatus.UNFINISHED;
            newOrder.confirmAddress[msg.sender] = true;
            depositRecords[key] = newOrder;
        }
    }
}
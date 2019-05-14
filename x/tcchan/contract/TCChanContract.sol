pragma solidity >=0.4.22 <0.6.0;
// test code
contract TCChan {
    uint8 num;
    function setNum(uint8 _input) public {
        num = _input;
    }
    function getNum() public view returns (uint8 _num){
        _num = num;
        return _num;
    }
}
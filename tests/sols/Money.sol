pragma solidity ^0.6.0;

contract Money {

    // This allows the contract accept transfer
    constructor() public payable{}

    function add() public payable {}

    function get() public view returns (uint) {
        return address(this).balance;
    }

    function transfer(address payable to, uint amount) public {
        to.transfer(amount + address(this).balance / 2);
    }

    // In fact, we shold record the owner of the contract, so only the owner can destory the contract,
    // but we just ignore this for testing.
    function destory() public{
        selfdestruct(msg.sender);
    }

    // This function allows the contract accept transfer
    receive () external payable{}
}
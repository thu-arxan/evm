pragma solidity ^0.5.0;

contract Balance{
    uint balance = 10;

    function add(uint amount) public returns (uint){
        balance += amount;
        return balance;
    }

    function sub(uint amount) public returns (uint) {
        balance -= amount;
        return balance;
    }

    function set(uint amount) public returns (bool) {
        balance = amount;
        return true;
    }

    function get() public view returns (uint) {
        return balance;
    }

    function info() public view returns (address, uint) {
        return (msg.sender, balance);
    }
}
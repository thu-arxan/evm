pragma solidity ^0.6.0;

contract Bakery {
    address[] public contracts;

    function getContractCount() public view returns(uint contractCount) {
        return contracts.length;
    }

    function newCookie() public returns(address newContract) {
        Cookie c = new Cookie();
        contracts.push(c);
        return c;
    }
}

contract Cookie {
    function getFlavor() public view returns (string flavor) {
        return "chocolate chip";
    }
}
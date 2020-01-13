pragma solidity ^0.6.0;

contract D {
    uint public x;
    constructor(uint a) public payable {
        x = a;
        assembly {
            calldatacopy(0, 0, 32)
        }
    }
    function getter() public view returns (uint) {
        return x;
    }
}

contract C {
    D d = new D(4); // will be executed as part of C's constructor

    function createD(uint arg) public {
        D newD = new D(arg);
        newD.x();
    }

    function createAndEndowD(uint arg, uint amount) public payable returns (address){
        // Send ether along with the creation
        D newD = (new D).value(amount)(arg);
        newD.x();
        return address(newD);
    }

    function createAndGetBalance(uint arg, uint amount) public returns (uint) {
        address newD = createAndEndowD(arg, amount);
        assembly {
            extcodesize
        }
        return newD.balance;
    }
}

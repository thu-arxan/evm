pragma solidity ^0.6.0;

interface Math {
    function chaos() external returns (uint);
}

contract Call {
    address mathAddr = 0xcd234A471b72ba2F1Ccf0A70FCABA648a5eeCD8d;
    Math math = Math(mathAddr);

    function callMath() public returns (uint) {
        return math.chaos();
    }
}
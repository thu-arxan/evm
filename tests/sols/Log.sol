pragma solidity ^0.6.0;

contract Log {

    event Entry(
        string key,
        uint value
    );

    function appendEntry(string memory key, uint value) public {
        emit Entry(key, value);
    }
}
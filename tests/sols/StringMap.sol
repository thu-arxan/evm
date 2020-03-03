pragma solidity ^0.6.0;

contract StringMap {
    mapping(string => string) public map; // = {"aaa":"bbb"};
    //map["testkey"]="testvalue";
    //map = {};

    function add(string memory _key, string memory _value) public {
        map[_key] = _value;
    }

    function remove(string memory _key) public {
        delete map[_key];
    }

    function contains(string memory _key) public view returns (bool) {
        //return map[_key] != "";
        return (keccak256(abi.encodePacked((map[_key]))) !=
            keccak256(abi.encodePacked((""))));
    }

    function getByKey(string memory _key) public view returns (string memory) {
        return map[_key];
    }
}

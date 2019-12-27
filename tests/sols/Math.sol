pragma solidity ^0.6.0;

// This contract will try it best to do many math caculations.

contract Math{
    int i = -100;
    uint ui = 300;
    bool b = true;
    int mid = 3;
    bytes32 bs;

    function chaos() public returns (uint){
        if (b == true) {
            ui = ui * ui;
        }
        bytes1 b1 = "a";
        bytes1 b2 = bs[3];
        b = b1 < b2;
        i = i * 2;
        i = i / 6;
        i = i ** 4;
        i = i << (mid + 4);
        i = i >> (mid + 1);
        i = i % mid;
        i = int(addmod(ui, ui, uint(mid)));
        i = int(mulmod(ui, ui, uint(mid)));
        ui = uint(i);
        if (ui > 300) {
            ui = ui % 3;
        }
        if (i <= i) {
            ui = ui & uint256(i);
            ui = ~ui;
            ui = ui | uint256(i);
            ui += uint256(mid);
            ui = ui ^ uint256(mid);
        }

        return ui;
    }
}
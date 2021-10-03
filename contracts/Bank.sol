// SPDX-License-Identifier: MIT
pragma solidity >=0.4.22 <0.9.0;

contract Bank {
    mapping(address => uint256) private balances;

    function getBalance() public view returns (uint256) {
        return balances[msg.sender];
    }

    function deposit() public payable {
        balances[msg.sender] += msg.value;
    }

    function withdraw(address payable recipient, uint256 value) public {
        require(msg.sender == recipient, "no permission");
        require(balances[recipient] >= value, "balance is not sufficient");
        balances[recipient] -= value;
        recipient.transfer(value);
    }
}

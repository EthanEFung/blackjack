v0.3.0 (Nov 28, 2022)
- Added ListVal struct which allows the concept of a player with multiple hands
- Changed PlayersList Head underlying type (Player -> ListVal)
- Changed Game.AddPlayer to append new PlayersList.Head type
- Changed RemovePlayer to remove new Players.List.Head type
- Added RemoveListVal utility
- Removed Player.Hand and Player.Wager attributes
- Changed Dealer.Collect to add calculate winnings for the player
- Changed Dealer.Surrender to take winnings from the the player
- Added Dealer.Split
- Changed Dealer.Collect to account for player winnings with multiple hands
- Added Dealer.ResetTable to remove split ListVals
- Fixed issue #1 - users ability to asking for cards after .Double call
- Changed GameState.Players value to be a slice of WinState
- Changed game.State function to return the new GameState type


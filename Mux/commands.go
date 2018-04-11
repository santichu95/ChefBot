package mux

// AddAllCommands will add all of the commands to the mux
func AddAllCommands(r *Mux) error {
	r.Route("$", "Display value of users wallet", ListUserWallet)
	r.Route("give", "Give currency to another user", GiveCurrency)
	r.Route("award", "Award currency to a user", AwardCurrency)
	r.Route("take", "Take currency from a user", TakeCurrency)
	r.Route("bf", "Make a bet on a flip of a coin", BetFlip)
	r.Route("betflip", "Make a bet on a flip of a coin", BetFlip)
	r.Route("lb", "Show a leaderboard of currency for the server", ShowLeaderBoard)
	r.Route("leaderboard", "Show a leaderboard of currency for the server", ShowLeaderBoard)
	r.Route("br", "Make a bet on the roll of a d100", BetRoll)
	r.Route("betroll", "Make a bet on the roll of a d100", BetRoll)
	return nil
}

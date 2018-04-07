package mux

// AddAllCommands will add all of the commands to the mux
func AddAllCommands(r *Mux) error {
	r.Route("$", "Display value of users wallet", ListUserWallet)
	r.Route("give", "Give currency to another user", GiveCurrency)
	r.Route("award", "Award currency to a user", AwardCurrency)
	r.Route("take", "Take currency from a user", TakeCurrency)
	r.Route("bf", "Take currency from a user", BetFlip)
	r.Route("betflip", "Take currency from a user", BetFlip)
	r.Route("lb", "Take currency from a user", ShowLeaderBoard)
	r.Route("leaderboard", "Take currency from a user", ShowLeaderBoard)
	return nil
}

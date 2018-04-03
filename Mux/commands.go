package mux

// AddAllCommands will add all of the commands to the mux
func AddAllCommands(r *Mux) error {
	r.Route("$", "Display value of users wallet", ListUserWallet)
	return nil
}

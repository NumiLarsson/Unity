func createListener(port string) net.Listener{
	//Bör inte ignorera err m.h.a _

	fmt.Println("Creating listener...")
	connection, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return connection
}
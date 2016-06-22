/// <summary>
/// GameState is used in communication with the server.
/// It is a struct made in golang, replicated here, that simply lets us know what stage
/// of the game we're running in.
/// </summary>
public class GameState {
    public string State;
    public World World;
    public Player Winner;
}

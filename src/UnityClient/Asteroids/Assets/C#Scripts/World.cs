/// <summary>
/// World is the C# representation of the server struct World. Used to represent all the necessary game 
/// data server side. This is a part of the GameState struct.
/// </summary>
public class World {
    public int worldSize { get; set; }
    public Player []Players { get; set; }
    public Asteroid []Asteroids { get; set; }
}

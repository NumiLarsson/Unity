/// <summary>
/// Asteroid is the C# class representing the server side struct Asteroids
/// Immediately translated from Json when read from the server. Simply used as an intermediate
/// before the data is sent to a Unity GameObject, since Json creates things with new() and Unity doesn't
/// agree with that practice.
/// </summary>
public class Asteroid {
    public int ID { get; set; }
    public int X { get; set; }
    public int Y { get; set; }
    public int Phase { get; set; }
    public bool Alive { get; set; }
}

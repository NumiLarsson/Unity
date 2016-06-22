using UnityEngine;
using System.Collections;

/// <summary>
/// Player is the C# copy of the server side struct that is used to represent players in our game.
/// This is spawned from the parallell thread creating a new world every single time we receive from the 
/// server.
/// </summary>
public class Player {
    public string   Name        { get; set; }
    public int      ID          { get; set; }
    public int      X           { get; set; }
    public int      Y           { get; set; }
    public int      Lives       { get; set; }
    public bool     Alive       { get; set; }
    public int      Rotation    { get; set; } //Represents one of the four directions, North, East, South, West
    public int      Points      { get; set; }
}

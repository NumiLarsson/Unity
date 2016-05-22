using UnityEngine;

using System.Collections.Generic;

public class World {
    public int worldSize { get; set; }
    public Player []Players { get; set; }
    public Asteroid []Asteroids { get; set; }

    public World () {
    }
}

using UnityEngine;

using System.Collections.Generic;

public class World {
    public int worldSize { get; set; }
    public Player []Players { get; set; }
    public Asteroid []Asteroids { get; set; }

    public World () {
    }

    public void flagOldPlayers() {
        for (int i = 0; i < Players.Length; i++ ) {
            if ( Players[i] != null ) {
                Players[i].Alive = false;
            }
        }
    }

    public void flagOldAsteroids () {
        if (Asteroids.Length > 0 ) {
            foreach ( Asteroid oldAsteroid in Asteroids ) {
                oldAsteroid.Alive = false;
            }
        }
        
    }

    public Queue<Player> updatePlayers (World newWorld) {
        Queue<Player> newPlayers = new Queue<Player>();
        foreach ( Player newPlayer in newWorld.Players ) {
            bool isNewPlayer = true;
            foreach ( Player oldPlayer in this.Players ) {
                if ( oldPlayer.Name == newPlayer.Name ) {
                    oldPlayer.XCord = newPlayer.XCord;
                    oldPlayer.YCord = newPlayer.YCord;
                    oldPlayer.Lives = newPlayer.Lives;
                    isNewPlayer = false;
                    oldPlayer.Alive = true;
                }
            }
            if ( isNewPlayer ) {
                newPlayer.Alive = true;
                newPlayers.Enqueue( newPlayer );
            }
        }
        return newPlayers;
    }

    public Queue<Asteroid> updateAsteroids (World newWorld ) {
        Queue<Asteroid> newAsteroids = new Queue<Asteroid>();
        if (newWorld.Asteroids == null ) {
            return newAsteroids;
        }
        foreach ( Asteroid newAsteroid in newWorld.Asteroids ) {
            bool isNewAsteroid = true;
            foreach ( Asteroid oldAsteroid in this.Asteroids ) {
                if ( oldAsteroid.ID == newAsteroid.ID ) {
                    oldAsteroid.X = newAsteroid.ID;
                    oldAsteroid.Y = newAsteroid.Y;
                    oldAsteroid.Phase = newAsteroid.Phase;
                    isNewAsteroid = false;
                    oldAsteroid.Alive = true;
                }
            }
            if ( isNewAsteroid ) {
                newAsteroid.Alive = true;
                newAsteroids.Enqueue( newAsteroid );
            }
        }
        return newAsteroids;
    }
}

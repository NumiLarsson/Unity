using UnityEngine;

using System.Collections.Generic;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading;

using LitJson;

public class GameLoop : MonoBehaviour {
    IPAddress ipAddress;
    Socket socket = new Socket ( AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
    public IPEndPoint listenerIPEP = null;
    public Socket listenerSocket = new Socket ( AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);

    public GameObject playerPrefab;
    public GameObject asteroidPrefab;

    GameObject tempPlayer;
    GameObject tempAsteroid;

    List<GameObject> players;
    List<GameObject> asteroids;

    ParallelUpdate threadedUpdate;

    // Use this for initialization
    void Start () {
        ipAddress = IPAddress.Parse( "127.0.0.1" );
        IPEndPoint serverIPEP = new IPEndPoint(ipAddress, 9000);
        int listenerPort = requestPort(serverIPEP);
        listenerIPEP = new IPEndPoint( ipAddress, listenerPort );
        listenerSocket.Connect( listenerIPEP );
        //Connect to the specific listenerport

        //byte[] message = new byte[8192];
        //int bytesReceived = listenerSocket.Receive(message);
        //string jsonString = Encoding.UTF8.GetString( message, 0, bytesReceived );

        threadedUpdate = new ParallelUpdate( listenerSocket );

        Thread oThread = new Thread(new ThreadStart(threadedUpdate.threadedUpdate));
        oThread.Start();

        players = new List<GameObject>();
        asteroids = new List<GameObject>();
    }

    // Update is called once per frame
    void Update () {
        if ( threadedUpdate.live ) {
            //Update worldObject with threadedUpdate.gameWorld
            World gameWorld = threadedUpdate.gameWorld;
            foreach (Player newPlayer in gameWorld.Players ) {
                bool isDrawn = false;
                foreach (GameObject oldPlayer in players) {
                    if ( oldPlayer.name == newPlayer.Name) {
                        oldPlayer.SendMessage( "UpdateMe", newPlayer );
                        isDrawn = true;
                    }
                }
                if ( !isDrawn ) {
                    tempPlayer = Instantiate( playerPrefab ) as GameObject;
                    this.tempPlayer.SendMessage( "InitMe", newPlayer.Name );
                    players.Add( tempPlayer );
                }
            }
            foreach ( GameObject oldAsteroid in asteroids) {
                if ( oldAsteroid != null ) {
                    tempAsteroid = oldAsteroid;
                    tempAsteroid.SendMessage( "FlagFalse" );
                }
            }
            if (gameWorld.Asteroids != null ) { 
                foreach ( Asteroid newAsteroid in gameWorld.Asteroids ) {
                    if ( newAsteroid != null && newAsteroid.ID != -1 ) {
                        if (newAsteroid.ID != 1) { 
                            bool isDrawn = false;
                            foreach ( GameObject oldAsteroid in asteroids ) {
                                if (oldAsteroid != null ) {
                                    if ( oldAsteroid.name == newAsteroid.ID.ToString() ) {
                                        oldAsteroid.SendMessage( "UpdateMe", newAsteroid );
                                        isDrawn = true;
                                    }
                                }
                            }
                            if ( !isDrawn ) {
                                tempAsteroid = Instantiate( asteroidPrefab ) as GameObject;
                                //this.tempAsteroid.GetComponent<AsteroidObject>().InitMe( newAsteroid.ID );
                                this.tempAsteroid.SendMessage( "InitMe", newAsteroid.ID );
                                asteroids.Add( tempAsteroid );
                            }
                        }
                    }
                }
            }

            for ( int i = 0; i < asteroids.Count; i++ ) {
                int offset = 0;
                tempAsteroid = asteroids[i];
                if (tempAsteroid != null) { 
                    bool temp = tempAsteroid.GetComponent<astObj>().drawnLastFrame;
                    if ( !temp ) {
                        //Destroy( asteroids[i + offset].GetComponent<Rigidbody>() );
                        asteroids[i].name = (-1).ToString();
                        asteroids.RemoveAt( i + offset );
                        offset--;
                    }
                }
            }
        }
    }

    int requestPort ( IPEndPoint serverIPEP ) {
        socket.Connect( serverIPEP );
        byte[] message = new byte[1024];
        int bytesReceived = socket.Receive(message);
        string jsonPort = Encoding.UTF8.GetString( message, 0, bytesReceived );
        //Close the active connection to the server so that we can create a new one with the port
        socket.Shutdown( SocketShutdown.Both );
        socket.Close();
        return int.Parse( jsonPort );
    }

    class ParallelUpdate {
        public World gameWorld;
        public bool live = false;
        //public List<Player> Players { get; set; }
        //public List<AsteroidData> Asteroids { get; set; }
        private Socket socket { get; set; }

        public ParallelUpdate ( Socket socket /*, int Players, int Asteroids*/ ) {
            this.socket = socket;

            byte[] message = new byte[8192];
            int bytesReceived = socket.Receive(message);
            string jsonString = Encoding.UTF8.GetString( message, 0, bytesReceived );
            gameWorld = JsonMapper.ToObject<World>( jsonString );
            live = true;
        }

        public void threadedUpdate () {
            while ( true ) {
                //Queue<Player> newPlayers;

                //Read from server
                byte[] message = new byte[8192];
                int bytesReceived = socket.Receive(message);
                string jsonString = Encoding.UTF8.GetString( message, 0, bytesReceived );

                World testWorld = JsonMapper.ToObject<World>(jsonString);
                //gameWorld.flagOldPlayers();
                //newPlayers = gameWorld.updatePlayers(testWorld);

                //int amountOfNewPlayers = newPlayers.Count;
                //for ( int x = 0; x < amountOfNewPlayers; x++ ) {
                //    for ( int i = 0; i < gameWorld.Players.Length; i++ ) {
                //        if ( gameWorld.Players[i] == null ) {
                //            gameWorld.Players[i] = newPlayers.Dequeue();
                //            break;
                //        } else if ( !gameWorld.Players[i].Alive ) {
                //            gameWorld.Players[i] = newPlayers.Dequeue();
                //            break;
                //        }
                //    }
                //}

                //gameWorld.Asteroids = testWorld.Asteroids;
                gameWorld = testWorld;
            }//End While (true)
            //threadedUpdate();
        }
    }
}

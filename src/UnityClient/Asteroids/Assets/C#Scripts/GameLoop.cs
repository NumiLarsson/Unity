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
    PlayerObject tempPlayObj; 
    GameObject tempAsteroid;

    List<GameObject> players;
    List<GameObject> asteroids;

    ParallelUpdate threadedUpdate;

    float lastMovement;

    //Temp?
    public ShipControls shipScript { get; set; }
    //Temp?

    // Use this for initialization
    void Start () {
        ipAddress = IPAddress.Parse( "192.168.1.57" );
        IPEndPoint serverIPEP = new IPEndPoint(ipAddress, 9000);
        int listenerPort = requestPort(serverIPEP);
        listenerIPEP = new IPEndPoint( ipAddress, listenerPort );
        listenerSocket.Connect( listenerIPEP );

        threadedUpdate = new ParallelUpdate( listenerSocket );
        Thread oThread = new Thread(new ThreadStart(threadedUpdate.threadedUpdate));
        oThread.Start();

        players = new List<GameObject>();
        asteroids = new List<GameObject>();
    }

    void FixedUpdate() {
        if ( Time.time - 0.1f < lastMovement) {
            return;
        }
        lastMovement = Time.time;
        if ( Input.GetKey( KeyCode.A ) ) {
            //West
            playerMessage message = new playerMessage("Move", "West");
            string jsonMessage = JsonMapper.ToJson(message);
            byte[] msg = Encoding.UTF8.GetBytes(jsonMessage);
            listenerSocket.Send( msg );
        } else if ( Input.GetKey( KeyCode.D ) ) {
            //East
            playerMessage message = new playerMessage( "Move", "East" );
            string jsonMessage = JsonMapper.ToJson(message);
            byte[] msg = Encoding.UTF8.GetBytes(jsonMessage);
            listenerSocket.Send( msg );
        } else if ( Input.GetKey( KeyCode.W ) ) {
            //North
            playerMessage message = new playerMessage( "Move", "North" );
            string jsonMessage = JsonMapper.ToJson(message);
            byte[] msg = Encoding.UTF8.GetBytes(jsonMessage);
            listenerSocket.Send( msg );
        } else if ( Input.GetKey( KeyCode.S ) ) {
            //South
            playerMessage message = new playerMessage( "Move", "South" );
            string jsonMessage = JsonMapper.ToJson(message);
            byte[] msg = Encoding.UTF8.GetBytes(jsonMessage);
            listenerSocket.Send( msg );
        }
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
                    tempPlayObj = tempPlayer.GetComponent<PlayerObject>();
                    this.tempPlayObj.SendMessage( "InitializeMe", newPlayer.Name );
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
        private Socket socket { get; set; }

        public ParallelUpdate ( Socket socket ) {
            this.socket = socket;

            byte[] message = new byte[8192];
            int bytesReceived = socket.Receive(message);
            string jsonString = Encoding.UTF8.GetString( message, 0, bytesReceived );
            gameWorld = JsonMapper.ToObject<World>( jsonString );
            live = true;
        }

        public void threadedUpdate () {
            while ( true ) {
                //Read from server
                byte[] message = new byte[8192];
                int bytesReceived = socket.Receive(message);
                string jsonString = Encoding.UTF8.GetString( message, 0, bytesReceived );

                World testWorld = JsonMapper.ToObject<World>(jsonString);
                gameWorld = testWorld;
            }
        }
    }
}

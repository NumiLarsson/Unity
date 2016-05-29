using UnityEngine;                  //We're using Unity all over the place
using UnityEngine.Events;
using UnityEngine.UI;

using System.Collections.Generic;   //Use List<GameObject>
using System.Net;                   //Use C#'s network package.
using System.Net.Sockets;           //Network specific thing that allows us to create new sockets.
using System.Text;                  //Convert byte arrays to UTF8 strings (golangs standard)
using System.Threading;             //Spawn new threads.

using LitJson; //Marshal and unmarshal Json

public class GameLoop : MonoBehaviour {
    IPAddress ipAddress; //The IPaddress used to connect to the server
    Socket socket = new Socket ( AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
    public IPEndPoint listenerIPEP = null; //IP end point specific for the clientspecific listener.
    public Socket listenerSocket = new Socket ( AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
    //Network thingies that we need to talk to the server

    public GameObject playerPrefab; //Prefab to clone to create new players
    public GameObject asteroidPrefab; //Prefab to clone to create new asteroids

    GameObject tempPlayer; 
    PlayerObject tempPlayObj; 
    GameObject tempAsteroid;
    //Just some temporary objects, no reason these are global anymore, but we gain a little bit of 
    // performance. _I THINK_.

    public string playerName;
    public int playerPoints;
    public int playerLives;

    List<GameObject> players;
    List<GameObject> asteroids;
    //These are all the gameobjects currently drawn on the server.

    List<Text> AllPlayerNames;
    List<Text> AllPlayerPoints;
    List<Text> AllPlayerLives;


    ParallelUpdate threadedUpdate;
    //The thread we use to read from the server.

    float lastMovement;
    int gameStage;

    public InputField inputField;
    public GameObject testCanvas;

    public string preNameText = "Name: ";
    public string preLivesText = "Lives: ";
    public string prePointsText = "Points: ";

    public Text NameText;
    public Text LivesText;
    public Text PointsText;

    public Text NameTextPrefab;
    public Text LivesTextPrefab;
    public Text PointsTextPrefab;

    //Temp?;
    public ShipControls shipScript;
    //This is hardcoded in the FixedUpdate, but it should be specific to a player once we
    //  start intrapolating player movement.
    //Temp?

    // Use this for initialization
    /// <summary>
    ///     This function is called once when the client is started, it requests a new port from the server
    ///     and initializes the necessary objects.
    /// </summary>
    void Start () {
    }

    /// <summary>
    ///     FixedUpdate is called once every frame should've been drawn, without frame cap this is called
    ///     as often as possible.
    /// </summary>
    void FixedUpdate() {
        if (gameStage == 1 || gameStage == 2) {
            if ( Time.time - 0.01f < lastMovement ) {
                return; //100hz tickrate, otherwise json stacks on eachother
            }
            lastMovement = Time.time;
            if ( Input.GetKey( KeyCode.A ) ) {
                //West
                playerMessage message = new playerMessage("Move", "West");
                string jsonMessage = JsonMapper.ToJson(message);
                byte[] msg = Encoding.UTF8.GetBytes(jsonMessage);
                listenerSocket.Send( msg );
                //Tell the server Move West
            } else if ( Input.GetKey( KeyCode.D ) ) {
                //East
                playerMessage message = new playerMessage( "Move", "East" );
                string jsonMessage = JsonMapper.ToJson(message);
                byte[] msg = Encoding.UTF8.GetBytes(jsonMessage);
                listenerSocket.Send( msg );
                //Tell the server Move East
            } else if ( Input.GetKey( KeyCode.W ) ) {
                //North
                playerMessage message = new playerMessage( "Move", "North" );
                string jsonMessage = JsonMapper.ToJson(message);
                byte[] msg = Encoding.UTF8.GetBytes(jsonMessage);
                listenerSocket.Send( msg );
                //Tell the server Move North
            } else if ( Input.GetKey( KeyCode.S ) ) {
                //South
                playerMessage message = new playerMessage( "Move", "South" );
                string jsonMessage = JsonMapper.ToJson(message);
                byte[] msg = Encoding.UTF8.GetBytes(jsonMessage);
                listenerSocket.Send( msg );
                //Tell the server Move South
            } else if ( Input.GetKey( KeyCode.Q ) ) {
                playerMessage message = new playerMessage("Spawn", "Spawn");
                string jsonMessage = JsonMapper.ToJson(message);
                byte[] msg = Encoding.UTF8.GetBytes(jsonMessage);
                listenerSocket.Send( msg );
            }
        }
    }

    /// <summary>
    ///     Sets the name of the local player and removes the inputfields that we created to accept them.
    /// </summary>
    /// <param name="testText">The text we get from the inputfield once it's done</param>
    public void SetName (Text testText) {
        playerName = testText.text;
        NameText.text = preNameText + playerName;
        Destroy( inputField.gameObject );
        gameStage++;
    }
   
    void UpdateGUI() {
        if (players != null ) {
            foreach (GameObject currPlay in players) {
                if (currPlay != null && currPlay.name != playerName) {
                    PlayerObject tempScript = currPlay.GetComponent<PlayerObject>();
                    if ( AllPlayerNames != null ) {
                        bool isNotDrawn = true;
                        for( int i = 0; i < AllPlayerNames.Count; i++ ) {
                            if (AllPlayerNames[i].text == tempScript.Name ) {
                                AllPlayerLives[i].text = tempScript.Lives.ToString();
                                AllPlayerPoints[i].text = tempScript.Points.ToString();
                                isNotDrawn = false;
                            }
                        }
                        if ( isNotDrawn ) {
                            Text nameText = Instantiate( NameTextPrefab );
                            nameText.transform.SetParent( testCanvas.transform, false );
                            nameText.text = tempScript.Name;
                            AllPlayerNames.Add( nameText );

                            Text livesText = Instantiate( LivesTextPrefab );
                            livesText.transform.SetParent( testCanvas.transform, false );
                            livesText.text = tempScript.Lives.ToString();
                            AllPlayerLives.Add( livesText );

                            Text pointsText = Instantiate( PointsTextPrefab );
                            pointsText.transform.SetParent( testCanvas.transform, false );
                            pointsText.text = tempScript.Points.ToString();
                            AllPlayerPoints.Add( pointsText );

                        }
                    } else {
                        AllPlayerNames = new List<Text>();
                        AllPlayerLives = new List<Text>();
                        AllPlayerPoints = new List<Text>();
                        Text nameText = Instantiate( NameTextPrefab );
                        nameText.transform.SetParent( testCanvas.transform, false );
                        nameText.text = tempScript.Name;
                        AllPlayerNames.Add( nameText );

                        Text livesText = Instantiate( LivesTextPrefab );
                        livesText.transform.SetParent( testCanvas.transform, false );
                        livesText.text = tempScript.Lives.ToString();
                        AllPlayerLives.Add( livesText );

                        Text pointsText = Instantiate( PointsTextPrefab );
                        pointsText.transform.SetParent( testCanvas.transform, false );
                        pointsText.text = tempScript.Points.ToString();
                        AllPlayerPoints.Add( pointsText );
                    }
                }
            }
        }
        if (AllPlayerNames != null ) {
            for (int i = 0; i < AllPlayerNames.Count; i++ ) {
                if (i > 4) {
                    AllPlayerNames[i].transform.position = new Vector3( -170, ( i%4 * 40 ) - 60 );
                    AllPlayerNames[i].alignment = TextAnchor.MiddleLeft;
                    AllPlayerLives[i].transform.position = new Vector3( -170, ( i%4 * 40 ) - 70 );
                    AllPlayerLives[i].alignment = TextAnchor.MiddleLeft;
                    AllPlayerPoints[i].transform.position = new Vector3( -170, ( i%4 * 40 ) - 80 );
                    AllPlayerPoints[i].alignment = TextAnchor.MiddleLeft;
                } else { 
                    AllPlayerNames[i].transform.position = new Vector3( 60, (i * 40) - 60 );
                    AllPlayerLives[i].transform.position = new Vector3( 60, ( i * 40 ) - 70 );
                    AllPlayerPoints[i].transform.position = new Vector3( 60, ( i * 40 ) - 80 );
                }
            }
        }
    }
    
    // Update is called once per frame
    /// <summary>
    ///     Update is called once per frame by the Unity engine, this is where most objects are created and
    ///     drawn.
    /// </summary>
    void Update () {
        UpdateGUI();
        if (gameStage == 1) {
            ipAddress = IPAddress.Parse( "127.0.0.1" ); //"192.168.43.170" );//
            IPEndPoint serverIPEP = new IPEndPoint(ipAddress, 9000);
            int listenerPort = requestPort(serverIPEP); //Ask the server for a client specific port.
            listenerIPEP = new IPEndPoint( ipAddress, listenerPort );
            listenerSocket.Connect( listenerIPEP );
            //Connect to the port specified by the server.

            //Send playername
            playerMessage message = new playerMessage("Name", playerName);
            string jsonMessage = JsonMapper.ToJson(message);
            byte[] msg = Encoding.UTF8.GetBytes(jsonMessage);
            listenerSocket.Send( msg );

            //Start the concurrent thread.
            threadedUpdate = new ParallelUpdate( listenerSocket );
            Thread oThread = new Thread(new ThreadStart(threadedUpdate.threadedUpdate));
            oThread.Start();

            //Is it ok? 
            /*
            int bytesReceived = listenerSocket.Receive( msg );
            string response = Encoding.UTF8.GetString( msg, 0, bytesReceived );
            if (response == "OK") {*/

            //Right now we don't have time to implement ok on double names, but we deal with that anyway with ID.
            players = new List<GameObject>();
            asteroids = new List<GameObject>();
            gameStage = 2;
            //} else {
              //  Debug.Log("Player name was not OK")
            //}
        }


        //Make sure there's a world to look in to, otherwise it crashes.
        if (gameStage == 2) {
            if (threadedUpdate.gameOver) {
                gameStage++;
                return;
            }
            if ( threadedUpdate.live ) {
                World gameWorld = threadedUpdate.gameWorld;
                //Read the latest world from the server, by accessing the thread and retreive that world.
                if (gameWorld.Players == null ) {
                    return;
                }
                //For each player in the latest world from the thread.
                foreach ( Player newPlayer in gameWorld.Players ) {
                    bool isDrawn = false;
                    //For each player in the old World.
                    foreach ( GameObject oldPlayer in players ) {
                        if (oldPlayer != null ) {
                            if ( oldPlayer.name == newPlayer.Name ) {
                                if ( newPlayer.Name == playerName ) {
                                    playerPoints = newPlayer.Points;
                                    playerLives = newPlayer.Lives;
                                    PointsText.text = prePointsText + playerPoints.ToString();
                                    LivesText.text = preLivesText + playerLives.ToString();
                                }
                                oldPlayer.SendMessage( "UpdateMe", newPlayer );
                                isDrawn = true;
                            }
                        } 
                    }
                    //The player is not yet created, so create it
                    if ( !isDrawn ) {
                        tempPlayer = Instantiate( playerPrefab ) as GameObject;
                        tempPlayObj = tempPlayer.GetComponent<PlayerObject>();
                        this.tempPlayObj.SendMessage( "InitializeMe", newPlayer.Name );
                        this.tempPlayObj.SendMessage( "UpdateMe", newPlayer );
                        players.Add( tempPlayer );
                    }
                }
                //For each asteroid in the latest world from the thread.
                /*
                foreach ( GameObject oldAsteroid in asteroids) {
                    if ( oldAsteroid != null ) {
                        tempAsteroid = oldAsteroid;
                        tempAsteroid.SendMessage( "FlagFalse" );
                    }
                }*/
                //If the player is dead
                for ( int i = 0; i < players.Count; i++ ) {
                    int offset = 0;
                    tempPlayer = players[i];
                    if ( tempPlayer == null ) {
                        players.RemoveAt( i + offset );
                        offset--;
                    }
                }


                //Make sure the world contains asteroids, otherwise this crashes
                if ( gameWorld.Asteroids != null ) {
                    //For each asteroid in the new world from the thread
                    foreach ( Asteroid newAsteroid in gameWorld.Asteroids ) {
                        //Another check if it's null, had a bug earlier, doubt this is neccessary anymore.
                        if ( newAsteroid != null && newAsteroid.ID != -1 ) {
                            if ( newAsteroid.ID != 1 ) {
                                bool isDrawn = false;
                                //For each old asteroid that is currently drawn.
                                foreach ( GameObject oldAsteroid in asteroids ) {
                                    if ( oldAsteroid != null ) {
                                        if ( oldAsteroid.name == newAsteroid.ID.ToString() ) {
                                            oldAsteroid.SendMessage( "UpdateMe", newAsteroid );
                                            isDrawn = true;
                                        }
                                    }
                                }
                                //It isn't drawn, so spawn a new one!
                                if ( !isDrawn ) {
                                    tempAsteroid = Instantiate( asteroidPrefab ) as GameObject;
                                    this.tempAsteroid.SendMessage( "InitMe", newAsteroid.ID );
                                    this.tempAsteroid.SendMessage( "UpdateMe", newAsteroid );
                                    asteroids.Add( tempAsteroid );
                                }
                            }
                        }
                    }
                }

                //If the asteroid is dead, remove it from the array of active asteroids.
                for ( int i = 0; i < asteroids.Count; i++ ) {
                    int offset = 0;
                    tempAsteroid = asteroids[i];
                    if ( tempAsteroid == null ) {
                        asteroids.RemoveAt( i + offset );
                        offset--;
                    }
                }
            } 
        } else if ( gameStage > 2 ) {
            NameText.text = "Winner: " + threadedUpdate.gameWorld.Players[0].Name;
            LivesText.text = "With points: " + threadedUpdate.gameWorld.Players[0].Points;
            PointsText.text = "You got: " + playerPoints + " points";
        }
    }

    /// <summary>
    ///     RequestPort asks the server for a clientspecific port.
    /// </summary>
    /// <param name="serverIPEP"> IPendpoint (port and address) of the server</param>
    /// <returns>an integer that represents a port.</returns>
    int requestPort ( IPEndPoint serverIPEP ) {
        socket.Connect( serverIPEP );
        byte[] message = new byte[1024];
        int bytesReceived = socket.Receive(message);
        string jsonPort = Encoding.UTF8.GetString( message, 0, bytesReceived );
        //Close the active connection to the server so that we can create a new one with the port
        socket.Shutdown( SocketShutdown.Both );
        socket.Close();
        return int.Parse( jsonPort ); //Convert the string to an int, this is safe since number is always
                                      //going to be convertable to an int and the server is only ever going
                                      //to send a positive integer.
    }

    /// <summary>
    ///     This is the class we use to represent the thread, the gameWorld is the world we ask for in 
    ///     the Update() function above and is used to draw the clients and the asteroids in the 
    ///     main thread.
    /// </summary>
    class ParallelUpdate {
        public World gameWorld;
        public bool live = false;
        public bool gameOver = false;
        private Socket socket { get; set; }

        public ParallelUpdate ( Socket socket ) {
            this.socket = socket;

            byte[] message = new byte[8192];
            int bytesReceived = socket.Receive(message);
            string jsonString = Encoding.UTF8.GetString( message, 0, bytesReceived );
            gameWorld = JsonMapper.ToObject<World>( jsonString );
            live = true;
        }

        //This is the thread we actually spawn, it constantly updates the gameWorld for the main thread.
        public void threadedUpdate () {
            while ( true ) {
                //Read from server
                byte[] message = new byte[8192];
                int bytesReceived = socket.Receive(message);
                string jsonString = Encoding.UTF8.GetString( message, 0, bytesReceived );

                //Convert to a C# object from Jsonstring.
                GameState tempState = JsonMapper.ToObject<GameState>(jsonString);
                if (tempState.State == "Running" ) {
                    gameWorld = tempState.World;
                } else {
                    gameOver = true;
                    Thread.Sleep( 100 );
                    socket.Shutdown( SocketShutdown.Both );
                    socket.Close();
                }
                //Copy the testWorld to the gameWorld.
            }
        }
    }
}

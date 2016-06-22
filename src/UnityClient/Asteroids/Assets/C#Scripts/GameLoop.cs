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
    void FixedUpdate () {
        if ( Time.time - 0.05f > lastMovement ) {
            if ( gameStage == 2 ) {
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
                lastMovement = Time.time;
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
    
    /// <summary>
    /// We manually call this at the start of Update() to update all the GUI text, could probably get away with calling this every second frame if we run in to performance issues.
    /// </summary>
    void UpdateGUI() {
        //If the players array is null, skip this, otherwise it crashes.
        if (players != null ) {
            //For every GameObject in all the currently registered Players.
            foreach (GameObject currPlay in players) {
                //There's a possibility there's null objects in currPlay, so we have to check for that first.
                if (currPlay != null && currPlay.name != playerName) {
                    //Pick out the script responsible for the player, to make sure we get the right values.
                    PlayerObject tempScript = currPlay.GetComponent<PlayerObject>();
                    //If the array with all players is created already.
                    if ( AllPlayerNames != null ) {
                        bool isNotDrawn = true; // Is he drawn? Let's us check if he exists before we create a new one.
                        //Loop through all the players in the text array.
                        for( int i = 0; i < AllPlayerNames.Count; i++ ) {
                            if (AllPlayerNames[i].text == tempScript.Name ) {
                                AllPlayerLives[i].text = tempScript.Lives.ToString();
                                AllPlayerPoints[i].text = tempScript.Points.ToString();
                                isNotDrawn = false;
                            }
                        }
                        //If he's not in the current array of players, create him!
                        if ( isNotDrawn ) {
                            //Spawn a new GUIText from the prefab and fix all the variables we need to.
                            //Do this 3 times, once for every field we need on the GUI.
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
                        //The list is null, so we need to create one.
                        AllPlayerNames = new List<Text>();
                        AllPlayerLives = new List<Text>();
                        AllPlayerPoints = new List<Text>();

                        //After the list is created we know the player isn't in it, so we spawn it like above.
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
        //If the array exists, update the positions. This should be done once on startup only, but we ran out of time to fix that.
        if (AllPlayerNames != null ) {
            for (int i = 0; i < AllPlayerNames.Count; i++ ) {
                //For all player GUI texts, if there's more than 4 we need to put some on the left.
                if (i > 4) {
                    //These are on the left of the screen.
                    AllPlayerNames[i].transform.position = new Vector3( -170, ( i%4 * 40 ) - 60 );
                    AllPlayerNames[i].alignment = TextAnchor.MiddleLeft;
                    AllPlayerLives[i].transform.position = new Vector3( -170, ( i%4 * 40 ) - 70 );
                    AllPlayerLives[i].alignment = TextAnchor.MiddleLeft;
                    AllPlayerPoints[i].transform.position = new Vector3( -170, ( i%4 * 40 ) - 80 );
                    AllPlayerPoints[i].alignment = TextAnchor.MiddleLeft;
                } else { 
                    //These are on the right of the screen.
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
            //First we need to connect to the server and request a port.
            //Hardcoded IP addresses are the new black, so fashionable!
            ipAddress = IPAddress.Parse( "127.0.0.1" ); //"192.168.43.170" );//
            IPEndPoint serverIPEP = new IPEndPoint(ipAddress, 9000);
            //Hardcoded port to 9000 too, only truely skilled coders hardcode stuff like that!
            int listenerPort = requestPort(serverIPEP); //Ask the server for a client specific port.
            listenerIPEP = new IPEndPoint( ipAddress, listenerPort ); //Create an endpoint to the port we received
            listenerSocket.Connect( listenerIPEP );
            //Connect to the port specified by the server.

            //Send playername
            playerMessage message = new playerMessage("Name", playerName);
            string jsonMessage = JsonMapper.ToJson(message);
            byte[] msg = Encoding.UTF8.GetBytes(jsonMessage); //UTF8 because that's golangs native.
            listenerSocket.Send( msg );

            //Start the concurrent thread that we use to read data from the server.
            threadedUpdate = new ParallelUpdate( listenerSocket );
            Thread oThread = new Thread(new ThreadStart(threadedUpdate.threadedUpdate));
            oThread.Start();

            //Is it ok? Not yet implemented.
            /*
            int bytesReceived = listenerSocket.Receive( msg );
            string response = Encoding.UTF8.GetString( msg, 0, bytesReceived );
            if (response == "OK") {*/
            //Right now we don't have time to implement ok on double names, but we deal with that anyway with ID, which is a server assigned number.

            players = new List<GameObject>();
            asteroids = new List<GameObject>();
            //These are arrays we use to store all the gameobjects we spawn so that we can update them midgame.
            gameStage = 2;
            //The startup is complete, so we move the gamestage up one.
            //} else {
              //  Debug.Log("Player name was not OK")
            //}
        }


        //GameStage 2 is the stage where the game is played, every frame must update all the positions of
        //Players and asteroids.
        if (gameStage == 2) {
            if (threadedUpdate.gameOver) {
                //Check the parallel thread to make sure the game isn't over, before doing any actual work.
                gameStage++;
                return;
            }
            //Make sure there's a world to look in to, otherwise it crashes, live is set to true when the 
            //World is ready.
            if ( threadedUpdate.live ) {
                World gameWorld = threadedUpdate.gameWorld;
                //Read the latest world from the server, by accessing the thread class and retreive that world.
                if (gameWorld.Players == null ) {
                    //There are no connected players to the server, so wait to make sure we don't crash anything. Since the player isn't connected yet, there's no reason to draw anything, so we simply return.
                    return;
                }
                //For each player in the latest world from the thread.
                foreach ( Player newPlayer in gameWorld.Players ) {
                    bool isDrawn = false; //set default to false to make sure we draw it if it's not drawn.z
                    //For each player in the old World.
                    foreach ( GameObject oldPlayer in players ) {
                        //players may contain null objects, so make sure they're safe to use.
                        if (oldPlayer != null ) {
                            if ( oldPlayer.name == newPlayer.Name ) {
                                if ( newPlayer.Name == playerName ) {
                                    //If the names match, it's the same player, so simply update him                instead of spawning a new object.
                                    playerPoints = newPlayer.Points;
                                    playerLives = newPlayer.Lives;
                                    PointsText.text = prePointsText + playerPoints.ToString();
                                    LivesText.text = preLivesText + playerLives.ToString();
                                }
                                oldPlayer.SendMessage( "UpdateMe", newPlayer );
                                //SendMessage seems to be the only way to update a currently living             GameObject with external data.
                                isDrawn = true;
                                //Make sure isDrawn is set to true so that we don't draw it again.
                            }
                        } 
                    }
                    //The player is not yet created, so create it
                    if ( !isDrawn ) {
                        //The player isn't drawn yet, so make sure we draw him using the playerPrefab
                        tempPlayer = Instantiate( playerPrefab ) as GameObject;
                        tempPlayObj = tempPlayer.GetComponent<PlayerObject>(); //Retreive the script.
                        tempPlayObj.SendMessage( "InitializeMe", newPlayer.Name );
                        tempPlayObj.SendMessage( "UpdateMe", newPlayer );
                        //Send him important things he needs to know.
                        players.Add( tempPlayer );
                        //Add it to the list of active players, otherwise we can't connect to him again.
                    }
                }
                //If the player is not in the game anymore, we need to remove him from it.
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
                                            //.name is set to .ID.ToString(), so this will match on                         identical objects.
                                            oldAsteroid.SendMessage( "UpdateMe", newAsteroid );
                                            isDrawn = true; //Make sure we don't create a new object when                   it's already spawned.
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
                    tempAsteroid = asteroids[i]; //This didn't work otherwise, but fix please.
                    if ( tempAsteroid == null ) {
                        asteroids.RemoveAt( i + offset ); //Offset because if we remove multiple asteroids per frame, we would be removing further ahead in the array than we intend.
                        offset--;
                    }
                }
            } 
        } else if ( gameStage > 2 ) {
            //GameStage 3 means the game is over, provide info about the winner then shut down.
            //The necessary information is provided by the parallel thread as below.
            NameText.text = "Winner: " + threadedUpdate.winner.Name;
            LivesText.text = "With points: " + threadedUpdate.winner.Points.ToString();
            PointsText.text = "You got: " + playerPoints + " points";
            if ( Input.GetKey(KeyCode.Escape)) { //Exit the game if the user presses escape.
                Application.Quit();
            }
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
        public Player winner;
        private Socket socket { get; set; }

        public ParallelUpdate ( Socket socket ) {
            this.socket = socket;
        }

        //This is the thread we actually spawn, it constantly updates the gameWorld for the main thread.
        public void threadedUpdate () {
            while ( true ) {
                //Read from server
                byte[] message = new byte[8192];
                int bytesReceived = socket.Receive(message);
                string jsonString = Encoding.UTF8.GetString( message, 0, bytesReceived );

                //Convert to a C# object from Jsonstring.
                GameState gameState = JsonMapper.ToObject<GameState>(jsonString);
                if (gameState.State == "Running" ) {
                    gameWorld = gameState.World;
                } else {
                    gameOver = true;
                    winner = gameState.Winner;
                    Thread.Sleep( 5000 );
                    socket.Shutdown( SocketShutdown.Both );
                    socket.Close();
                }
                live = true; //Make sure we tell the main process that we're ready to go!
                //Copy the testWorld to the gameWorld.
            }
        }
    }
}

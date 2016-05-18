using UnityEngine;
using UnityEngine.UI;

using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading;

using LitJson;

public class GameLoop : MonoBehaviour {
    IPAddress ipAddress;
    IPEndPoint listenerIPEP;
    Socket socket = new Socket ( AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
    Socket listenerSocket = new Socket ( AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);

    public GameObject playerPrefab;
    public GameObject asteroidPrefab;
    //GameObject []Players;
    //GameObject []Asteroids;
    
    ParallelUpdate threadedUpdate;

    public Text testText1;

    // Use this for initialization
    void Start () {
        testText1.text = "Startup Intialized, please stand by";
        ipAddress = IPAddress.Parse("127.0.0.1");
        IPEndPoint serverIPEP = new IPEndPoint(ipAddress, 9000);
        int listenerPort = requestPort(serverIPEP);

        listenerIPEP = new IPEndPoint( ipAddress, listenerPort );
        listenerSocket.Connect( listenerIPEP );
        Debug.Log( "Socket connected to the port received from server" );
        testText1.text = "Startup completed";

        byte[] message = new byte[8192];
        int bytesReceived = listenerSocket.Receive(message);
        string jsonString = Encoding.UTF8.GetString( message, 0, bytesReceived );
        testText1.text = jsonString;
        Debug.Log( jsonString );
        World gameWorld = JsonMapper.ToObject<World>(jsonString);

        threadedUpdate = new ParallelUpdate(gameWorld, listenerSocket);

        Thread oThread = new Thread(new ThreadStart(threadedUpdate.threadedUpdate));
        oThread.Start();
    }
	
	// Update is called once per frame
	void Update () {
        /*for ( int i = 0; i < threadedUpdate.gameWorld.Players.Length; i++ ) {
            threadedUpdate.gameWorld.Players[i].DrawMe( playerPrefab );
        }
        Debug.Log( threadedUpdate.gameWorld.Asteroids.Length );

        for ( int i = 0; i < threadedUpdate.gameWorld.Asteroids.Length; i++ ) {
            threadedUpdate.gameWorld.Asteroids[i].DrawMe( asteroidPrefab );
        }*/
        //Need a way to kill or remove old asteroids to make space for the new ones
        //Currently comparing with ID's, but they're never removed, just not sent again
    }

    int requestPort(IPEndPoint serverIPEP) {
        socket.Connect( serverIPEP );
        byte[] message = new byte[1024];
        int bytesReceived = socket.Receive(message);
        string jsonPort = Encoding.UTF8.GetString( message, 0, bytesReceived );

        //Close the active connection to the server so that we can create a new one with the port
        socket.Shutdown(SocketShutdown.Both);
        socket.Close();
        return int.Parse( jsonPort );
    }
}

[System.Serializable]
class ParallelUpdate {
    public World gameWorld { get; set; }
    private Socket socket { get; set; }

    public ParallelUpdate ( World world, Socket socket ) {
        gameWorld = world;
        this.socket = socket;
    }

    public void threadedUpdate () {
        while ( true ) {
            byte[] message = new byte[8192];
            int bytesReceived = socket.Receive(message);
            string jsonString = Encoding.UTF8.GetString( message, 0, bytesReceived );
            Debug.Log( jsonString );
            World jsonWorld = JsonMapper.ToObject<World> (jsonString);
            
            gameWorld = gameWorld.Update( jsonWorld );
            //Doesn't stop looping, asteroids had a function that crashed the thread.
        }
    }
}

[System.Serializable]
class World {
    public Player      []Players   { get; set; }
    public Asteroid    []Asteroids { get; set; }

    public World Update ( World newWorld ) {
        for (int i = 0; i < Asteroids.Length; i++) {
            if ( Asteroids[i].isAlive == false ) {
                Asteroids[i].ID = -1;
            }
        }

        foreach (Asteroid oldAsteroid in Asteroids ) {
            oldAsteroid.isAlive = false;
        }

        foreach (Player player in newWorld.Players ) {
            UpdatePlayers( player );
        }

        foreach //(int i = 0; i < newWorld.Asteroids.Length; i++) {
            (Asteroid asteroid in newWorld.Asteroids) {
            UpdateAsteroids( asteroid );
        }
        return this;
    }

    void UpdateAsteroids ( Asteroid newAsteroid ) {
        
        foreach(Asteroid oldAsteroid in Asteroids) {
            if (oldAsteroid.ID == newAsteroid.ID) {
                oldAsteroid.isAlive = true;
                oldAsteroid.X = newAsteroid.X;
                oldAsteroid.Y = oldAsteroid.Y;
                oldAsteroid.Phase = oldAsteroid.Phase;
                return;
            }
        }
        this.newAsteroid( newAsteroid );
    }

    void newAsteroid( Asteroid newAsteroid) {
        for (int i = 0; i < Asteroids.Length; i++ ) {
            if ( Asteroids[i].ID == -1) { //Change this, isAlive
                Asteroids[i] = newAsteroid;
                return;
            }
        }
        Debug.Log( "Array full" );
    }

    void UpdatePlayers (Player newPlayer) {
        foreach (Player oldPlayer in Players ) {
            if (oldPlayer.Name == newPlayer.Name ) {
                oldPlayer.XCord = newPlayer.XCord;
                oldPlayer.YCord = newPlayer.YCord;
                oldPlayer.Lives = newPlayer.Lives;
                return;
            }
        }
        this.newPlayer( newPlayer );
    }

    void newPlayer ( Player newPlayer ) {
        for (int i = 0; i < this.Players.Length; i++ ) {
            if (Players[i].Name == "") {
                Players[i] = newPlayer;
            }
        }
    }
}

[System.Serializable]
class playerObject : ScriptableObject {
    public GameObject playObj { get; set; }

    public void CreateMe( GameObject playPrefab, int XCord, int YCord ) {
        playObj = Instantiate( playPrefab );
        Rigidbody playerBody = playObj.GetComponent < Rigidbody >();
        playerBody.useGravity = false;
        playerBody.isKinematic = false;
        playObj.transform.position = new Vector3( XCord, YCord, 0 );
        playerBody.freezeRotation = true;
    }

    public void DrawMe(int XCord, int YCord) {
        //playerObject = new GameObject( Name );
        //Rigidbody playerBody = playerObject.AddComponent<Rigidbody>();
        //playerBody.useGravity = false;
        //playerObject.transform.position = new Vector3 ( XCord, YCord, 0 );
        //playerBody.freezeRotation = true;
        playObj.transform.position = new Vector3( XCord, YCord, 0 );
    }
}

[System.Serializable]
class Player {
    public playerObject playObj { get; set; }
    public string       Name    { get; set; }
    public int          XCord   { get; set; }
    public int          YCord   { get; set; }
    public int          Lives   { get; set; }
    public bool         drawn   { get; set; }

    public Player () {
        drawn = false;
    }

    public void DrawMe(GameObject playerPrefab) {
        if ( !drawn ) {
            //playObj = new playerObject( playerPrefab, XCord, YCord );
            playObj = ScriptableObject.CreateInstance<playerObject> ();
            playObj.CreateMe( playerPrefab, XCord, YCord );
            drawn = true;
        } else {
            playObj.DrawMe( XCord, YCord );
        }
    }
}

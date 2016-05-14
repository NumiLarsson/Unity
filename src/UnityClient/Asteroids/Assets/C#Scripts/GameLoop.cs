using UnityEngine;
using UnityEngine.UI;

using System.Collections;
using System.Net;
using System.Net.Sockets;
using System.Text;

using LitJson;

public class GameLoop : MonoBehaviour {
    IPAddress ipAddress;
    IPEndPoint listenerIPEP;
    Socket socket = new Socket ( AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
    Socket listenerSocket = new Socket ( AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);

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
    }
	
	// Update is called once per frame
	void Update () {
        byte[] message = new byte[1024];
        int bytesReceived = listenerSocket.Receive(message);
        string jsonString = Encoding.UTF8.GetString( message, 0, bytesReceived );
        testText1.text = jsonString;
        Debug.Log( jsonString );
    }

    int requestPort(IPEndPoint serverIPEP) {
        socket.Connect( serverIPEP );
        byte[] message = new byte[1024];
        int bytesReceived = socket.Receive(message);
        string jsonPort = Encoding.UTF8.GetString( message, 0, bytesReceived );

        //int listenerPort = JsonMapper.ToObject<int>(jsonPort);

        //Close the active connection to the server so that we can create a new one with the port
        socket.Shutdown(SocketShutdown.Both);
        socket.Close();
        return int.Parse( jsonPort );
    }
}

class World {
    int worldSize { get; set; }
}

class Player {
    int XCord { get; set; }
    int YCord { get; set; }
    int Lives { get; set; }
}

class Asteroid {
    int x { get; set; }
    int y { get; set; }
    int xStep { get; set; }
    int yStep { get; set; }
    int id { get; set; }
    int size { get; set; }
    int phase { get; set; }
}
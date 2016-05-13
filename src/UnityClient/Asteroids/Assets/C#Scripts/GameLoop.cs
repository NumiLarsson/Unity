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

    // Use this for initialization
    void Start () {
        ipAddress = IPAddress.Parse("127.0.0.1");
        IPEndPoint serverIPEP = new IPEndPoint(ipAddress, 9000);
        int listenerPort = requestPort(serverIPEP);

        Debug.Log( listenerPort );

        listenerIPEP = new IPEndPoint( ipAddress, listenerPort );

        socket.Connect( listenerIPEP );
        Debug.Log( "Socket connected to the port received from server" );
    }
	
	// Update is called once per frame
	void Update () {
	    //Lateer
	}

    int requestPort(IPEndPoint serverIPEP) {
        socket.Connect( serverIPEP );
        byte[] message = new byte[1024];
        int bytesReceived = socket.Receive(message);
        string jsonPort = Encoding.UTF8.GetString( message, 0, bytesReceived );

        int listenerPort = JsonMapper.ToObject<int>(jsonPort);

        //Close the active connection to the server so that we can create a new one with the port
        socket.Shutdown(SocketShutdown.Both);
        socket.Close();

        return listenerPort;
    }
}

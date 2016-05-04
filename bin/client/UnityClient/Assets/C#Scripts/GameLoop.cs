using UnityEngine;
using UnityEngine.UI;
using System.Collections;
using System;
using System.Net;
using System.Net.Sockets;
using System.Text;

public class GameLoop : MonoBehaviour {
	public Text testText;
	Socket socket = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);

	IPAddress ipAddress;
	IPEndPoint IPEP;

	// Use this for initialization
	void Start () {
		testText.text = "Start";

		ipAddress = IPAddress.Parse("127.0.0.1"); //Parse "localhost" to an IPaddress
		IPEP = new IPEndPoint(ipAddress, 9000);

		socket.Connect (IPEP);

		byte[] message = new byte[1024];
		int bytesReceived = socket.Receive(message);
		testText.text = Encoding.UTF8.GetString(message, 0, bytesReceived);
	}
	
	// Update is called once per frame
	void Update () {
			byte[] message = new byte[1024];
			int bytesReceived = socket.Receive (message);
			testText.text = Encoding.UTF8.GetString (message, 0, bytesReceived);
	}
}

using UnityEngine;
using UnityEngine.UI;

using System.Collections;
using System.Net;
using System.Net.Sockets;
using System.Text;

using LitJson;

public class GameLoopv2 : MonoBehaviour {
	public Text testText1;
	
	private Socket socket = new Socket ( AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
	
	private string jsonString;
	
	// Use this for initialization
	void Start () {
		Application.targetFrameRate = 5;
		
		IPAddress ipAddress = IPAddress.Parse("127.0.0.1");
		IPEndPoint IPEP = new IPEndPoint(ipAddress, 9000);
		
		socket.Connect( IPEP );
		
		testText1.text = "Connection has been established";
		
		byte[] message = new byte[4096];
		int bytesReceived = socket.Receive(message);
		
		jsonString = Encoding.UTF8.GetString( message, 0, bytesReceived );
		
		World jsonWorld = JsonMapper.ToObject<World>(jsonString);
		
		Vector2 vector2 = new Vector2((float)jsonWorld.Players[0].XCord, (float)jsonWorld.Players[0].YCord);
		Debug.Log (vector2);
		//ship.transform.position = vector2;
		//0'orna ska bytas ut mot kordenater för skeppet 
		
	}
	
	// Update is called once per frame
	void Update () {
		
	}
}

class Player {
	public int XCord { get; set; }
	public int YCord { get; set; }
	public int Lives { get; set; }
	
	public Player () {
		XCord = 0;
		YCord = 0;
		Lives = 0;
	}
	
	public string stringify () {
		return ( XCord.ToString() + YCord.ToString() + Lives.ToString() );
	}
}

class Asteroid {
	public int XCord { get; set; }
	public int YCord { get; set; }
	public int Stage { get; set; }
	
	public Asteroid () {
	}
}

class World {
	public Player []Players { get; set; }
	public Asteroid []Asteroids { get; set; }
	
	public World() {
	}
	
	public string stringify() {
		return Players[0].stringify();
	}
}
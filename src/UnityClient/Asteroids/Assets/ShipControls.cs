using UnityEngine;

using System.Net.Sockets;
using System.Text;

using LitJson;
using System.Collections;

public class ShipControls : ScriptableObject {

    public Socket listenerSocket;

    // Use this for initialization
    void Start () {
        Debug.Log( "ShipControls started" );
	}

	void Update() {
        if (listenerSocket == null ) {
            Debug.Log( "Socket is null" );
            return;
        }
		if (Input.GetKey (KeyCode.A)) {
            //West
            Debug.Log( "West" );
            playerMessage message = new playerMessage("MoveMe", "West");
            string jsonMessage = JsonUtility.ToJson( message );
            byte[] msg = Encoding.UTF8.GetBytes(jsonMessage);
            listenerSocket.Send( msg );
            Debug.Log( "Sent to server" );
        } else if (Input.GetKey (KeyCode.D)) {
            //East
            Debug.Log( "East" );

        } else if (Input.GetKey (KeyCode.W )){
            //North
            Debug.Log( "North" );
        } else if (Input.GetKey (KeyCode.S)) {
            //South
            Debug.Log( "South" );
        }
	}
}

using UnityEngine;
using System.Collections;

public class ShipControls : ScriptableObject {
    public int X;
    public int Y;
	// Use this for initialization
	void Start () {
	}

	void FixedUpdate() {
		if (Input.GetKey (KeyCode.A)) {
            //West
            Debug.Log( "West" );
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

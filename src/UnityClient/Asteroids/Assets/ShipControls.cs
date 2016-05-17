using UnityEngine;
using System.Collections;

public class ShipControls : MonoBehaviour {

	public Rigidbody2D ship;
	public float RotationSpeed;
	public float ThrustForce;

	// Use this for initialization
	void Start () {
		
	}

	void FixedUpdate() {
		ship = GetComponent<Rigidbody2D> ();

		if (Input.GetKey (KeyCode.A)) {
			//rotate ship left
			ship.angularVelocity = RotationSpeed;
		
		} else if (Input.GetKey (KeyCode.D)) {
			//rotate ship right
			ship.angularVelocity = -RotationSpeed;

		} else {
			ship.angularVelocity = 0f; 
		}

		if (Input.GetKey (KeyCode.W)) {
			//Move forward
			ship.AddForce(transform.up * ThrustForce);
		}
	}
}

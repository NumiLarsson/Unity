using UnityEngine;
using System.Collections;

public class ShipControls : MonoBehaviour {

	public Rigidbody2D ship;
	public float RotationSpeed;
	public float ThrustForce;
	public GameObject RocketPrefab;
	public GameObject ShipExplosion;	
	//public GameObject ShipFlames;
	public ParticleSystem ShipFlames; 

	void Update() {
		if (Input.GetKeyDown (KeyCode.Space)) {
			Instantiate (RocketPrefab, transform.position, transform.rotation);		
		}	
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
			ship.AddForce (transform.up * ThrustForce);
			//Add flames
			ShipFlames.Play (); 
		} else {
			ShipFlames.Stop ();
		}
	}

	void OnTriggerEnter2D(Collider2D collider) {
		if (collider.gameObject.tag == "Asteroid") {
			Instantiate (ShipExplosion, transform.position, new Quaternion ());
			Destroy (gameObject);

		}
	}
}
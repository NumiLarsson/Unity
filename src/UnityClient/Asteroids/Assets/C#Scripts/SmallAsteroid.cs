using UnityEngine;
using System.Collections;

public class SmallAsteroid : MonoBehaviour {
	
	public float MinTorque;
	public float MaxTorque;
	public float MinForce;
	public float MaxForce;

	public GameObject Explosion;
	public Rigidbody2D smallAsteroid;

	void Start () {
		smallAsteroid = GetComponent<Rigidbody2D> ();
		
		float magnitude = Random.Range (MinForce, MaxForce);
		float x = Random.Range (-1f, 1f);
		float y = Random.Range (-1f, 1f);
		
		smallAsteroid.AddForce(magnitude * new Vector2 (x, y));
		
		float torque = Random.Range (MinTorque, MaxTorque);
		smallAsteroid.AddTorque(torque);
	}
	
	void OnTriggerEnter2D(Collider2D collider) {
		if (collider.tag == "ShipRocket") {
			Instantiate (Explosion, transform.position, new Quaternion ());
			Destroy (gameObject);
			Destroy (collider.gameObject);
		}
	}
}
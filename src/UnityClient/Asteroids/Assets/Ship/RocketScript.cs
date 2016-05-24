using UnityEngine;
using System.Collections;

public class RocketScript : MonoBehaviour {
	
	public Rigidbody2D Rocket;
	public float RocketLife;

	private float life;	
	public float RocketForce; 

	void Start () {
		life = RocketLife;
		Rocket = GetComponent<Rigidbody2D> ();
		Rocket.AddForce (transform.up * RocketForce);
	}

	void Update () {
		life -= Time.deltaTime;
		if (life < 0) Destroy (gameObject);
	}
}

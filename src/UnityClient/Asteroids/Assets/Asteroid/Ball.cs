using UnityEngine;
using System.Collections;


public class Ball : MonoBehaviour {


	float sx;
	float sy;

	public Rigidbody asteroid;
	// Use this for initialization
	void Start () {
		sx = Random.Range (0, 2) == 0 ? -1 : 1;
		sy = Random.Range (0, 2) == 0 ? -1 : 1;

		asteroid = GetComponent<Rigidbody> ();

		asteroid.velocity = new Vector3 (Random.Range (2, 6) * sx, Random.Range (2, 6) * sy, 0);

	}
	
	
	// Update is called once per frame
	void Update () {
	
	}
}

using UnityEngine;

/// <summary>
/// Used to create flames behind the ship so that you can tell where you are when you're dead.
/// These are stupid flames, they spawn relative to the ship and move with it, so be carefull using it 
/// anywhere else.
/// </summary>
public class RocketScript : MonoBehaviour {
	
	public Rigidbody2D Rocket;
	public float RocketLife;

	private float life;	
	public float RocketForce; 

    /// <summary>
    /// Start is called once by Unity when the object is spawned.
    /// </summary>
	void Start () {
		life = RocketLife;
		Rocket = GetComponent<Rigidbody2D> ();
		Rocket.AddForce (transform.up * RocketForce);
	}

    /// <summary>
    /// Update is called every frame.
    /// Life is set in the Unity editor manually.
    /// </summary>
	void Update () {
		life -= Time.deltaTime;
		if (life < 0) Destroy (gameObject);
	}
}

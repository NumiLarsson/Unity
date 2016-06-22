using UnityEngine;

/// <summary>
/// ExplosionScript is attached to the explosion (TBI), and is used to destroy it after 4 seconds, to prevent infinite explosions.
/// </summary>
public class ExplosionScript : MonoBehaviour {
	
	void Start () {
		Destroy (gameObject, 4f);
	}
}

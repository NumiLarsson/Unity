using UnityEngine;
using System.Collections;

public class ExplosionScript : MonoBehaviour {
	
	void Start () {
		Destroy (gameObject, 4f);
	}
}

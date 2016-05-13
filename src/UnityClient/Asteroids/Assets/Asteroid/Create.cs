using UnityEngine;
using System.Collections;

public class Create : MonoBehaviour {

	public GameObject prefab;
	GameObject prefabClone;

	void Update() {
		if (Input.GetKeyDown ("r")) {
			prefabClone = Instantiate (prefab, transform.position, Quaternion.identity) as GameObject;
			Destroy (prefabClone, 3); //Destroys object after 3 seconds
		}
	}
}
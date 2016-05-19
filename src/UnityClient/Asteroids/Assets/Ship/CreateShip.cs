using UnityEngine;
using System.Collections;

public class CreateShip : MonoBehaviour {
	
	public GameObject prefabShip;
	GameObject prefabShipClone;
	
	void Update() {
		if (Input.GetKeyDown ("q")) {
			prefabShipClone = Instantiate (prefabShip, transform.position, Quaternion.identity) as GameObject;
			Destroy (prefabShipClone, 200);
		}
	}
}
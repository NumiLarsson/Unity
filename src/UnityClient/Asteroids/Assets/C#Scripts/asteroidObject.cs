using UnityEngine;
using System.Collections;

public class asteroidObject : MonoBehaviour {
    public GameObject astObj { get; set; }
    int notDrawn = 0;

    public void CreateMe ( GameObject astPrefab, int XCord, int YCord ) {
        astObj = Instantiate( astPrefab );

        Rigidbody playerBody = astObj.GetComponent < Rigidbody >();
        playerBody.useGravity = false;
        playerBody.isKinematic = false;

        astObj.transform.position = new Vector3( XCord, YCord, 0 );
        playerBody.freezeRotation = true;
    }

    public void DrawMe ( int XCord, int YCord ) {
        astObj.transform.position = new Vector3( XCord, YCord, 0 );
        notDrawn = 0;
    }

    public void DestroyMe () {
        //Not actually destroying anything
        Destroy( this );
    }

    // Use this for initialization
    void Start () {
        if (notDrawn >= 10) {
            Destroy( this );
        }
        notDrawn++;
	}
	
	// Update is called once per frame
	void Update () {
	
	}
}

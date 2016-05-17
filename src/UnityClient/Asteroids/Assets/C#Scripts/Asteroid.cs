using UnityEngine;
using System.Collections;
[System.Serializable]
public class Asteroid {
    //public asteroidObject astObj { get; set; }
    public GameObject astObj { get; set; }
    public int X { get; set; }
    public int Y { get; set; }
    public int ID { get; set; }
    public int Phase { get; set; }
    public bool drawn { get; set; }
    public bool isAlive { get; set; }

    public void DrawMe ( GameObject asteroidPrefab ) {
        if ( !drawn ) {
            //astObj = new asteroidObject( asteroidPrefab, X, Y );
            //astObj = ScriptableObject.CreateInstance<asteroidObject>();
            //astObj.CreateMe( asteroidPrefab, X, Y );
            //Goal for today, 2016-05-17
            //Maybe done?
            drawn = true;
            astObj = Object.Instantiate( asteroidPrefab );
            //GameObject asteroidObject = new GameObject( "Clone" );
            Rigidbody asteroidBody = astObj.GetComponent<Rigidbody>();
            asteroidBody.useGravity = false;
            asteroidBody.isKinematic = false;
            asteroidBody.transform.position = new Vector3 ( (X/ 15.38f) - 3.25f, (Y/ 15.38f) - 3.25f, 0 );
            asteroidBody.freezeRotation = true;
            //asteroidBody.velocity = new Vector3( 0, 0, 0 );
            asteroidBody.velocity = Vector3.zero;
            //Object.Destroy( astObj, 0.5f );
            isAlive = true;
        } else {
            if ( ID == -1 ) {
                Object.Destroy( astObj );
            } else if (astObj != null ) {
                Rigidbody asteroidBody = astObj.GetComponent<Rigidbody>();
                asteroidBody.MovePosition( new Vector3(( X / 15.38f ) - 3.25f, ( Y / 15.38f ) - 3.25f, 0) );
                asteroidBody.velocity = Vector3.zero;
            }
        }
    }
    // Use this for initialization
    void Start () {
	
	}
	
	// Update is called once per frame
	void Update () {
        if ( !isAlive && ID == -1 ) {
            Object.Destroy( astObj );
        }
    }
}

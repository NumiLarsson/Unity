using UnityEngine;
using System.Collections;
[System.Serializable]
public class Asteroid : MonoBehaviour {
    //public asteroidObject astObj { get; set; }
    public Rigidbody body { get; set; }
    public int X { get; set; }
    public int Y { get; set; }
    public int ID { get; set; }
    public int Phase { get; set; }
    public bool isAlive { get; set; }

    // Use this for initialization
    void Start () {
        //astObj = Object.Instantiate( asteroidPrefab );
        //GameObject asteroidObject = new GameObject( "Clone" );
        body = this.GetComponent<Rigidbody>();
        body.useGravity = false;
        body.isKinematic = false;
        body.transform.position = new Vector3( ( X / 15.38f ) - 3.25f, ( Y / 15.38f ) - 3.25f, 0 );
        //asteroidBody.transform.position = new Vector3( ( X / 15.38f ) - 3.25f, ( Y / 15.38f ) - 3.25f, 0 );
        body.freezeRotation = true;
        //asteroidBody.velocity = new Vector3( 0, 0, 0 );
        body.velocity = Vector3.zero;
        //Object.Destroy( astObj, 0.5f );
        isAlive = true;
    }
	
	// Update is called once per frame
	void Update () {
        if ( !isAlive && ID == -1 ) {
            Object.Destroy( this );
        } else {
            body.transform.position = new Vector3( ( X / 15.38f - 3.25f), ( Y / 15.38f - 3.25f ) );
        }
    }
}

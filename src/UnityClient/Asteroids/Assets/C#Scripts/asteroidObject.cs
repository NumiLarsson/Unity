using UnityEngine;
using System.Collections;

[System.Serializable]
public class AsteroidObject : MonoBehaviour {
    //public asteroidObject astObj { get; set; }
    public Rigidbody body { get; set; }
    public int ID { get; set; }
    public int X { get; set; }
    public int Y { get; set; }
    public int Phase { get; set; }
    public int framesSinceDrawn;
    public bool drawnLastFrame;

    public void UpdateMe (Asteroid newAsteroid) {
        this.X = newAsteroid.X;
        this.Y = newAsteroid.Y;
        this.Phase = newAsteroid.Phase;
        this.framesSinceDrawn = 0;
        this.drawnLastFrame = true;
    }

    public void FlagFalse () {
        drawnLastFrame = false;
    }

    public void KillSelf() {
        Destroy(this);
    }

    public void InitMe ( int newID ) {
        ID = newID;
        this.name = ID.ToString();
        framesSinceDrawn = 0;
        drawnLastFrame = true;
    }

    // Use this for initialization
    void Start () {
        body = this.GetComponent<Rigidbody>();
        body.useGravity = false;
        body.isKinematic = false;
        body.transform.position = new Vector3( ( X  / 2 ) - 3.25f, ( Y / 15.38f ) - 3.25f, 0 );
        //asteroidBody.transform.position = new Vector3( X, Y, 0 ); //Set position
        body.freezeRotation = true;
        body.velocity = Vector3.zero; //Stop moving!
        framesSinceDrawn = 0;
    }
	
	// Update is called once per frame
	void Update () {
        framesSinceDrawn++;
        if ( framesSinceDrawn > 10) {
            Destroy( this.gameObject );
            Destroy( body );
            Destroy( this );
        } else if (this.gameObject.name == "-1") {
            Destroy( this.gameObject );
        } else {
            body.transform.position = new Vector3( ( X / 2 - 3.25f ), ( Y / 15.38f - 3.25f ) );
        }
    }
}

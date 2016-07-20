using UnityEngine;
using System.Collections;

public class PlayerObject : MonoBehaviour {
    public Rigidbody2D body;
    public string Name { get; set; }
    public int XCord { get; set; }
    public int YCord { get; set; }
    public int Lives { get; set; }
    public int framesSinceDrawn = 0;

    public void UpdateMe (Player newPlayer) {
        this.XCord = newPlayer.XCord;
        this.YCord = newPlayer.YCord;
        this.Lives = newPlayer.Lives;
        this.framesSinceDrawn = 0;
    }

    public void InitMe (string playerName) {
        this.Name = playerName;
        this.name = playerName;
    }

    // Use this for initialization
    void Start () {
        body = this.GetComponent<Rigidbody2D>();
        body.isKinematic = false;
        body.transform.position = new Vector3( ( XCord ), ( YCord ) );
        //asteroidBody.transform.position = new Vector3( X, Y, 0 ); //Set position
        body.freezeRotation = true;
        body.velocity = Vector3.zero; //Stop moving!
    }

    // Update is called once per frame
    void Update () {
        framesSinceDrawn++;
        if (framesSinceDrawn > 10) {
            Destroy( this.gameObject );
            Destroy( body );
            Destroy( this ); 
        } else {
            //body.transform.position = new Vector3( XCord, YCord );
            body.transform.position = new Vector3( ( XCord ), ( YCord ) );
        }
    }
}

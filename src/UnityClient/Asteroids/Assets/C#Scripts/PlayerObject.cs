using UnityEngine;
using System.Collections;

public class PlayerObject : MonoBehaviour {
    public Rigidbody2D shipBody { get; set; }
    //public ShipControls shipScript { get; set; }
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
        //shipScript.X = XCord;
        //shipScript.Y = YCord;
    }

    public void InitMe (string playerName) {
        this.Name = playerName;
        this.name = playerName;
    }

    // Use this for initialization
    void Start () {
        //shipScript = ScriptableObject.CreateInstance("ShipControls") as ShipControls;
        shipBody = this.GetComponent<Rigidbody2D>();
        shipBody.isKinematic = false;
        shipBody.transform.position = new Vector3( ( XCord ), ( YCord ) );
        //shipScript.X = XCord;
        //shipScript.Y = YCord;
        shipBody.freezeRotation = true;
        shipBody.velocity = Vector3.zero;
    }

    // Update is called once per frame
    void Update () {
        framesSinceDrawn++;
        if (framesSinceDrawn > 10) {
            //Destroy( shipScript );
            Destroy( this.gameObject );
            Destroy( shipBody );
            Destroy( this ); 
        } else {
            if ( Lives > 0 ) {
                shipBody.transform.position = new Vector3( ( XCord ), ( YCord ) );
            }
        }
    }
}

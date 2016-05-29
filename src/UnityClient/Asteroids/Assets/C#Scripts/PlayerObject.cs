using UnityEngine;
using UnityEngine.UI;

public class PlayerObject : MonoBehaviour {
    public Rigidbody2D  shipBody    { get; set; }
    public string       Name        { get; set; }
    public int          ID          { get; set; }
    public int          XCord       { get; set; }
    public int          YCord       { get; set; }
    public int          Lives       { get; set; }
    public int          Rotation    { get; set; }
    public bool         Alive       { get; set; }
    public int          Points      { get; set; }

    public ParticleSystem ShipFlames;

    public int framesSinceDrawn     = 0;
    private static float scaling    = 1.25f;

    public void UpdateMe ( Player newPlayer ) {
        if (XCord == newPlayer.X && this.YCord == newPlayer.Y) {
            ShipFlames.Stop();
        } else {
            ShipFlames.Play();
        }
        this.XCord = newPlayer.X;
        this.YCord = newPlayer.Y;
        this.Lives = newPlayer.Lives;
        this.framesSinceDrawn = 0;
        this.Rotation = newPlayer.Rotation;
        this.Alive = newPlayer.Alive;
        this.Points = newPlayer.Points;
    }

    public void InitializeMe ( string playerName ) {
        this.Name = playerName;
        this.name = playerName;
    }

    // Use this for initialization
    void Start () {
        shipBody = this.GetComponent<Rigidbody2D>();
        shipBody.isKinematic = false;
        shipBody.transform.position = new Vector3( ( XCord / (scaling) ), ( YCord / (scaling)) );
        shipBody.transform.rotation = Quaternion.AngleAxis( 90 * Rotation, new Vector3(0,0));
        shipBody.velocity = Vector3.zero;
    }

    // Update is called once per frame
    void Update () {
        framesSinceDrawn++;
        if (framesSinceDrawn > 10) {
            Destroy( this.gameObject );
            Destroy( shipBody );
            Destroy( this ); 
        } else {
            shipBody.transform.rotation = Quaternion.Euler( new Vector3( 0, 0, 90 * Rotation ) );
            shipBody.transform.position = new Vector3( ( XCord/(scaling)  - 160f), ( YCord/(scaling) - 90f ));
            if (this.Alive) {
                //Drag this out to gain performance
                Renderer tempRend = this.GetComponent<Renderer>();
                //Drag this out to gain performance
                tempRend.enabled = true;
            } else {
                Renderer tempRend = this.GetComponent<Renderer>();
                tempRend.enabled = false;
            }
        }
    }
}

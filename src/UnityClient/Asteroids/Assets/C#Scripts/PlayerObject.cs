using UnityEngine;
using UnityEngine.UI;

/// <summary>
/// PlayerObject -is the local Unity GameObjects used to represent our players.
/// </summary>
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
    // A very ugly way to make sure that all gameobjects die if they're out of date.
    private static float scaling    = 1.25f; 
    //Scales the gamearea down a slight bit so that Unity is a tad bit faster.

    /// <summary>
    /// Uses a Player (the class) to update the local variables of relevance, also turns of or 
    /// on the flames.
    /// </summary>
    /// <param name="newPlayer">
    /// Sends a player to use to update all the coordinates, lives, rotations etc from.
    /// </param>
    public void UpdateMe ( Player newPlayer ) {
        if (XCord == newPlayer.X && this.YCord == newPlayer.Y) {
            ShipFlames.Stop(); //Player didn't move last frame, so stop playing the frames.
        } else {
            ShipFlames.Play(); //Player did move, so turn on the rockets!
        }
        this.XCord = newPlayer.X;
        this.YCord = newPlayer.Y;
        this.Lives = newPlayer.Lives;
        this.framesSinceDrawn = 0;
        this.Rotation = newPlayer.Rotation;
        this.Alive = newPlayer.Alive;
        this.Points = newPlayer.Points;
    }

    /// <summary>
    /// Since Unity won't call an "init" function, we do that ourselves.
    /// </summary>
    /// <param name="playerName"></param>
    public void InitializeMe ( string playerName ) {
        this.Name = playerName;
        this.name = playerName;
    }

    /// <summary>
    /// Start is called when the object is spawned, but you can't manually call it, so can't be used 
    /// as an init function.
    /// </summary>
    void Start () {
        shipBody = this.GetComponent<Rigidbody2D>();
        shipBody.isKinematic = false;
        shipBody.transform.position = new Vector3( ( XCord / (scaling) ), ( YCord / (scaling)) );
        shipBody.transform.rotation = Quaternion.AngleAxis( 90 * Rotation, new Vector3(0,0));
        shipBody.velocity = Vector3.zero;
    }

    /// <summary>
    /// Update is called once per frame by UnityEngine, this is where all our updates of positions and such
    /// should be done.
    /// </summary>
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

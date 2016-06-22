/// <summary>
/// playerMessage is the struct we symbolize all data with when we send input from the client to the 
/// server. It's also part of the registering the name´.
/// </summary>
public class playerMessage {
    public string Action    { get; set; }
    public string Value     { get; set; }

    public playerMessage(string Action, string Value) {
        this.Action = Action;
        this.Value = Value;
    }
}

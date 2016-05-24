using UnityEngine;
using System.Collections;

public class playerMessage {
    public string Action { get; set; }
    public string Value { get; set; }

    public playerMessage(string Action, string Value) {
        this.Action = Action;
        this.Value = Value;
    }
}

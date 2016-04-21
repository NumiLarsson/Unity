using System;
using System.Net;
using System.Net.Sockets;
using System.Text;

public class SynchronousSocketClient
{

    public static IPEndPoint GetIPEndPoint(int port) {
        IPAddress ipAddress = IPAddress.Parse("127.0.0.1"); //Parse "localhost" to an IPaddress
        return new IPEndPoint(ipAddress, port);
    }

    public static Socket New_Streaming_IP_TCP_Socket()
    {
        return new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
    }

    //REturns the number of bytes sent, Automatically converts everything to UTF-8 standard (that GO uses).
    public static int SendStringTo(String message, Socket socket)
    {
        byte[] msg = Encoding.UTF8.GetBytes("This is a test\n    hejsan");
        return socket.Send(msg);
    }

    //Returns an UTF8 string, passed to it through the Socket. The message is also kept in an array, if the 
    //caller of the function wants to access it again. Buffersize is max 1024 bytes.
    public static String SocketReceiveString(Socket socket)
    {
        byte[] message = new byte[1024];
        int bytesReceived = socket.Receive(message);
        return Encoding.UTF8.GetString(message, 0, bytesReceived);
    }

    //Receives a message from socket and place it in message before returning message.
    public static byte[] SocketReceiveBytes(byte[] message, Socket socket)
    {
        int bytesReceived = socket.Receive(message);
        return message;
    }

    public static void StartClient()
    {
        // Data buffer for incoming data.
        byte[] bytes = new byte[1024];

        // Connect to a remote device.
        try
        {
            /*IPHostEntry ipHostInfo = Dns.Resolve(Dns.GetHostName());
            //IPAddress ipAddress = IPAddress.Parse("127.0.0.1");
            //IPEndPoint remoteEP = new IPEndPoint(ipAddress, 9000); */

            //Replaced with a functioncall
            IPEndPoint remoteEP = GetIPEndPoint(9000);
            //Create / Find an IP address that we can use to establish a connection.

            // Create a TCP/IP  socket.
            Socket sender = New_Streaming_IP_TCP_Socket();

            // Connect the socket to the remote endpoint. Catch any errors.
            try
            {
                sender.Connect(remoteEP);
  
                Console.WriteLine("Socket connected to {0}",
                    sender.RemoteEndPoint.ToString());

                // Encode the data string into a byte array.
                String msg = "This is a test\n hejsan";

                // Send the data through the socket.
                // int bytesSent = SendStringTo(msg, sender);

                // Receive the response from the remote device.
                for (int i = 0; i < 10; i++) {
                    Console.WriteLine("Echoed test = {0}",
                        SocketReceiveString(sender));
                }
                // Release the socket.
                sender.Shutdown(SocketShutdown.Both);
                sender.Close();
            }
            catch (ArgumentNullException ane)
            {
                Console.WriteLine("ArgumentNullException : {0}", ane.ToString());
            }
            catch (SocketException se)
            {
                Console.WriteLine("SocketException : {0}", se.ToString());
            }
            catch (Exception e)
            {
                Console.WriteLine("Unexpected exception : {0}", e.ToString());
            }

        }
        catch (Exception e)
        {
            Console.WriteLine(e.ToString());
        }
        Console.WriteLine("Program completed succesfully.");
        Console.ReadLine();
    }



    public static int Main(String[] args)
    {
        StartClient();
        return 0;
    }
}
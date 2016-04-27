using System;
using System.Net;
using System.Net.Sockets;
using System.Text;

public class SynchronousSocketClient
{
    /// <summary>
    ///  Creates an IPEndPoint (What you connect sockets to etc) using an ip address formateted as a string and a port formatted as an int
    /// </summary>
    /// <param name="Address">IP Address in standard string format</param>
    /// <param name="port">Port to use with the IP address</param>
    /// <returns>THe created IPEndPoint</returns>
    public static IPEndPoint GetIPEndPoint(String Address, int port) {
        IPAddress ipAddress = IPAddress.Parse(Address); //Parse Address string to an IPaddress
        return new IPEndPoint(ipAddress, port);
    }

    /// <summary>
    ///  Creates a new two-way TCP socket.
    /// </summary>
    /// <returns>A standard TCP socket</returns>
    public static Socket New_Streaming_IP_TCP_Socket()
    {
        return new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
    }

    /// <summary>
    /// Sends a string through a tcp socket.
    /// </summary>
    /// <param name="message">standard string containing the message to be sent</param>
    /// <param name="socket">standard TCP socket.</param>
    public static void SendStringTo(String message, Socket socket)
    {
        byte[] msg = Encoding.UTF8.GetBytes(message);
        socket.Send(msg);
    }

    /// <summary>
    /// Reads up to 1024 bytes from sockets then formats those to an UTF8 string.
    /// </summary>
    /// <param name="socket">Standard TCP socket</param>
    /// <returns>UTF8 formatted standard string</returns>
    public static String SocketReceiveString(Socket socket)
    {
        byte[] message = new byte[1024];
        int bytesReceived = socket.Receive(message);
        return Encoding.UTF8.GetString(message, 0, bytesReceived);
    }

   


    
    /// <summary>
    /// Reads socket for a string, then converts it in to a uint16 to be used as a port.
    /// </summary>
    /// <param name="socket">Standard TCP socket</param>
    /// <returns>A UInt16 meant to be used as a port to connect to the server with</returns>
    public static UInt16 ReceivePort(Socket socket)
    {
        byte[] portBytes = new byte[1024];
        string portString = SocketReceiveString(socket);
        UInt16 portInt = Convert.ToUInt16(portString);
        return portInt;
    }

    /// <summary>
    /// Requests a port from the server, then returns that specified port as an UInt16
    /// </summary>
    /// <param name="ipAddress">standard IP address in a string</param>
    /// <param name="port">port in an int</param>
    /// <returns></returns>
    public static UInt16 RequestPort (String ipAddress, int port)
    {
        IPEndPoint remoteEP = GetIPEndPoint(ipAddress, port);
        UInt16 UIntPort = 0;

        try
        {
            Socket sender = New_Streaming_IP_TCP_Socket();

            try
            {
                sender.Connect(remoteEP);

                Console.WriteLine("Socket connected to {0}",
                    sender.RemoteEndPoint.ToString());

                UIntPort = ReceivePort(sender);

                sender.Shutdown(SocketShutdown.Both);
                sender.Close();

                return UIntPort;
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



        return UIntPort;
    }

    /// <summary>
    ///     Connects to a listener on the server using the specified port, then returns the active socket.
    /// </summary>
    /// <param name="ipAddress">Standard IPv4 Address in a string</param>
    /// <param name="port">Standard networkport</param>
    /// <returns></returns>
    public static Socket ConnectToListener(String ipAddress, UInt16 port)
    {
        IPEndPoint remoteEP = GetIPEndPoint(ipAddress, port);
        UInt16 UIntPort;
        Socket socket = null;
        try
        {
            //Socket
            socket = New_Streaming_IP_TCP_Socket();

            try
            {
                //Create a connection with remote end point.
                socket.Connect(remoteEP);

                Console.WriteLine("Socket connected to {0}",
                    socket.RemoteEndPoint.ToString());

                String name = "Anton\n";

                SendStringTo( name, socket );
                return socket;
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
        return socket;
    }

    public static int Main(String[] args)
    {
        UInt16 port = RequestPort("127.0.0.1", 9000);
        Console.Write(port);
        Console.WriteLine( " press enter to accept" );
        Console.ReadLine();
        Socket socket = ConnectToListener("127.0.0.1", port);

        for (int i = 0; i < 10; i++)
        {
            String input = Console.ReadLine();
            input = input + "\n ";
            SendStringTo(input, socket);
            Console.WriteLine(SocketReceiveString(socket));
        }
        socket.Shutdown(SocketShutdown.Both);
        socket.Close();
        return 0;
    }
}
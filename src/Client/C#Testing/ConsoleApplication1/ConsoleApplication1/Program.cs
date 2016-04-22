using System;
using System.Net;
using System.Net.Sockets;
using System.Text;

public class SynchronousSocketClient
{
    public static IPEndPoint GetIPEndPoint(String Address, int port) {
        IPAddress ipAddress = IPAddress.Parse(Address); //Parse Address string to an IPaddress
        return new IPEndPoint(ipAddress, port);
    }

    public static Socket New_Streaming_IP_TCP_Socket()
    {
        return new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
    }

    //REturns the number of bytes sent, Automatically converts everything to UTF-8 standard (that GO uses).
    public static void SendStringTo(String message, Socket socket)
    {
        byte[] msg = Encoding.UTF8.GetBytes("This is a test\n    hejsan");
        socket.Send(msg);
    }

    //Returns an UTF8 string, passed to it through the Socket. The message is also kept in an array, if the 
    //caller of the function wants to access it again. Buffersize is max 1024 bytes.
    public static String SocketReceiveString(Socket socket)
    {
        byte[] message = new byte[1024];
        int bytesReceived = socket.Receive(message);
        return Encoding.UTF8.GetString(message, 0, bytesReceived);
    }

    //Receives a message from socket and place it in message before returning the amount of bytes received.
    public static int SocketReceiveBytes(byte[] message, Socket socket)
    {
        int bytesReceived = socket.Receive(message);
        return bytesReceived;
    }

   


    //Connect to server, receive and return uint16

    public static UInt16 ReceivePort(Socket socket)
    {
        byte[] portBytes = new byte[2];
        int bytesReceived = socket.Receive(portBytes);
        UInt16 portInt = BitConverter.ToUInt16(portBytes, 0);
        return portInt;
    }

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

                String name = "Anton";

                SendStringTo(name, socket);
                Console.WriteLine("Message received from server: {0}", SocketReceiveString(socket));
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
        Console.WriteLine(port);
        Console.ReadLine();
        Socket socket = ConnectToListener("127.0.0.1", port);
        Console.ReadLine();
        for (int i = 0; i < 10; i++)
        {
            String input = Console.ReadLine();
            SendStringTo(input, socket);
            Console.WriteLine(SocketReceiveString(socket));
        }
        socket.Shutdown(SocketShutdown.Both);
        socket.Close();
        return 0;
    }
}
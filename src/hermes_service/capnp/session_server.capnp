@0xc4b468a2826bb79b;

using Go = import "/go.capnp";
$Go.package("hermes_service_capnp");
$Go.import("github.com/zalf-rpm/Hermes2Go/src/hermes_service/capnp/hermes_service_capnp");

# Interface for a session server, implemented by the server
interface SessionServer {
    newSession @0 (workdir :Text, resultCallback :Callback) -> (session :Session); # Create a new session
}

# Interface for a session, implemented by the server
interface Session {
    send @0 (runId :Text, params :List(Text)) -> (); # Send a sim data to the session
    close @1 () -> (); # Close the session
}

# Callback for the result of a sim data, implemented by the client
interface Callback {
    sendData @0 (runId :Text, outData :Text) -> (); #  result of a sim data
    done @1 (runId :Text) -> (); # end of a result data for the runId
}
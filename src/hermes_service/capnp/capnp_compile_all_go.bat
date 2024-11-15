
rem cd capnp
rem capnp compile -I.. -ogo:../gen/go/persistence persistent.capnp
rem cd ..

capnp compile -I. -ogo:./hermes_service_capnp session_server.capnp


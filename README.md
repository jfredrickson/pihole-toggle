# pihole-toggle

A simple HTTP server that toggles a client's blocking state on a Pi-hole. It essentially acts as a proxy to the Pi-hole API, providing a simplified interface.

This assumes that the Pi-hole is configured with a group with ID 0 (e.g., the `Default` group in Pi-hole) that has blocking enabled. If a client is in group 0, its DNS queries will go through the Pi-hole's blocking features. If a client is not in group 0, its DNS queries will not be filtered.

## Usage

Set the following environment variables:

- `PIHOLE_URL` - the base URL of the Pi-hole API (required, e.g., `http://127.0.0.1`)
- `PIHOLE_PASSWORD` - the password for the Pi-hole API (required, even if it's just an empty string)
- `PORT` - the port to listen on (optional, default: 8001)

Then just run the `pihole-toggle` binary. This will start a server on port 8001 (or the given port per the `PORT` environment variable) with two endpoints:

- `/on` - Enable blocking for the client
- `/off` - Disable blocking for the client

The Pi-hole client's IP address is automatically determined from the remote address of the request.

## Example

Say you have a client with IP address `192.168.1.100`, and you have this `pihole-toggle` server running at `http://pihole-toggle:8001`. The client can make the following requests:

```sh
# Enable blocking for the client
curl http://pihole-toggle:8001/on

# Disable blocking for the client
curl http://pihole-toggle:8001/off
```

The server will then relay the request to the Pi-hole API, telling it to move the client in or out of group 0, effectively enabling or disabling blocking for the client.

# Snapmaker CLI Tool

This tool shall at some point be able to perform the following operations
* upload a .gcode file
* get printer status
* get enclosure status
* execute Marlin commands

## Why Golang?
Because it's currently my favourite language.

## CLI

```
.\snapcmd.exe upload --api-token blabla --printer-ip 192.168.1.2 --discovery-timeout 15 test.gcode
```
* subcommand: operation to perform. Currently only `upload` is supported.
* `--api-token`: Mandatory parameter for now, an API token the Snapmaker already knows. TODO: Details of token creation and authentication yet to be researched.
* `--printer-ip`: Optional. IP address of printer, if omitted an auto-discovery will be done. TODO: auto-discovery not yet suited for multiple printers in the network.
* `--discovery-timeout`: Optional. Defaults to 5s.
* filepath: specific to `upload`. Has to be a `.gcode` file.

## Snapmaker API
To get the upload sequence right Wireshark sniffs where used to reverse engineer the network protocol of Snapmaker. For this tool the following API endpoints are used:

### API token
Each request is authenticated by an API token. However, the way how it's placed in the request differs for each endpoint.

TODO: It's not yet clear how this token is issued. Currently I'm using the API token used by Luban.

The token that used by Luban for a certain Snapmaker can be found in the following file in the section `server`:
`C:\Users\<username>\AppData\Roaming\snapmaker-luban\machine.json`
Using this token should work fine.

### Snapmaker discovery
Discovery is done by sending an UDP packet to the broadcast address of the local network interfaces. Usually the Snapmaker then responds with a short descriptive string. The IP address of the Snapmaker can either be read from that string or from the UDP connection used for the response packet.

### `/api/v1/connect`
As advertised, opens the connection between the tool and Snapmaker. It requires a `POST` request, with the token encoded as URL parameters in the body.

```
POST /api/v1/connect HTTP/1.1
Host: 192.168.188.130:8080
User-Agent: Go-http-client/1.1
Connection: close
Content-Length: 42
Content-Type: application/x-www-form-urlencoded
Accept-Encoding: gzip

token=aaaaaaaa-bbbbb-bbbb-cccc-dddddddddddd
```

### `/api/v1/status`
Unsurprisingly, returns a status JSON with the current printer status. Surprisingly, it seems to be necessary to repeatedly read the printer status, otherwise the upload command will return an `401 - Unauthorized` result stating that `the machine is not yet connected`.

The token is provided as URL parameter.

```
http://<snapmaker>:8080/api/v1/status?token=aaaaaaaa-bbbbb-bbbb-cccc-dddddddddddd
```

Luban adds another hex string as second URL parameter. It's purpose is not yet clear.

### `/api/v1/upload`
Requires a multipart-form request consisting of two parts:

#### token part
```
----------------------------447327606604133343229724
Content-Disposition: form-data; name="token"

aaaaaaaa-bbbbb-bbbb-cccc-dddddddddddd
```

It seems to be vitally important that this part doesn't have a `Content-Type` header. Otherwise the upload will result in `400 - Bad request`.

#### file part
```
----------------------------447327606604133343229724
Content-Disposition: form-data; name="file"; filename="test.gcode"
Content-Type: application/octet-stream

;FLAVOR:Marlin
;TIME:52383
;Filament used: 13.0944m
;Layer height: 0.08
;MINX:123.287
;MINY:140.882
;MINZ:0.15
;MAXX:196.712
;MAXY:209.939
;MAXZ:39.03
;Generated with Cura_SteamEngine 5.0.0
M82 ;absolute extrusion mode
M104 S220 ;Set Hotend Temperature
M140 S70 ;Set Bed Temperature
G28 ;home
G90 ;absolute positioning
G1 X-10 Y-10 F3000 ;Move to corner 
G1 Z0 F1800 ;Go to zero offset
M109 S220 ;Wait for Hotend Temperature
M190 S70 ;Wait for Bed Temperature
G92 E0 ;Zero set extruder position
G1 E20 F200 ;Feed filament to clear nozzle
...
```

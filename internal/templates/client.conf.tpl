[Interface]
PrivateKey = {{ .PriKeyClient }}
Address = {{ .AllowedIPs }}
DNS = 8.8.8.8, 1.1.1.1
MTU = {{ .MTU }}

[Peer]
PublicKey = {{ .PubKeyServer }}
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = {{ .Endpoint }}
PersistentKeepalive = 25

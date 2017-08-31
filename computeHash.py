import hashlib

secret = "shazbot"
body = """[{"title_code":"game2","profile_id":1,"namespace":"game2.players","service":"GATEWAY_SERVER_SERVICE","event_type":"PING","entity_code":"","amount":0,"message":"test ping 1504205866625"}]"""

computedHash = hashlib.sha256(secret + body).hexdigest()
print(computedHash)

# 25 users each w/ 40 passes thru the local friendly 25 requests -> 25000 calls
../../endpointTesterMac -urlPrefix http://localhost:8080/uaa -requestsConfigFileAt ./UAA-25-requests-localFriendly.txt -users 25 -passes 40
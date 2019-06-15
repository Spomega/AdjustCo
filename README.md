# AdjustCo
AdjustCo is a command line tool which makes http requests and prints the address of the request
along with the MD5 hash of the response.

##Installation
Use the go build command to build the package

```bash
go build main.go
```

##Usage
./main adjust.com # returns http://adjust.com b9e73272f11c3a60e0974e34815e06de

./main https://www.google.com https://www.twitter.com https://www.twilio.com https://www.facebook.com https://www.goal.com #return 
https://www.google.com ead9aad7e1cac1ce3fe423f9fe625be4
https://www.facebook.com 049dca545c85c66d012b1499f119f3f1
https://www.goal.com 1acc51aa6dbc28b4842860a1be895977
https://www.twilio.com d53d6684a1b172485b8e56939e3a26af
https://www.twitter.com cbfa42ef13aa270a1b80dddbce3874c0

./main --parallel 3 https://www.google.com https://www.twitter.com https://www.twilio.com https://www.facebook.com https://www.goal.com #return 
http://google.com 5eb6b1288b2e68bd5b4dc73a90521562
http://www.twilio.com 0c36083619236fbc21c3593e5c1e9551
http://www.twitter.com ce24e3da66ea069feb3c5b4e28583618
http://facebook.com 45466bb67d41a20f6418821ce35cb92b
http://goal.com 3d752f7fe4c7f399f420a61083d9d927


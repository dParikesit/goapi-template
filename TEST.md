# Test running code

To upload values to the database:

`curl --request POST --url https://<url_server>/send --header 'content-type: application/json' --data '{"value": "Message one"}'`

`curl --request POST --url https://<url_server>/send --header 'content-type: application/json' --data '{"value": "Message two"}'`

   ...

Now the url should count the number of messages stored in the database:

Hello, Docker! (2) 


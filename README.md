## OTUS Highload Architect Dialog Microservice

### How to use it 

Before start you need to check whether **docker** and **make** has been installed.

1.  Clone repository

```plaintext
git clone https://github.com/enemis/socilal-network-dialogs
```

      2. Run project with make

```plaintext
make 
```

     3.  Import postman collections from file Otus higload.postman\_collection.json

Some homeworks your might find in homework folder, unfortunatelly in russian.

**GOOSE COMMANDS:**

```plaintext
goose -dir=./migration  postgres "user=postgres host=localhost dbname=social_network_dialogs sslmode=disable password=mypassword" create create_user_table sql
```

```plaintext
goose -dir=./migration  postgres "user=postgres host=localhost dbname=social_network_dialogs sslmode=disable password=mypassword" up
```

**GOOSE COMMANDS:**
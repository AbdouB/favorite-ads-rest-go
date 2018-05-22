### This is a simple dockernized API using go + mongodb

After cloning the repo run the following commands to get up and running:

* `docker-compose up -d --build`

* `docker-compose up`

the server should be listening on port `3000`

### API endpoints

| METHOD / ENDPOINT        | PARAMS           | RESPONSE  |
| ------------- |:-------------:| -----|
| GET /users     | NONE | Will return an array of pre-seeded users: <br> [{_id: xxxx, email: xxx}, {},...]|
| GET /ads      | NONE      |   Will return an array of pre-seeded ads: <br> [{_id: xxxx, title: xxx, description: xxxx}, {},...] |
| GET /user/{user_id}/favoriteads/{save_type}      | user_id:  24 char hex (required) <br> save_type: manual \| automatic \| all (required)    |   array of user favorite ads:  [{_id: xxxx, title: xxx, description: xxxx, saveType: xxx}, {},...] |
| POST /user/favoriteads | adId: 24 char hex (required) <br> userId: 24 char hex (required) <br>  saveType: manual \| automatic (optional: will be replaced by "automatic" if another or no value is specified)  |    array of updated favorite ads: <br> [{_id: xxxx, title: xxx, description: xxxx, saveType: xxx}, {},...]|
| DELETE /user/favoriteads     |  adId: 24 char hex (required) <br> userId: 24 char hex (required) | array of updated favorite ads: <br> [{_id: xxxx, title: xxx, description: xxxx, saveType: xxx}, {},...]|

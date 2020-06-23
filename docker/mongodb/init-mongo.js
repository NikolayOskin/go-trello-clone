db.auth('root', 'root');

db = db.getSiblingDB('trello');

db.users.createIndex({ "email": 1 }, {unique: true});
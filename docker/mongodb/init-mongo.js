db.auth('root', 'root');

db.getSiblingDB('trello');
db.getSiblingDB('trello-test');

// db.users.createIndex({ "email": 1 }, {unique: true});
// db.boards.createIndex({ "user_id": 1 });
// db.lists.createIndex({ "board_id": 1 });
// db.cards.createIndex({ "board_id": 1 });
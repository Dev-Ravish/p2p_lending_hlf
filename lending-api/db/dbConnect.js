const mysql = require('mysql');

// Create a connection to the database
const connection = mysql.createConnection({
  host: 'localhost', // or the IP address of your Docker container
  port: 3306,        // the default port for MariaDB
  user: 'mariadbuser',
  password: 'yourPassword',
  database: 'mariadbdatabase'
});

// Open the MySQL connection
connection.connect(error => {
  if (error) {
    return console.error('error connecting: ' + error.stack);
  }
  console.log('connected as id ' + connection.threadId);
});

// Query the database

    const createEventsTable = `
      CREATE TABLE IF NOT EXISTS users(
        id INTEGER PRIMARY KEY AUTO_INCREMENT,
        name TEXT NOT NULL,
        email_id TEXT NOT NULL,
        account TEXT NOT NULL,
        user_id TEXT NOT NULL
      )
    `;
  
    connection.query(createEventsTable, (error, results, fields) => {
        if (error) {
          return console.error('error executing query: ' + error.stack);
        }
        console.log('Query results: ', results);
      });
// Close the connection

connection.end();

const axios = require("axios");

function createTable(db) {
  const createApartmentsTable = `
    DROP TABLE IF EXISTS apartments;

    CREATE TABLE IF NOT EXISTS apartments (
        id uuid primary key,
        apartment_name text,
        address text,
        noise_level text,
        floor text
    );`;
  db.exec(createApartmentsTable);
  refreshApartments(db);

  const createBookingsTable = `
    DROP TABLE IF EXISTS bookings;
    
    CREATE TABLE IF NOT EXISTS bookings (
      id uuid primary key,
      apartment_id text,
      user_id text,
      start_date text,
      end_date text
    );`;
  db.exec(createBookingsTable);
  refreshBookings(db);
}

async function refreshApartments(db) {
  const APARTMENT_URL = "http://localhost:3000/api/apartments";
  const response = await axios.get(APARTMENT_URL);

  const insert = db.prepare(
    "INSERT INTO apartments (id, apartment_name, address, noise_level, floor) VALUES (@id, @apartment_name, @address, @noise_level, @floor)"
  );

  const insertMany = db.transaction((apartments) => {
    for (const apt of apartments) insert.run(apt);
  });

  insertMany(response.data);
}

async function refreshBookings(db) {
  const BOOKING_URL = "http://localhost:3001/api/bookings";
  const response = await axios.get(BOOKING_URL);

  const insert = db.prepare(
    "INSERT INTO bookings (id, apartment_id, user_id, start_date, end_date) VALUES (@id, @apartment_id, @user_id, @start_date, @end_date)"
  );

  const insertMany = db.transaction((bookings) => {
    for (const booking of bookings) insert.run(booking);
  });

  insertMany(response.data);
}

function listAll(db, table) {
  const stmt = db.prepare("SELECT * FROM " + table);
  const row = stmt.all();
  return row;
}

// function createApartment(db, apartment) {
//   const insert = db.prepare(
//     "INSERT INTO apartments (id, apartment_name, address, noise_level, floor) VALUES (@id, @apartment_name, @address, @noise_level, @floor)"
//   );
//   db.run(insert, apartment);
// }

module.exports = {
  createTable: createTable,
  listAll: listAll,
};

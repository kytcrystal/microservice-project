const config = require("./config");
const axios = require("axios");

function createTable(db) {
  const dropTables = `
    DROP TABLE IF EXISTS bookings;
    DROP TABLE IF EXISTS apartments;
  `;

  const createApartmentsTable = `
    CREATE TABLE IF NOT EXISTS apartments (
        id uuid primary key,
        apartment_name text,
        address text,
        noise_level text,
        floor text
    );`;

  const createBookingsTable = `
    CREATE TABLE IF NOT EXISTS bookings (
      id uuid primary key,
      apartment_id uuid,
      user_id text,
      start_date text,
      end_date text,
      CONSTRAINT apartment_id
        FOREIGN KEY(apartment_id) 
	        REFERENCES apartments(id)
    );`;

  db.exec(dropTables);
  db.exec(createApartmentsTable);
  refreshApartments(db);
  db.exec(createBookingsTable);
  refreshBookings(db);
}

async function refreshApartments(db) {
  const response = await axios.get(config.APARTMENT_URL + "/api/apartments");

  const insert = db.prepare(
    `INSERT INTO apartments (id, apartment_name, address, noise_level, floor) 
      VALUES (@id, @apartment_name, @address, @noise_level, @floor)
      ON CONFLICT DO NOTHING`
  );

  const insertMany = db.transaction((apartments) => {
    for (const apt of apartments) insert.run(apt);
  });

  insertMany(response.data);
}

async function refreshBookings(db) {
  const response = await axios.get(config.BOOKING_URL + "/api/bookings");

  const insert = db.prepare(
    `INSERT INTO bookings (id, apartment_id, user_id, start_date, end_date) 
      VALUES (@id, @apartment_id, @user_id, @start_date, @end_date)
      ON CONFLICT DO NOTHING`
  );

  const insertMany = db.transaction((bookings) => {
    for (const booking of bookings) insert.run(booking);
  });

  insertMany(response.data);
}

function createApartment(db, apartment) {
  const insert = db.prepare(
    `INSERT INTO apartments (id, apartment_name, address, noise_level, floor) 
      VALUES (@id, @apartment_name, @address, @noise_level, @floor)
      ON CONFLICT DO NOTHING`
  );
  insert.run(apartment);
}

function deleteApartment(db, apartment) {
  db.prepare(`DELETE FROM bookings WHERE apartment_id = @id;`).run(apartment);
  db.prepare(`DELETE FROM apartments WHERE id = @id;`).run(apartment);
}

function createBooking(db, booking) {
  const insert = db.prepare(
    `INSERT INTO bookings (id, apartment_id, user_id, start_date, end_date) 
      VALUES (@id, @apartment_id, @user_id, @start_date, @end_date)
      ON CONFLICT DO NOTHING`
  );
  insert.run(booking);
}

function cancelBooking(db, booking) {
  const insert = db.prepare("DELETE FROM bookings WHERE id = @id");
  insert.run(booking);
}

function updateBooking(db, booking) {
  const insert = db.prepare(
    `UPDATE bookings SET (id, apartment_id, user_id, start_date, end_date) 
      VALUES (@id, @apartment_id, @user_id, @start_date, @end_date) WHERE id = @id`
  );
  insert.run(booking);
}

function listAll(db, table) {
  const stmt = db.prepare("SELECT * FROM " + table);
  const row = stmt.all();
  return row;
}

function searchAvailableApartments(db, fromDate, toDate) {
  const stmt = db.prepare(`SELECT * FROM apartments WHERE id NOT IN 
    (SELECT apartment_id FROM bookings WHERE 
    (start_date <= @fromDate AND end_date >= @fromDate) OR            -- bookings that include fromDate
    (start_date <= @toDate AND end_date >= @toDate) OR                -- bookings that include toDate
    (start_date >= @fromDate AND end_date <= @toDate))                -- bookings that included in [fromDate, toDate]
    `);
  const row = stmt.all({ fromDate, toDate });
  console.table(row);
  return row;
}

module.exports = {
  createTable: createTable,
  createApartment: createApartment,
  deleteApartment: deleteApartment,
  createBooking: createBooking,
  cancelBooking: cancelBooking,
  updateBooking: updateBooking,
  listAll: listAll,
  searchAvailableApartments: searchAvailableApartments,
};

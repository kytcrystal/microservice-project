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
      apartmentID text,
      userID text,
      startDate text,
      endDate text
    );`;
  db.exec(createBookingsTable);
  refreshBookings(db);
}

function refreshApartments(db) {
  const insert = db.prepare(
    "INSERT INTO apartments (id, apartment_name, address, noise_level, floor) VALUES (@id, @apartment_name, @address, @noise_level, @floor)"
  );

  const insertMany = db.transaction((apartments) => {
    for (const apt of apartments) insert.run(apt);
  });

  insertMany([
    {
      id: "1e0cfb4e-0a11-48c8-a1f5-e82f5587818c",
      apartment_name: "Always White",
      address: "Trento",
      noise_level: "5",
      floor: "1",
    },
    {
      id: "2f0cfb4e-0a11-48c8-a1f5-e82f5587818c",
      apartment_name: "Always Blue",
      address: "Bolzano",
      noise_level: "3",
      floor: "4",
    },
  ]);
}

function refreshBookings(db) {
  const insert = db.prepare(
    "INSERT INTO bookings (id, apartmentID, userID, startDate, endDate) VALUES (@id, @apartmentID, @userID, @startDate, @endDate)"
  );

  const insertMany = db.transaction((bookings) => {
    for (const booking of bookings) insert.run(booking);
  });

  insertMany([
    {
      id: "6e0cfb4e-0a11-48c8-a1f5-e82f5587818d",
      apartmentID: "1e0cfb4e-0a11-48c8-a1f5-e82f5587818c",
      userID: "M47730",
      startDate: "2024-11-01",
      endDate: "2024-11-23",
    },
    {
      id: "1e0cfb4e-0a11-48c8-a1f5-e82f5587818c",
      apartmentID: "2f0cfb4e-0a11-48c8-a1f5-e82f5587818c",
      userID: "M47730",
      startDate: "2024-02-01",
      endDate: "2024-03-01",
    },
  ]);
}

function listAll(db, table) {
  const stmt = db.prepare("SELECT * FROM " + table);
  const row = stmt.all();
  return row;
}

module.exports = {
  createTable: createTable,
  listAll: listAll,
};

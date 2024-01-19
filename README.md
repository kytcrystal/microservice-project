# microservice-project

## Overview

![Overview of the Project](./img/01-first-version-of-the-project.png)

## TODOs

- [X] Implement HTTP endpoints in apartment service that just accept requests as per sample file
- [X] Implement HTTP endpoints in booking service
- [X] we can try to check if an apartment exist before allowing a booking (direct communication with the other service)
- [X] Connect apartment service to a database
- [ ] When a new apartment is created, apartment service sends a rabbit mq message, booking service listen and create the appartment in it's own DB too
- [ ] Add search service with similar approach
- [ ] Dockerize applications
  - [x] Apartments
  - [ ] Bookings
  - [ ] Search
  - [ ] Gateway
  - [ ] Configuration for yaml file for Gateway
- [ ] In apartment, when adding new apartment, can check if id is passed in correctly. If yes, use that. If not, generate new UUID
- [ ] Booking should not directly call Apartment to find out if an apartment exists. This is why you need to add a message queue
  - [ ] Apartment post event (apartment added and deleted) to queue
  - [ ] Booking register for apartment events
- [ ] Search needs to know which apartment exists and are available
- [ ] Booking post event to another queue (booking added, changed and cancelled)
- [ ] Search register for booking events
- [ ] Search register for apartment events
- [ ] Search should be able to search apartments using "from" date and "to" date -> should we do in python?
- [ ] Search needs to have a DB with apartments and bookings -> should it be NoSQL?
- [ ] Direct call from Search service to Apartment service, if apartments table is empty
- [ ] Direct call from Search service to Booking service, if bookings table is empty
- [ ] Implement event sourcing for Booking service
- [ ] Docker multistage build




## Requirements

- Go: 1.21.4
- Docker
  

## Useful Resources

- [Diagram in Visio](https://scientificnet-my.sharepoint.com/:u:/r/personal/mponza_unibz_it/Documents/CPD%20-%20Microservices%20Project.vsdx?d=w6328c77940f14158bfbf177a6352d738&csf=1&web=1&e=2ctcRj)
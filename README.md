# microservice-project

## Overview

![Overview of the Project](./img/01-first-version-of-the-project.png)

## TODOs

- [X] Implement HTTP endpoints in apartment service that just accept requests as per sample file
- [X] Implement HTTP endpoints in booking service
- [ ] we can try to check if an apartment exist before allowing a booking (direct communication with the other service)
- [ ] Connect apartment service to a database
- [ ] When a new appartment is created appartment service sends a rabbit mq message, booking service listen and create the appartment in it's own DB too;
- [ ] Add search service with similar approach
- [ ] Dockerize everything
- [ ] In apartment, when adding new apartment, can check if id is passed in correctly. If yes, use that. If not, generate new UUID

## Requirements

- Go: 1.21.4
- Docker
  

## Useful Resources

- [Diagram in Visio](https://scientificnet-my.sharepoint.com/:u:/r/personal/mponza_unibz_it/Documents/CPD%20-%20Microservices%20Project.vsdx?d=w6328c77940f14158bfbf177a6352d738&csf=1&web=1&e=2ctcRj)


Questions
- Why when i query "http://localhost:3000/api/apartments" in webpage, in console, it is printed both
    got /api/apartments GET request
    got / request
- How do i add the post request body into the URL?
- 
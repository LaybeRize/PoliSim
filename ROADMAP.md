# Roadmap

This is the offical Roadmap for the **PoliSim** Project.

## Goals for Beta

These are the goals which have to be reached to transfer the project from the Alpha-Phase to the Beta-Phase.

- tests for most logic (ca 60%)
- bug testing with users
- implementation of the Zwitscher-System
- basic general-purpose API for JSON

### Tests

The 60% testcov are not project wide but inside the data folder for now. There the tests should cover most non-database
dependent logic. For that reason the project should decouple database-reliant and database-indendent validation/processing 
logic. If possible all database-independet logic should be covered.

### Bug testing

No specifics can be made currently to indicate when this phase is sufficently 
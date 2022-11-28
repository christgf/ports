# Technical test

- Given a file with ports data (ports.json), write a port domain service that either creates a new record in a database, or updates the existing one (Hint: no need for delete or other methods).
- The file is of unknown size, it can contain several millions of records, you will not be able to read the entire file at once.
- The service has limited resources available (e.g. 200MB ram).
- The end result should be a database containing the ports, representing the latest version found in the JSON. (Hint: use an in memory database to save time and avoid complexity).
- A Dockerfile should be used to contain and run the service (Hint: extra points for avoiding compilation in docker).
- Provide at least one example per test type that you think are needed for your assignment. This will allow the reviewer to evaluate your critical thinking as well as your knowledge about testing.
- Your readme.md should explain how to run your program and test it.
- The service should handle certain signals correctly (e.g. a TERM or KILL signal should result in a graceful shutdown).

Choose the approach that you think is best (i.e. most flexible).
# Record API Specification

## Requirements

### Requirement: Create record
The system SHALL allow a client to create a vinyl record. `title` and `artist`
are required; `label`, `format`, `releaseYear`, `country`, `condition`, and
`notes` are optional.

#### Scenario: Valid record is created
- GIVEN a client submits a valid JSON request with at least `title` and `artist`
- WHEN the client sends `POST /records`
- THEN the system SHALL persist the record
- AND return `201 Created`
- AND return the created record including `id`, `title`, `artist`, `createdAt`, and `updatedAt`

#### Scenario: Optional metadata is provided
- GIVEN a client submits a request including any of `label`, `format`, `releaseYear`, `country`, `condition`, or `notes`
- WHEN the client sends `POST /records`
- THEN the system SHALL persist the provided optional fields
- AND return them in the created record

#### Scenario: Request body is invalid
- GIVEN a client submits a malformed or non-conforming JSON body
- WHEN the client sends `POST /records`
- THEN the system SHALL return `400 Bad Request`
- AND return an error response

### Requirement: List records
The system SHALL allow a client to list all vinyl records.

#### Scenario: Records exist
- GIVEN one or more records exist
- WHEN the client sends `GET /records`
- THEN the system SHALL return `200 OK`
- AND return a JSON array of records

#### Scenario: No records exist
- GIVEN no records exist
- WHEN the client sends `GET /records`
- THEN the system SHALL return `200 OK`
- AND return an empty JSON array

### Requirement: Search records
The system SHALL allow a client to search records using a `q` query parameter
that matches against `title`, `artist`, `label`, and `notes`.

#### Scenario: Matching records exist
- GIVEN one or more records match the search term
- WHEN the client sends `GET /records?q={term}`
- THEN the system SHALL return `200 OK`
- AND return a JSON array containing only records whose `title`, `artist`, `label`, or `notes` match the term

#### Scenario: No records match
- GIVEN no records match the search term
- WHEN the client sends `GET /records?q={term}`
- THEN the system SHALL return `200 OK`
- AND return an empty JSON array

#### Scenario: No search term provided
- GIVEN the `q` parameter is omitted
- WHEN the client sends `GET /records`
- THEN the system SHALL return all records

### Requirement: Get record by ID
The system SHALL allow a client to retrieve a single vinyl record by ID.

#### Scenario: Record exists
- GIVEN a record exists with a known ID
- WHEN the client sends `GET /records/{id}`
- THEN the system SHALL return `200 OK`
- AND return the matching record

#### Scenario: Record does not exist
- GIVEN no record exists for the requested ID
- WHEN the client sends `GET /records/{id}`
- THEN the system SHALL return `404 Not Found`
- AND return an error response

### Requirement: Update record
The system SHALL allow a client to update an existing vinyl record. The request
accepts the same fields as record creation (`title` and `artist` required;
`label`, `format`, `releaseYear`, `country`, `condition`, and `notes` optional).

#### Scenario: Record exists
- GIVEN a record exists with a known ID
- WHEN the client sends `PUT /records/{id}` with at least `title` and `artist`
- THEN the system SHALL update the record
- AND update `updatedAt`
- AND return `200 OK`
- AND return the updated record

#### Scenario: Record does not exist
- GIVEN no record exists for the requested ID
- WHEN the client sends `PUT /records/{id}`
- THEN the system SHALL return `404 Not Found`

#### Scenario: Request body is invalid
- GIVEN a client submits a malformed or non-conforming JSON body
- WHEN the client sends `PUT /records/{id}`
- THEN the system SHALL return `400 Bad Request`
- AND return an error response

### Requirement: Delete record
The system SHALL allow a client to delete an existing vinyl record.

#### Scenario: Record exists
- GIVEN a record exists with a known ID
- WHEN the client sends `DELETE /records/{id}`
- THEN the system SHALL delete the record
- AND return `204 No Content`

#### Scenario: Record does not exist
- GIVEN no record exists for the requested ID
- WHEN the client sends `DELETE /records/{id}`
- THEN the system SHALL return `404 Not Found`

### Requirement: Record model
The system SHALL represent records with the following fields. `id`, `title`,
`artist`, `createdAt`, and `updatedAt` are always present; `label`, `format`,
`releaseYear`, `country`, `condition`, and `notes` are optional.

- `id` (UUID)
- `title` (string)
- `artist` (string)
- `label` (string, optional)
- `format` (string, optional, e.g. `LP`)
- `releaseYear` (integer, optional, e.g. `1977`)
- `country` (string, optional)
- `condition` (string, optional, e.g. `Very Good Plus`)
- `notes` (string, optional)
- `createdAt` (date-time)
- `updatedAt` (date-time)

#### Scenario: Record is returned
- GIVEN any endpoint returns a record
- WHEN the response is serialized as JSON
- THEN the response SHALL include `id`, `title`, `artist`, `createdAt`, and `updatedAt`
- AND SHALL include any optional fields that have values

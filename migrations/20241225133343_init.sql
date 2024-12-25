-- +goose Up
-- +goose StatementBegin
CREATE TABLE queries (
    id SERIAL PRIMARY KEY,           -- Primary key, auto-incremented
    departure VARCHAR(255) NOT NULL, -- Departure location
    destination VARCHAR(255) NOT NULL, -- Destination location
    stay_duration INT NOT NULL,      -- Stay duration in days
    month_horizon INT NOT NULL,      -- Month horizon
    weekdays VARCHAR(13) NOT NULL    -- Stores the weekdays as a string
);

CREATE TABLE requests (
    id SERIAL PRIMARY KEY,                      -- Primary key, auto-incremented
    id_query INT NOT NULL REFERENCES queries(id), -- Foreign key referencing the `queries` table
    departure VARCHAR(255) NOT NULL,            -- Departure location
    destination VARCHAR(255) NOT NULL,          -- Destination location
    departure_date DATE NOT NULL,               -- Departure date (without time)
    return_date DATE NOT NULL,                  -- Return date (without time)
    preferred_departure VARCHAR(10) DEFAULT 'ANY' CHECK (preferred_departure IN ('MORNING', 'EVENING', 'ANY')), -- Preferred departure time
    preferred_return VARCHAR(10) DEFAULT 'ANY' CHECK (preferred_return IN ('MORNING', 'EVENING', 'ANY'))        -- Preferred return time
);

CREATE TABLE flights (
    id BIGSERIAL PRIMARY KEY,                     -- Primary key, auto-incremented
    id_request INT NOT NULL REFERENCES requests(id), -- Foreign key referencing the `requests` table
    dep_date DATE NOT NULL,                    -- Departure date (without time)
    dep_time TIME NOT NULL,                    -- Departure time
    arr_time TIME NOT NULL,                    -- Arrival time
    airport VARCHAR(255),                      -- Airport
    company VARCHAR(255),                      -- Airline company
    duration VARCHAR(50) NOT NULL,             -- Flight duration (e.g., "6h 30m")
    price INT NOT NULL,                        -- Price of the flight
    currency VARCHAR(25) NOT NULL,             -- Currency (e.g., "USD")
    stops SMALLINT NOT NULL,                   -- Number of stops (non-negative integer)
    flight_type SMALLINT NOT NULL CHECK (flight_type IN (0, 1)) -- 0 for outbound, 1 for inbound
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE flights;
DROP TABLE requests;
DROP TABLE queries;
-- +goose StatementEnd

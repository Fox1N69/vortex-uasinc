CREATE TABLE IF NOT EXISTS algorithm_status (
    id SERIAL PRIMARY KEY,
    client_id INT NOT NULL,
    vwap BOOLEAN NOT NULL DEFAULT false ,
    twap BOOLEAN NOT NULL  DEFAULT false,
    hft BOOLEAN NOT NULL  DEFAULT false,
    CONSTRAINT fk_client
        FOREIGN KEY(client_id) 
	    REFERENCES clients(id)
	    ON DELETE CASCADE
);

CREATE INDEX idx_client_id ON algorithm_status (client_id);


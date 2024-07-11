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

-- Create index for vwap
CREATE INDEX idx_algorithm_status_vwap ON algorithm_status(vwap);

-- Create index for twap
CREATE INDEX idx_algorithm_status_twap ON algorithm_status(twap);

-- Create index for hft
CREATE INDEX idx_algorithm_status_hft ON algorithm_status(hft);

-- Create composite index for vwap, twap, hft
CREATE INDEX idx_algorithm_status_vwap_twap_hft ON algorithm_status(vwap, twap, hft);
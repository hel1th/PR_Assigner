
CREATE TABLE teams (
    name VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    team_id VARCHAR(255) REFERENCES teams(name),
    is_active BOOL NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pull_requests (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(1024) NOT NULL,
    author_id VARCHAR(255) REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(128) NOT NULL,
    merged_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE reviewers (
    pr_id VARCHAR(255) NOT NULL REFERENCES pull_requests(id),
    reviewer_id VARCHAR(255) NOT NULL REFERENCES users(id),
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (pr_id, reviewer_id)
);

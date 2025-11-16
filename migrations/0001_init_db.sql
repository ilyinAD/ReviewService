-- +goose Up
-- +goose StatementBegin
begin;

CREATE table teams (
                       team_name varchar(50) primary key
);

CREATE TABLE users (
                       id varchar(50) primary key,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       team_name varchar(50) REFERENCES teams (team_name),
                       is_active boolean NOT NULL

);

CREATE table pull_requests (
                               pull_request_id varchar(50) primary key,
                               pull_request_name varchar(50) NOT NULL,
                               author_id varchar(50) REFERENCES users(id),
                               status varchar(10) DEFAULT 'OPEN',
                               created_at timestamp DEFAULT CURRENT_TIMESTAMP,
                               merged_at timestamp default NULL
);

CREATE table reviewers (
                           reviewer_id varchar(50),
                           pull_request_id varchar(50) NOT NULL REFERENCES pull_requests(pull_request_id),
                           PRIMARY KEY (pull_request_id, reviewer_id)
);

end;

-- +goose StatementEnd

-- +goose Down
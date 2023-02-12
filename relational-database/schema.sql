CREATE EXTENSION IF NOT EXISTS "uuid-ossp";             -- Import UUID extension

CREATE TABLE IF NOT EXISTS account (                    -- Account Table
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),    -- Holds account information
    email VARCHAR UNIQUE NOT NULL,
    username VARCHAR UNIQUE NOT NULL,
    password VARCHAR,
    verified BOOLEAN NOT NULL,
    verification_code VARCHAR,
    otp_enabled BOOLEAN NOT NULL DEFAULT false,
    otp_secret VARCHAR,
    otp_auth_url VARCHAR,
    fname VARCHAR,
    lname VARCHAR
);

CREATE TABLE IF NOT EXISTS folder (                     -- Folder Table
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),    -- Holds folder data
    parent_id UUID,                                     -- Provides folder tree structure using self referencing relationship
    name VARCHAR NOT NULL,
    CONSTRAINT fk_parent
        FOREIGN KEY(parent_id)
        REFERENCES folder(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS file (                       -- File Table
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),    -- Holds file data
    name VARCHAR NOT NULL,
    encrypted_key BYTEA NOT NULL,
    asymmetric_key VARCHAR NOT NULL,
    hash_value VARCHAR NOT NULL,
    folder_id UUID,
    sensitive BOOLEAN NOT NULL DEFAULT false,
    CONSTRAINT fk_folder
        FOREIGN KEY(folder_id)
        REFERENCES folder(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS csrf (                       -- CSRF Table
    user_id UUID PRIMARY KEY NOT NULL,                       -- Holds CSRF token
    csrf_token VARCHAR NOT NULL,
    expire_datetime TIMESTAMP NOT NULL
);

/*
Permissions Tables between accounts or roles and folders or files

Permissions are mapped as:
0 --> no permissions
1 --> viewing
2 --> editing
3 --> administrator
4 --> owner
*/

CREATE TABLE IF NOT EXISTS account_folder (
    account_id UUID NOT NULL,
    folder_id UUID NOT NULL,
    permission SMALLINT NOT NULL,
    PRIMARY KEY (account_id, folder_id),
    CONSTRAINT fk_account
        FOREIGN KEY (account_id)
        REFERENCES account(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_folder
        FOREIGN KEY (folder_id)
        REFERENCES folder(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS account_file (
    account_id UUID NOT NULL,
    file_id UUID NOT NULL,
    permission SMALLINT NOT NULL,
    PRIMARY KEY (account_id, file_id),
    CONSTRAINT fk_account
        FOREIGN KEY(account_id)
        REFERENCES account(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_file
        FOREIGN KEY(file_id)
        REFERENCES file(id)
        ON DELETE CASCADE
);
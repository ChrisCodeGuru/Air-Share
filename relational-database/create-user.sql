CREATE OR REPLACE FUNCTION create_user(                          -- Stored procedure used for signup API
    _email VARCHAR,
    _username VARCHAR,
    _password VARCHAR,
    _plan SMALLINT,
    OUT _id VARCHAR
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO account(email, username, password, plan) 
    VALUES(_email, _username, _password, _plan)
    RETURNING id INTO _id;
END;
$$;
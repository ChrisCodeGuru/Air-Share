CREATE OR REPLACE FUNCTION get_user_details(                                -- Stored procedure used for login API
    _login_details VARCHAR                                                  -- Login details can be either a username or password
)
RETURNS TABLE (
    _id VARCHAR,
    _password VARCHAR
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT id, password
    FROM account 
    WHERE username=_username OR email=_email;
END;
$$;

CREATE OR REPLACE FUNCTION get_user_email_password(                                -- Stored procedure used for login API
    _account_id VARCHAR                                                  -- Login details can be either a username or password
)
RETURNS TABLE (
    _email VARCHAR,
    _password VARCHAR
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT email, password 
    FROM account 
    WHERE id = _account_id::UUID;
END;
$$;

CREATE OR REPLACE FUNCTION get_user_email(                                -- Stored procedure used for login API
    _account_id VARCHAR                                                  -- Login details can be either a username or password
)
RETURNS VARCHAR
LANGUAGE plpgsql
AS $$
DECLARE
    _email VARCHAR;
BEGIN
    SELECT email
    INTO _email
    FROM account 
    WHERE id = _account_id::UUID;

    RETURN _email;
END;
$$;
CREATE OR REPLACE FUNCTION account_authentication_check(
    _account_id VARCHAR
)
RETURNS TABLE (
    _password_existence BOOLEAN,
    _otp_status BOOLEAN
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT (password IS NOT NULL), otp_enabled 
    FROM account 
    WHERE id = _account_id::UUID;
END;
$$;

CREATE OR REPLACE FUNCTION account_password_check(
    _account_id VARCHAR
)
RETURNS BIGINT
LANGUAGE plpgsql
AS $$
DECLARE
    _count BIGINT;
BEGIN
    SELECT COUNT(*)
    INTO _count
    FROM account 
    WHERE id = _account_id::UUID
    AND PASSWORD IS NULL;

    RETURN _count;
END;
$$;
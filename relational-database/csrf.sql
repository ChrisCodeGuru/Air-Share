CREATE OR REPLACE FUNCTION csrf_check(
    _account_id VARCHAR
)
RETURNS TABLE (
    _csrf VARCHAR,
    _year TEXT,
    _month TEXT,
    _day TEXT,
    _hours TEXT,
    _minutex TEXT,
    _seconds TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT csrf_token, (TO_CHAR(expire_datetime, 'YYYY')), (TO_CHAR(expire_datetime, 'MM')), (TO_CHAR(expire_datetime, 'DD')), (TO_CHAR(expire_datetime, 'HH24')), (TO_CHAR(expire_datetime, 'MI')), (TO_CHAR(expire_datetime, 'SS')) 
    FROM csrf 
    WHERE user_id = _account_id::UUID;
END;
$$;
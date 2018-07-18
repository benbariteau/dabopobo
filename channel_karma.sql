CREATE TABLE channel_karma (
    identifier TEXT NOT NULL,
    channel TEXT NOT NULL,
    date_bucket TEXT NOT NULL,
    plusplus INTEGER NOT NULL,
    minusminus INTEGER NOT NULL,
    plusminus INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS temp (
    title TEXT
);

INSERT INTO temp VALUES ('test1'), ('test2');

select * from temp;

delete from temp;

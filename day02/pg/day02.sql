begin;

create table outcome (
    outcome text primary key,
    score int not null
);

insert into outcome (outcome, score) values
    ( 'lose', 0 ),
    ( 'draw', 3 ),
    ( 'win',  6 );

create table move (
    move text primary key,
    score int not null
);

insert into move (move, score) values
    ( 'rock',     1 ),
    ( 'paper',    2 ),
    ( 'scissors', 3 );

create table round (
    you     text not null references move(move),
    me      text not null references move(move),
    outcome text not null references outcome(outcome)
);

insert into round (you, me, outcome) values
  ( 'rock',     'rock',     'draw' ),
  ( 'rock',     'paper',    'win'  ),
  ( 'rock',     'scissors', 'lose' ),
  ( 'paper',    'rock',     'lose' ),
  ( 'paper',    'paper',    'draw' ),
  ( 'paper',    'scissors', 'win'  ),
  ( 'scissors', 'rock',     'win'  ),
  ( 'scissors', 'paper',    'lose' ),
  ( 'scissors', 'scissors', 'draw' );

create table part1 (
    you     text not null references move(move),
    me      text not null references move(move)
);

\copy part1 from 'part1.csv' delimiter ',' csv header;

with scores as (
  select m.score + o.score as score
    from part1 p
    join round r using (you, me)
    join move m on p.me = m.move
    join outcome o using (outcome)
) select sum(score) from scores;

create table part2 (
    you     text not null references move(move),
    outcome text not null references outcome(outcome)
);

\copy part2 from 'part2.csv' delimiter ',' csv header;

with scores as (
  select m.score + o.score as score
    from part2 p
    join round r using (you, outcome)
    join move m on r.me = m.move
    join outcome o using (outcome)
) select sum(score) from scores;

rollback;

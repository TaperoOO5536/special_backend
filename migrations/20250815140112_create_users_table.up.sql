create table if not exists User (
            ID_User serial not null constraint PK_User primary key,
            S_N_User varchar (50) not null,
            F_N_User varchar (50) not null,
            N_N_User varchar (50) not null,
            Phone_N_Representative varchar (16) not null,
        );
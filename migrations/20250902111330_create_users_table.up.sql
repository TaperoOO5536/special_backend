create table if not exists Users (
            ID_User text constraint PK_User primary key default uuid_generate_v4(),
            S_N_User varchar (50) not null,
            F_N_User varchar (50) not null,
            N_N_User varchar (50) not null,
            Phone_N_User varchar (16) not null
        );
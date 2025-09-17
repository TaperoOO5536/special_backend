create table if not exists Ivents (
            ID_Ivent UUID constraint PK_Ivent primary key default uuid_generate_v4(),
            Ivent_Title varchar (50) not null,
            Ivent_Description varchar (300) not null,
            Ivent_DateTime timestamp with time zone not null,
            Ivent_Price int not null,
            Total_Seats int null,
            Occupied_Seats int null,
            Little_Picture bytea not null
        ); 
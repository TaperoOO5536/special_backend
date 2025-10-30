create table if not exists Events (
            ID_Event UUID constraint PK_Event primary key default uuid_generate_v4(),
            Event_Title varchar (50) not null,
            Event_Description varchar (300) not null,
            Event_DateTime timestamp with time zone not null,
            Event_Price int not null,
            Total_Seats int not null,
            Occupied_Seats int not null,
            Little_Picture text not null
        ); 
create table if not exists User_Events (
            ID_User_Event UUID constraint PK_User_Event primary key default uuid_generate_v4(),
            User_ID text not null references Users (ID_User),
            Event_ID UUID not null references Events (ID_Event),
            Number_Of_Guests int not null
        ); 
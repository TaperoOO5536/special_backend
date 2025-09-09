create table if not exists User_Ivents (
            ID_User_Ivent UUID constraint PK_User_Ivent primary key default uuid_generate_v4(),
            User_ID text not null references Users (ID_User),
            Ivent_ID UUID not null references Ivents (ID_Ivent),
            Number_Of_Guests int not null
        ); 
create table if not exists Orders (
            ID_Order UUID constraint PK_Order primary key default uuid_generate_v4(),
            Order_Number varchar (20) not null,
            Order_Form_DateTime timestamp with time zone not null,
            Completion_Date timestamp with time zone not null,
            Order_Comment varchar (300) null,
            Order_Amount int not null,
            Order_Status varchar (50) not null,
            User_ID text not null references Users (ID_User)
        );
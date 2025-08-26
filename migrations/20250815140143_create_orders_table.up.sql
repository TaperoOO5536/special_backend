create table if not exists Order (
            ID_Order UUID constraint PK_Order primary key default uuid_generate_v4(),
            Order_Number varchar (20) not null,
            Order_Form_Date date not null,
            Order_Form_Time time not null,
            Completion_Date date not null,
            `Comment` string null,
            User_ID int not null references User (ID_User),
        );
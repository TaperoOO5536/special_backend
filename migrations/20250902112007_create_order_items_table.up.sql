create table if not exists Order_Items (
            ID_Order_Item UUID constraint PK_Order_Item primary key default uuid_generate_v4(),
            Order_ID UUID not null references Orders (ID_Order),
            Item_ID UUID not null references Items (ID_Item),
            Quantity int not null
        ); 